package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Bixor-Engine/backend/internal/models"
	"github.com/gin-gonic/gin"
)

// UpdateSettings godoc
// @Summary Update user settings
// @Description Update user's preferences (language, timezone)
// @Tags Authorization
// @Accept json
// @Produce json
// @Security BackendSecret
// @Security BearerAuth
// @Param request body models.UpdateSettingsRequest true "Settings update data"
// @Success 200 {object} models.UserResponse "User settings updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/settings/update [post]
func (h *AuthHandler) UpdateSettings(c *gin.Context) {
	token, err := h.valToken(c)
	if err != nil {
		return
	}

	var req models.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_failed",
			"message": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	query := `
		UPDATE users 
		SET language = $1, timezone = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, first_name, last_name, username, email, email_status, 
				  phone_number, phone_status, address, city, country, 
				  role, status, kyc_status, twofa_enabled, last_login_at, 
				  language, timezone, global_balance, created_at, updated_at
	`

	var user models.UserResponse
	err = h.DB.QueryRow(query, req.Language, req.Timezone, token.UserID).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.EmailStatus,
		&user.PhoneNumber, &user.PhoneStatus, &user.Address, &user.City, &user.Country,
		&user.Role, &user.Status, &user.KYCStatus, &user.TwoFAEnabled, &user.LastLoginAt,
		&user.Language, &user.Timezone, &user.GlobalBalance, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to update settings",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Settings updated successfully",
		"user":    user,
	})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change user's password with verification of current password
// @Tags Authorization
// @Accept json
// @Produce json
// @Security BackendSecret
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "Password change data"
// @Success 200 {object} map[string]string "Password changed successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid current password"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/security/password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	token, err := h.valToken(c)
	if err != nil {
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_failed",
			"message": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// 1. Get current password hash
	var currentHash string
	err = h.DB.QueryRow("SELECT password FROM users WHERE id = $1", token.UserID).Scan(&currentHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to retrieve user data",
		})
		return
	}

	// 2. Verify current password
	match, err := models.VerifyPassword(req.CurrentPassword, currentHash)
	if err != nil || !match {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_password",
			"message": "Current password is incorrect",
		})
		return
	}

	// 3. Hash new password
	newHash, err := models.HashPassword(req.NewPassword, models.DefaultArgonParams())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "hashing_error",
			"message": "Failed to process new password",
		})
		return
	}

	// 4. Update password
	_, err = h.DB.Exec("UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2", newHash, token.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to update password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// ToggleTwoFA godoc
// @Summary Toggle 2FA status
// @Description Enable or disable 2FA. Requires OTP verification code.
// @Tags Authorization
// @Accept json
// @Produce json
// @Security BackendSecret
// @Security BearerAuth
// @Param request body models.ToggleTwoFARequest true "2FA toggle data"
// @Success 200 {object} models.UserResponse "2FA status updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid OTP or validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/security/2fa [post]
func (h *AuthHandler) ToggleTwoFA(c *gin.Context) {
	token, err := h.valToken(c)
	if err != nil {
		return
	}

	var req models.ToggleTwoFARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_failed",
			"message": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// 1. Verify OTP Code
	// Get latest unused 2FA OTP for this user
	var otp models.OTP
	err = h.DB.QueryRow(`
		SELECT id, code, expires_at
		FROM otps
		WHERE user_id = $1 AND type = '2fa' AND used = FALSE
		ORDER BY created_at DESC
		LIMIT 1
	`, token.UserID).Scan(&otp.ID, &otp.Code, &otp.ExpiresAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "otp_required",
			"message": "Please request a 2FA verification code first",
		})
		return
	}

	if time.Now().After(otp.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "otp_expired",
			"message": "OTP code has expired",
		})
		return
	}

	if otp.Code != req.Code {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_otp",
			"message": "Invalid OTP code",
		})
		return
	}

	// 2. Mark OTP as used
	h.DB.Exec("UPDATE otps SET used = TRUE, updated_at = NOW() WHERE id = $1", otp.ID)

	// 3. Update User 2FA Status
	query := `
		UPDATE users 
		SET twofa_enabled = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, first_name, last_name, username, email, email_status, 
				  phone_number, phone_status, address, city, country, 
				  role, status, kyc_status, twofa_enabled, last_login_at, 
				  language, timezone, global_balance, created_at, updated_at
	`

	var user models.UserResponse
	err = h.DB.QueryRow(query, req.Enable, token.UserID).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.EmailStatus,
		&user.PhoneNumber, &user.PhoneStatus, &user.Address, &user.City, &user.Country,
		&user.Role, &user.Status, &user.KYCStatus, &user.TwoFAEnabled, &user.LastLoginAt,
		&user.Language, &user.Timezone, &user.GlobalBalance, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to update 2FA status",
		})
		return
	}

	message := "2FA enabled successfully"
	if !req.Enable {
		message = "2FA disabled successfully"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"user":    user,
	})
}
