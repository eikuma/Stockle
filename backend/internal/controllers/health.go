package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/eikuma/stockle/backend/internal/database"
	"github.com/eikuma/stockle/backend/internal/config"
)

type HealthController struct {
	config *config.Config
}

func NewHealthController(cfg *config.Config) *HealthController {
	return &HealthController{
		config: cfg,
	}
}

type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version"`
	Checks    map[string]HealthCheck `json:"checks"`
}

type HealthCheck struct {
	Status   string        `json:"status"`
	Duration time.Duration `json:"duration"`
	Error    string        `json:"error,omitempty"`
}

func (h *HealthController) Health(c *gin.Context) {
	start := time.Now()
	
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: start,
		Version:   "1.0.0",
		Checks:    make(map[string]HealthCheck),
	}
	
	// Database health check
	dbStart := time.Now()
	dbCheck := HealthCheck{
		Status:   "healthy",
		Duration: 0,
	}
	
	if err := database.IsHealthy(); err != nil {
		dbCheck.Status = "unhealthy"
		dbCheck.Error = err.Error()
		response.Status = "unhealthy"
	}
	
	dbCheck.Duration = time.Since(dbStart)
	response.Checks["database"] = dbCheck
	
	// API health check
	apiCheck := HealthCheck{
		Status:   "healthy",
		Duration: time.Since(start),
	}
	response.Checks["api"] = apiCheck
	
	// Configuration health check
	configCheck := HealthCheck{
		Status:   "healthy",
		Duration: time.Nanosecond, // Config check is instant
	}
	
	if h.config == nil {
		configCheck.Status = "unhealthy"
		configCheck.Error = "configuration not loaded"
		response.Status = "unhealthy"
	}
	
	response.Checks["config"] = configCheck
	
	// Set HTTP status based on overall health
	statusCode := http.StatusOK
	if response.Status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}
	
	c.JSON(statusCode, response)
}

func (h *HealthController) Ready(c *gin.Context) {
	// Check if all required components are ready
	checks := []func() error{
		database.IsHealthy,
	}
	
	for _, check := range checks {
		if err := check(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  err.Error(),
			})
			return
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

func (h *HealthController) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}