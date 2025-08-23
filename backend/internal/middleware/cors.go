package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/eikuma/stockle/backend/internal/config"
)

func CORS(cfg *config.Config) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowMethods:     cfg.Server.CORS.AllowedMethods,
		AllowHeaders:     cfg.Server.CORS.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	
	// Development mode: allow specific origins for credentials support
	if cfg.Server.Environment == "development" {
		corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
	} else {
		corsConfig.AllowOrigins = cfg.Server.CORS.AllowedOrigins
	}
	
	return cors.New(corsConfig)
}