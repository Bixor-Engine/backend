package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type APIHandler struct{}

func NewAPIHandler() *APIHandler {
	return &APIHandler{}
}

// LandingPage godoc
// @Summary Landing page
// @Description Get basic information about the Bixor Trading Engine
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Landing page message"
// @Router / [get]
func (h *APIHandler) LandingPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Powered by Bixor-Engine Go Trading Engine - info@bixor.io",
	})
}

// APIInfo godoc
// @Summary API information
// @Description Get detailed information about all available services and endpoints
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "API information with services and endpoints"
// @Router /api/v1/info [get]
func (h *APIHandler) APIInfo(c *gin.Context) {
	version := os.Getenv("API_VERSION")
	if version == "" {
		version = "1.0.0"
	}

	c.JSON(http.StatusOK, gin.H{
		"api_version": "v1",
		"services": []gin.H{
			{
				"service":     "api-service",
				"version":     version,
				"description": "REST API service for trading operations",
				"endpoints": gin.H{
					"health": "/api/v1/health",
					"status": "/api/v1/status",
					"info":   "/api/v1/info",
				},
			},
			{
				"service":     "database",
				"version":     version,
				"description": "Database service",
				"endpoints": gin.H{
					"check_via": "/api/v1/health",
				},
			},
		},
		"documentation": "https://docs.bixor.io",
		"support":       "info@bixor.io",
	})
}
