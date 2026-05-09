package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	logger := GetLogger()
	logger.Infoln("Starting App")
	cfg, err := GetConfig(logger)

	if err != nil {
		logger.Errorln("Fatal error during conf loading;")
		logger.Errorln(err)
		logger.Error("Stopping execution")

		logger.Exit(1)
	}

	svc := NewService(logger, cfg)

	// Initialize map service with GeoJSON file
	if err := svc.LoadGeoJSONMapFile("map.geojson"); err != nil {
		logger.Errorf("could not load map file: %s", err.Error())
		logger.Exit(1)
	}

	r := gin.Default()

	r.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"https://torneodeirioni.it", "https://giovaniinformazione.it", "http://localhost:4321"},
			AllowMethods:     []string{"GET", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
	)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})
	r.GET("/get", svc.TrovaRionePerIndirizzo)

	r.Run(cfg.ServeAddress)
}
