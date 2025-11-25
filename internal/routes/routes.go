package routes

import (
	"database/sql"
	"net/http"

	"github.com/Bixor-Engine/backend/internal/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(db *sql.DB) *gin.Engine {
	router := gin.Default()

	// Initialize handlers
	apiHandler := handlers.NewAPIHandler()
	healthHandler := handlers.NewHealthHandler(db)
	authHandler := handlers.NewAuthHandler(db)

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Landing page
	router.GET("/", apiHandler.LandingPage)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Alternative docs route
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// API version 1 routes
	v1 := router.Group("/api/v1")
	{
		// Health and status endpoints
		v1.GET("/health", healthHandler.HealthCheck)
		v1.GET("/status", healthHandler.GetStatus)
		v1.GET("/info", apiHandler.APIInfo)

		// Authentication endpoints
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.GET("/me", authHandler.GetCurrentUser)

			// OTP endpoints (require authentication)
			otp := auth.Group("/otp")
			{
				otp.POST("/request", authHandler.RequestOTP)
				otp.POST("/verify", authHandler.VerifyOTP)
			}

			// Future auth endpoints will be added here
			// auth.POST("/logout", authHandler.Logout)
			// auth.POST("/forgot-password", authHandler.ForgotPassword)
		}

		// Future routes will be added here
		// v1.GET("/users", userHandler.GetUsers)
		// v1.POST("/orders", orderHandler.CreateOrder)
		// v1.GET("/markets", marketHandler.GetMarkets)
	}

	return router
}
