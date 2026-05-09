package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
	"github.com/sirupsen/logrus"
)

const GEOCODE_BASE_URL = "https://maps.googleapis.com/maps/api/geocode/json"

func NewService(log *logrus.Logger, cfg *AppConfig) *AppService {
	return &AppService{
		log:      log,
		cfg:      cfg,
		features: make(map[string]*geojson.Feature),
		rioneMap: make(map[string]string),
	}
}

func (s *AppService) LoadGeoJSONMapFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read GeoJSON file: %w", err)
	}

	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		return fmt.Errorf("failed to parse GeoJSON: %w", err)
	}

	// Index features and extract rione names
	for i, feature := range fc.Features {
		if feature == nil || feature.Geometry == nil {
			continue
		}

		// Extract rione name from feature properties
		var rioneName string
		if name, ok := feature.Properties["Name"]; ok {
			rioneName = fmt.Sprintf("%v", name)
		}

		// Use rione name as ID, fallback to generic ID
		id := rioneName
		if id == "" {
			id = fmt.Sprintf("feature_%d", i)
		}

		s.features[id] = feature
		s.rioneMap[id] = rioneName

		s.log.Debugf("loaded feature: %s", id)
	}

	s.log.Infof("loaded %d features (rioni) from GeoJSON", len(s.features))
	return nil
}

// findRioneForPoint checks which rione polygon contains the given point
// Returns the rione name if found, empty string otherwise
func (s *AppService) findRioneForPoint(point Point) string {
	orbPoint := orb.Point{point.Longitude, point.Latitude}
	s.log.Infof("checking point: lon=%.6f, lat=%.6f", point.Longitude, point.Latitude)

	for rioneID, feature := range s.features {
		if feature.Geometry == nil {
			s.log.Debugf("  %s: no geometry", rioneID)
			continue
		}

		// Try pointer first, then value
		var poly *orb.Polygon
		if p, ok := feature.Geometry.(*orb.Polygon); ok {
			poly = p
		} else if p, ok := feature.Geometry.(orb.Polygon); ok {
			poly = &p
		}

		if poly == nil {
			s.log.Debugf("  %s: not a polygon (type: %T)", rioneID, feature.Geometry)
		}

		contains := planar.PolygonContains(*poly, orbPoint)
		if contains {
			rioneName := s.rioneMap[rioneID]
			s.log.Infof("point (%.6f, %.6f) found in rione: %s", point.Latitude, point.Longitude, rioneName)
			return rioneName
		}
	}

	s.log.Infof("point (%.6f, %.6f) not found in any rione", point.Latitude, point.Longitude)
	return ""
}

// isMapLoaded checks if any features have been loaded
func (s *AppService) isMapLoaded() bool {
	return len(s.features) > 0
}

func (s *AppService) geocode(address string) (*GeocodeResponse, error) {
	endpoint := URLWithParams(GEOCODE_BASE_URL, map[string]string{
		"address": fmt.Sprintf("%s %s", address, "Carmiano, Lecce"),
		"key":     s.cfg.MapsApiKey,
	})

	s.log.Debugf("executing GET %s", endpoint)
	res, err := http.Get(endpoint)

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non 200 status code %d", res.StatusCode)
	}

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var geocode GeocodeResponse
	if err := json.NewDecoder(res.Body).Decode(&geocode); err != nil {
		return nil, err
	}

	return &geocode, nil
}

func (s *AppService) processRioneAddressRequest(c *gin.Context) (*RioneAddressRequest, error) {
	addr := c.Query("address")

	if len(addr) == 0 {
		return nil, fmt.Errorf("missing required 'address' param")
	}

	return &RioneAddressRequest{Address: addr}, nil
}

func isInCarmiano(g GeocodeResult) bool {
	for _, comp := range g.AddressComponents {
		if comp.LongName == "Carmiano" {
			return true
		}
	}

	return false
}

func (s *AppService) parseGeocodeResult(g *GeocodeResult) RioneAddressResponse {
	var street, number string

	for _, comp := range g.AddressComponents {
		if Includes(comp.Types, "route") {
			street = comp.LongName
		}

		if Includes(comp.Types, "street_number") {
			number = comp.LongName
		}
	}

	rione := ""
	if street != "" {
		rione = s.findRioneForPoint(Point{
			Latitude:  g.Geometry.Location.Lat,
			Longitude: g.Geometry.Location.Lng,
		})
	}

	return RioneAddressResponse{
		Street:    street,
		Number:    number,
		Latitude:  g.Geometry.Location.Lat,
		Longitude: g.Geometry.Location.Lng,
		Rione:     rione,
	}
}

func (s *AppService) TrovaRionePerIndirizzo(c *gin.Context) {
	req, err := s.processRioneAddressRequest(c)

	if !s.isMapLoaded() {
		s.log.Errorln("mappa non caricata, impossibile elaborare richiesta")
		c.JSON(http.StatusInternalServerError, ErrorBody("Servizio non disponibile"))
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorBody(err.Error()))
		return
	}

	geocode, err := s.geocode(req.Address)
	if err != nil {
		s.log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, ErrorBody("Servizio non disponibile"))
		return
	}

	if geocode.Status == GeocodeStatusZeroResults || len(geocode.Results) == 0 {
		c.JSON(http.StatusNotFound, ErrorBody("Indirizzo non trovato"))
		return
	}

	var found *GeocodeResult
	for _, geoResult := range geocode.Results {
		if isInCarmiano(geoResult) {
			found = &geoResult
			break
		}
	}

	if found == nil {
		c.JSON(http.StatusNotFound, ErrorBody("Indirizzo non trovato a Carmiano"))
		return
	}

	parsed := s.parseGeocodeResult(found)

	if parsed.Rione == "" {
		s.log.Errorf("impossibile determinare rione: %+v", parsed)
		c.JSON(http.StatusNotFound, ErrorBody("Impossibile determinare Rione"))
		return
	}

	c.JSON(http.StatusOK, parsed)
}
