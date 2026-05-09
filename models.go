package main

import (
	"github.com/paulmach/orb/geojson"
	"github.com/sirupsen/logrus"
)

type GeocodeStatus string

const (
	GeocodeStatusOK          GeocodeStatus = "OK"
	GeocodeStatusZeroResults GeocodeStatus = "ZERO_RESULTS"
)

type GeocodeResponse struct {
	Status  GeocodeStatus   `json:"status"`
	Results []GeocodeResult `json:"results"`
}

type GeocodeResult struct {
	FormattedAddress  string             `json:"formatted_address"`
	PlaceID           string             `json:"place_id"`
	Types             []string           `json:"types"`
	AddressComponents []AddressComponent `json:"address_components"`
	Geometry          Geometry           `json:"geometry"`
	NavigationPoints  []NavigationPoint  `json:"navigation_points"`
}

type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type Geometry struct {
	Location     LatLng `json:"location"`
	LocationType string `json:"location_type"`
	Viewport     Bounds `json:"viewport"`
	Bounds       Bounds `json:"bounds"`
}

type Bounds struct {
	Northeast LatLng `json:"northeast"`
	Southwest LatLng `json:"southwest"`
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type NavigationPoint struct {
	Location NavLatLng `json:"location"`
}

// Navigation points use "latitude"/"longitude" instead of "lat"/"lng"
type NavLatLng struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type RioneAddressRequest struct {
	Address string
}

type RioneAddressResponse struct {
	Street    string  `json:"street"`
	Number    string  `json:"number"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Rione     string  `json:"rione"`
}

type Point struct {
	Latitude  float64
	Longitude float64
}

type AppService struct {
	log      *logrus.Logger
	cfg      *AppConfig
	features map[string]*geojson.Feature // ID to feature mapping
	rioneMap map[string]string           // feature ID to rione name
}
