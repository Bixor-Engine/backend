package handlers

import (
	"net/http"
	"strings"

	"github.com/Bixor-Engine/backend/internal/models"
	"github.com/gin-gonic/gin"
)

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update user's personal information
// @Tags Authorization
// @Accept json
// @Produce json
// @Security BackendSecret
// @Security BearerAuth
// @Param request body models.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} models.UserResponse "User profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/profile/update [post]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	token, err := h.valToken(c)
	if err != nil {
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_failed",
			"message": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Update user in database
	query := `
		UPDATE users 
		SET first_name = $1, last_name = $2, phone_number = $3, 
			address = $4, city = $5, country = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING id, first_name, last_name, username, email, email_status, 
				  phone_number, phone_status, address, city, country, 
				  role, status, kyc_status, twofa_enabled, last_login_at, 
				  language, timezone, global_balance, created_at, updated_at
	`

	var user models.UserResponse
	err = h.DB.QueryRow(query,
		req.FirstName, req.LastName, req.PhoneNumber,
		req.Address, req.City, req.Country,
		token.UserID,
	).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.EmailStatus,
		&user.PhoneNumber, &user.PhoneStatus, &user.Address, &user.City, &user.Country,
		&user.Role, &user.Status, &user.KYCStatus, &user.TwoFAEnabled, &user.LastLoginAt,
		&user.Language, &user.Timezone, &user.GlobalBalance, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to update profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

// valToken is a helper to validate token and return claims
func (h *AuthHandler) valToken(c *gin.Context) (*models.JWTClaims, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "missing_authorization",
			"message": "Authorization header is required",
		})
		return nil, models.ErrInvalidToken
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := models.ValidateAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_token",
			"message": "Invalid or expired access token",
		})
		return nil, err
	}
	return claims, nil
}
