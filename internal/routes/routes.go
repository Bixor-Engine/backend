package routes

import (
	"database/sql"
	"net/http"

	"github.com/Bixor-Engine/backend/internal/handlers"
	"github.com/Bixor-Engine/backend/internal/middleware"
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
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Backend-Secret, X-API-Secret")

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
		// ============================================
		// PUBLIC ROUTES - No authentication required
		// ============================================
		public := v1.Group("")
		public.Use(middleware.PublicMiddleware())
		{
			// Health and status endpoints (public for monitoring)
			public.GET("/health", healthHandler.HealthCheck)
			public.GET("/status", healthHandler.GetStatus)
			public.GET("/info", apiHandler.APIInfo)
		}

		// ============================================
		// PROTECTED BY BACKEND SECRET - Frontend requests
		// ============================================
		protected := v1.Group("")
		protected.Use(middleware.BackendSecretMiddleware())
		{
			// Authentication endpoints (frontend uses backend secret)
			auth := protected.Group("/auth")
			{
				auth.POST("/register", authHandler.Register)
				auth.POST("/login", authHandler.Login)
				auth.POST("/refresh", authHandler.RefreshToken)
				auth.GET("/me", authHandler.GetCurrentUser)

				// OTP endpoints (require backend secret + JWT token)
				otp := auth.Group("/otp")
				{
					otp.POST("/request", authHandler.RequestOTP)
					otp.POST("/verify", authHandler.VerifyOTP)
				}

				// Logout endpoint (requires JWT token)
				auth.POST("/logout", authHandler.Logout)

				// Future auth endpoints will be added here
				// auth.POST("/forgot-password", authHandler.ForgotPassword)
			}

			// Future protected routes (frontend access)
			// protected.GET("/users", userHandler.GetUsers)
			// protected.POST("/orders", orderHandler.CreateOrder)
			// protected.GET("/markets", marketHandler.GetMarkets)
		}

		// ============================================
		// PERSONAL API ROUTES - User token based (Future)
		// ============================================
		personal := v1.Group("/personal")
		personal.Use(middleware.UserTokenMiddleware())
		{
			// Future personal API endpoints
			// personal.GET("/trades", personalHandler.GetTrades)
			// personal.POST("/orders", personalHandler.CreateOrder)
			// personal.GET("/balance", personalHandler.GetBalance)
		}
	}

	return router
}
