package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type AppConfig struct {
	Port int
	ServeAddress string
	MapsApiKey string
}

func GetConfig(logger *logrus.Logger) (*AppConfig, error) {
	godotenv.Load()

	// Port

	port := os.Getenv("PORT")
	if (len(port) == 0) {
		logger.Warnln("PORT env var not set; defaulting to 8080")
		port = "8080"
	}

	intPort, err := strconv.Atoi(port)
	if (err != nil) {
		logger.Warnf("invalid port %s provided; defaulting to 8080", port)
		intPort = 8080
	}

	// Maps API Key
	apiKey := os.Getenv("MAPS_API_KEY")
	if (len(apiKey) == 0) {
		return nil, fmt.Errorf("MAPS_API_KEY not provided")
	}

	return &AppConfig{
		Port: intPort,
		ServeAddress: fmt.Sprintf(":%d", intPort),
		MapsApiKey: apiKey,
	}, nil
}