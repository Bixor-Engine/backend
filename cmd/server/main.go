package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/Bixor-Engine/backend/docs"
	"github.com/Bixor-Engine/backend/internal/routes"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type App struct {
	DB *sql.DB
}

// @title Bixor Trading Engine API
// @version 1.0.0
// @description High-performance trading backend API for cryptocurrency exchange operations
// @contact.email info@bixor.io
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
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

	// Setup routes
	router := routes.SetupRoutes(app.DB)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Go service starting on port %s", port)
	log.Printf("Swagger UI available at: http://localhost:%s/swagger/index.html", port)
	log.Printf("API Documentation at: http://localhost:%s/docs", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func (app *App) initDB() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required. Please create a .env file (copy from .env.example)")
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
