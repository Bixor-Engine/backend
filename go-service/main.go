package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type App struct {
	DB *sql.DB
	Router *gin.Engine
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Timestamp time.Time         `json:"timestamp"`
	Database  DatabaseHealth    `json:"database"`
	Details   map[string]string `json:"details"`
}

type DatabaseHealth struct {
	Status      string        `json:"status"`
	ResponseTime time.Duration `json:"response_time_ms"`
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	app := &App{}
	
	// Initialize database
	if err := app.initDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer app.DB.Close()

	// Initialize router
	app.initRoutes()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Go service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, app.Router))
}

func (app *App) initDB() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://bixor_user:bixor_pass@localhost:5432/bixor?sslmode=disable"
	}

	var err error
	app.DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return err
	}

	// Test connection
	if err = app.DB.Ping(); err != nil {
		return err
	}

	log.Println("Database connection established")
	return nil
}

func (app *App) initRoutes() {
	app.Router = gin.Default()

	// CORS middleware
	app.Router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Health check endpoint
	app.Router.GET("/health", app.healthCheck)

	// API version 1 routes
	v1 := app.Router.Group("/api/v1")
	{
		v1.GET("/status", app.getStatus)
	}
}

func (app *App) healthCheck(c *gin.Context) {
	start := time.Now()
	
	// Check database health
	dbHealth := DatabaseHealth{
		Status: "healthy",
	}
	
	if err := app.DB.Ping(); err != nil {
		dbHealth.Status = "unhealthy"
		log.Printf("Database health check failed: %v", err)
	}
	
	dbHealth.ResponseTime = time.Since(start)

	response := HealthResponse{
		Status:    "healthy",
		Service:   "bixor-go-service",
		Timestamp: time.Now(),
		Database:  dbHealth,
		Details: map[string]string{
			"version": "1.0.0",
			"uptime":  time.Since(start).String(),
		},
	}

	// If database is unhealthy, mark overall status as unhealthy
	if dbHealth.Status == "unhealthy" {
		response.Status = "unhealthy"
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (app *App) getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Bixor Go Service is running",
		"timestamp": time.Now(),
		"service": "bixor-go-service",
		"version": "1.0.0",
	})
} 