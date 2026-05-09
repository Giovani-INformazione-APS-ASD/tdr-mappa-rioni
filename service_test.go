package main

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

type TestCase struct {
	address       string
	expectedRione string
	description   string
}

func mustBe(s1, s2 string) TestCase {
	return TestCase{
		address:       s1,
		expectedRione: s2,
	}
}

// TODO: Trovare dei Test più critici
var testCases = []TestCase{
	// Aia
	mustBe("Bar Piccolo", "Aia"),
	mustBe("Ottoemezzo", "Aia"),
	mustBe("Via Goldoni", "Aia"),
	mustBe("Viaggia con Noi", "Aia"),
	mustBe("Sapori di Mare del Salento", "Aia"),
	mustBe("Elite Legno", "Aia"),
	mustBe("Civoleva 196", "Aia"),
	mustBe("Lorenzauto", "Aia"),
	mustBe("Il Nido", "Aia"),
	mustBe("Pneumatic Center", "Aia"),

	// Centro
	mustBe("Li Birbanti", "Centro"),
	mustBe("Bracia", "Centro"),
	mustBe("Piazza Assunta", "Centro"),
	mustBe("Bar Diaz", "Centro"),
	mustBe("Farmacia Barbagallo Edvige", "Centro"),
	mustBe("Bar Rizzo", "Centro"),
	mustBe("Via Cesare Battisti 14", "Centro"),
	mustBe("Barberia Petrelli", "Centro"),
	mustBe("Via Gorizia", "Centro"),
	mustBe("Centro Polivalente per Minori NOI", "Centro"),

	// Immacolata
	mustBe("Vie Lecce 89", "Immacolata"),
	mustBe("Cineteatro Lumiere", "Immacolata"),
	mustBe("Via Libertini 2/B", "Immacolata"),
	mustBe("Primopiano", "Immacolata"),
	mustBe("Diesel Cafè", "Immacolata"),
	mustBe("Pasticceria Perrone", "Immacolata"),
	mustBe("Chiesa dell'Immacolata", "Immacolata"),
	mustBe("Via Immacolata", "Immacolata"),
	mustBe("Via Stazione", "Immacolata"),
	mustBe("Via Villafranca 2", "Immacolata"),

	// Mendule Mare
	mustBe("SP12 21/e", "Mendule Mare"),
	mustBe("Via Grassi 46", "Mendule Mare"),
	mustBe("Via Gentile 25", "Mendule Mare"),
	mustBe("Via Ada Negri 1", "Mendule Mare"),
	mustBe("La Maglianese", "Mendule Mare"),
	mustBe("Piazza degli Eroi", "Mendule Mare"),
	mustBe("Parco della Scienza", "Mendule Mare"),
	mustBe("Via Trappeto", "Mendule Mare"),
	mustBe("RA Home Solutions", "Mendule Mare"),
	mustBe("Via Gentile", "Mendule Mare"),

	// Pitrignani
	mustBe("La Bella Vita", "Pitrignani"),
	mustBe("Via Venezia", "Pitrignani"),
	mustBe("Via Papa Pio XII", "Pitrignani"),
	mustBe("Ecocentro di Carmiano", "Pitrignani"),
	mustBe("Via Fratelli Carioli 4", "Pitrignani"),
	mustBe("Via Filippo Turati 8", "Pitrignani"),
	mustBe("D'Agostino Elettrodomestici", "Pitrignani"),
	mustBe("Ristorante Pizzeria Salento", "Pitrignani"),
	mustBe("Moda Più", "Pitrignani"),

	// Quatarari
	mustBe("Via Giacomo Leopardi", "Quatarari"),
	mustBe("Nuova OMET", "Quatarari"),
	mustBe("Colorificio 2C", "Quatarari"),
	mustBe("Via Donizetti 4", "Quatarari"),
	mustBe("Via Po 45", "Quatarari"),
	mustBe("Via Monteroni 15", "Quatarari"),
	mustBe("Via Garda 7", "Quatarari"),
	mustBe("Via Pitero Mascagni 6", "Quatarari"),
	mustBe("Via Marco Polo 7a", "Quatarari"),
	mustBe("Via Santa Caterina da Siena 99", "Quatarari"),

	// San Giovanni
	mustBe("Via Pietro Micca", "San Giovanni"),
	mustBe("Via Sara Librando 20", "San Giovanni"),
	mustBe("Eurospin", "San Giovanni"),
	mustBe("Piazza Paolino Arnesano", "San Giovanni"),
	mustBe("Via Nazario Sauro", "San Giovanni"),
	mustBe("Via Firenze", "San Giovanni"),
	mustBe("Via Bologna", "San Giovanni"),
	mustBe("Circolo Tennis Grande Slam", "San Giovanni"),
	mustBe("Via Giorgione", "San Giovanni"),
	mustBe("Scuola Maurizio Arnesano", "San Giovanni"),

	// Saraceni
	mustBe("Via Carso", "Saraceni"),
	mustBe("Via Gagliardina", "Saraceni"),
	mustBe("Via Garibaldi", "Saraceni"),
	mustBe("Via Don Donato Franco", "Saraceni"),
	mustBe("Panetteria La Genuina", "Saraceni"),
	mustBe("Via Aldo Moro", "Saraceni"),
	mustBe("Via Alcide de Gasperi", "Saraceni"),
	mustBe("Fantasy Cartolibreria", "Saraceni"),
	mustBe("Burlesq Fashionable", "Saraceni"),
	mustBe("Via Salvo D'Acquisto 3", "Saraceni"),
}

func setupService(t *testing.T) *AppService {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	mapsApiKey := os.Getenv("MAPS_API_KEY")
	if mapsApiKey == "" {
		logger.Errorln("MAPS_API_KEY environment variable not set")
		t.Fail()
	}

	cfg := &AppConfig{
		MapsApiKey:   mapsApiKey,
		ServeAddress: "localhost:8080",
	}

	svc := NewService(logger, cfg)

	// Load the GeoJSON map
	if err := svc.LoadGeoJSONMapFile("map.geojson"); err != nil {
		t.Fatalf("Failed to load GeoJSON map: %v", err)
	}

	return svc
}

func TestAddressToRioneMapping(t *testing.T) {
	svc := setupService(t)

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Geocode the address
			geocode, err := svc.geocode(tc.address)
			if err != nil {
				t.Errorf("Geocoding failed: %v", err)
				return
			}

			if geocode.Status != GeocodeStatusOK || len(geocode.Results) == 0 {
				t.Errorf("No geocoding results for: %s", tc.address)
				return
			}

			// Verify address is in Carmiano
			result := geocode.Results[0]
			if !isInCarmiano(result) {
				t.Errorf("Address not in Carmiano: %s", tc.address)
				return
			}

			// Parse and check rione
			parsed := svc.parseGeocodeResult(&result)

			if parsed.Rione != tc.expectedRione {
				t.Errorf("Rione mismatch for %s: expected %q, got %q",
					tc.address, tc.expectedRione, parsed.Rione)
			}

			t.Logf("✓ %s -> %s (%.6f, %.6f)",
				parsed.Street, parsed.Rione, parsed.Latitude, parsed.Longitude)
		})
	}
}
