package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// BackendSecretMiddleware validates the backend secret from request header
// This is used to protect API routes that should only be accessible from the frontend
func BackendSecretMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		backendSecret := os.Getenv("BACKEND_SECRET")
		if backendSecret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "backend_secret_not_configured",
				"message": "Backend secret is not configured. Please set BACKEND_SECRET in environment variables.",
			})
			c.Abort()
			return
		}

		// Get secret from header (X-Backend-Secret or X-API-Secret)
		secret := c.GetHeader("X-Backend-Secret")
		if secret == "" {
			secret = c.GetHeader("X-API-Secret")
		}

		if secret == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "missing_backend_secret",
				"message": "Backend secret is required. Include X-Backend-Secret header in your request.",
			})
			c.Abort()
			return
		}

		if secret != backendSecret {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_backend_secret",
				"message": "Invalid backend secret",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// PublicMiddleware marks routes as public (no authentication required)
// This is just a passthrough, but useful for documentation
func PublicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// UserTokenMiddleware will be used in the future for personal API tokens
// For now, it's a placeholder that can be extended later
func UserTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement user token validation when personal API is implemented
		// For now, this is a placeholder
		c.Next()
	}
}
