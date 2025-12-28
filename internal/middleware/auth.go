package middleware

import (
	"net/http"
	"os"

	"github.com/Bixor-Engine/backend/internal/models"
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

// UserTokenMiddleware validates the JWT access token for personal API routes
func UserTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing_token", "message": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if it starts with "Bearer "
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token_format", "message": "Authorization header must be Bearer token"})
			c.Abort()
			return
		}

		tokenString := authHeader[7:]
		claims, err := models.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token", "message": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user ID and other claims in context
		c.Set("userID", claims.UserID.String())
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}
