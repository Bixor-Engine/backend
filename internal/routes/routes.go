package routes

import (
	"database/sql"
	"net/http"

	"github.com/Bixor-Engine/backend/internal/handlers"
	"github.com/Bixor-Engine/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(db *sql.DB) *gin.Engine {
	router := gin.Default()

	// Initialize handlers
	apiHandler := handlers.NewAPIHandler()
	healthHandler := handlers.NewHealthHandler(db)
	authHandler := handlers.NewAuthHandler(db)
	currencyHandler := handlers.NewCurrencyHandler(db)
	walletHandler := handlers.NewWalletHandler(db)
	transactionHandler := handlers.NewTransactionHandler(db)

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

	// Load HTML templates
	router.LoadHTMLGlob("internal/templates/*")

	// API Documentation landing page
	router.GET("/docs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "docs.html", gin.H{})
	})
	router.GET("/docs/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs")
	})
	router.GET("/docs/index.html", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs")
	})

	// Swagger spec endpoints (JSON)
	swaggerHandler := handlers.NewSwaggerHandler()
	router.GET("/docs/public.json", swaggerHandler.GetPublicSwaggerSpec)
	router.GET("/docs/private.json", swaggerHandler.GetPrivateSwaggerSpec)
	router.GET("/docs/personal.json", swaggerHandler.GetPersonalSwaggerSpec)

	// Swagger UI endpoints
	router.GET("/docs/public", func(c *gin.Context) {
		c.HTML(http.StatusOK, "swagger.html", gin.H{
			"Title":   "Public API",
			"specURL": "/docs/public.json",
		})
	})
	router.GET("/docs/private", func(c *gin.Context) {
		c.HTML(http.StatusOK, "swagger.html", gin.H{
			"Title":   "Private API",
			"specURL": "/docs/private.json",
		})
	})
	router.GET("/docs/personal", func(c *gin.Context) {
		c.HTML(http.StatusOK, "swagger.html", gin.H{
			"Title":   "Personal API",
			"specURL": "/docs/personal.json",
		})
	})

	// Redirect /swagger to /docs for documentation landing page
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs")
	})
	router.GET("/swagger/*any", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs")
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

			// Currency endpoints (public - no authentication required)
			currency := public.Group("/currency")
			{
				currency.GET("", currencyHandler.GetCoins)
				currency.GET("/:ticker", currencyHandler.GetCoinByTicker)
			}
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

				// Profile management
				auth.POST("/profile/update", authHandler.UpdateProfile)
				auth.POST("/settings/update", authHandler.UpdateSettings)

				// Security management
				security := auth.Group("/security")
				{
					security.POST("/password", authHandler.ChangePassword)
					security.POST("/2fa", authHandler.ToggleTwoFA)
				}

				// Future auth endpoints will be added here
				// auth.POST("/forgot-password", authHandler.ForgotPassword)
			}

			// Authenticated User Routes (require both Secret + JWT)
			userRoutes := protected.Group("")
			userRoutes.Use(middleware.UserTokenMiddleware())
			{
				userRoutes.GET("/wallets", walletHandler.GetWallets)
				userRoutes.GET("/transactions", transactionHandler.GetTransactions)
			}
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
