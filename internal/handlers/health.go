package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	DB *sql.DB
}

type ServiceHealth struct {
	Service   string            `json:"service"`
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Details   map[string]string `json:"details"`
}

type ServiceStatus struct {
	Service   string    `json:"service"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		DB: db,
	}
}

func (h *HealthHandler) getAPIVersion() string {
	version := os.Getenv("API_VERSION")
	if version == "" {
		return "1.0.0"
	}
	return version
}

func (h *HealthHandler) checkDatabaseHealth() (bool, string, time.Duration) {
	start := time.Now()

	if h.DB == nil {
		return false, "Database connection not initialized", time.Since(start)
	}

	if err := h.DB.Ping(); err != nil {
		log.Printf("Database health check failed: %v", err)
		return false, "Database connection failed", time.Since(start)
	}

	return true, "Database connection successful", time.Since(start)
}

// HealthCheck godoc
// @Summary Health check for all services
// @Description Check the health status of all services including API and database
// @Tags Monitoring
// @Accept json
// @Produce json
// @Success 200 {array} ServiceHealth "All services are healthy"
// @Failure 503 {array} ServiceHealth "One or more services are unhealthy"
// @Router /api/v1/health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	timestamp := time.Now()
	version := h.getAPIVersion()

	// Check API service
	apiService := ServiceHealth{
		Service:   "api-service",
		Status:    "healthy",
		Timestamp: timestamp,
		Details: map[string]string{
			"version": version,
		},
	}

	// Check database service - real check
	isHealthy, message, responseTime := h.checkDatabaseHealth()
	databaseService := ServiceHealth{
		Service:   "database",
		Status:    "healthy",
		Timestamp: timestamp,
		Details: map[string]string{
			"version":       version,
			"response_time": responseTime.String(),
		},
	}

	if !isHealthy {
		databaseService.Status = "unhealthy"
		databaseService.Details["error"] = message
	}

	services := []ServiceHealth{apiService, databaseService}

	// If any service is unhealthy, return 503
	for _, service := range services {
		if service.Status == "unhealthy" {
			c.JSON(http.StatusServiceUnavailable, services)
			return
		}
	}

	c.JSON(http.StatusOK, services)
}

// GetStatus godoc
// @Summary Get status of all services
// @Description Get the current operational status of all services
// @Tags Monitoring
// @Accept json
// @Produce json
// @Success 200 {array} ServiceStatus "All services are active"
// @Failure 503 {array} ServiceStatus "One or more services are inactive"
// @Router /api/v1/status [get]
func (h *HealthHandler) GetStatus(c *gin.Context) {
	timestamp := time.Now()
	version := h.getAPIVersion()

	// Check API service status
	apiStatus := ServiceStatus{
		Service:   "api-service",
		Status:    "active",
		Message:   "API Service is running",
		Timestamp: timestamp,
		Version:   version,
	}

	// Check database service status - real check
	isActive, message, _ := h.checkDatabaseHealth()
	databaseStatus := ServiceStatus{
		Service:   "database",
		Status:    "active",
		Message:   message,
		Timestamp: timestamp,
		Version:   version,
	}

	if !isActive {
		databaseStatus.Status = "inactive"
		databaseStatus.Message = message
	}

	services := []ServiceStatus{apiStatus, databaseStatus}

	// If any service is inactive, return 503
	for _, service := range services {
		if service.Status == "inactive" {
			c.JSON(http.StatusServiceUnavailable, services)
			return
		}
	}

	c.JSON(http.StatusOK, services)
}
