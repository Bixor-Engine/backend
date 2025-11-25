package handlers

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/Bixor-Engine/backend/internal/models"
	"github.com/Bixor-Engine/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	DB *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{
		DB: db,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User registration data"
// @Success 201 {object} models.UserResponse "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - validation errors"
// @Failure 409 {object} map[string]interface{} "Conflict - user already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_failed",
			"message": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Check if username or email already exists
	if exists, err := h.checkUserExists(req.Username, req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to check user existence",
		})
		return
	} else if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "user_exists",
			"message": "Username or email already exists",
		})
		return
	}

	// Hash the password using Argon2i
	hashedPassword, err := models.HashPassword(req.Password, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "password_hash_error",
			"message": "Failed to hash password",
		})
		return
	}

	// Parse referred_by UUID if provided
	var referredByUUID *uuid.UUID
	if req.ReferredBy != nil && *req.ReferredBy != "" {
		parsedUUID, err := uuid.Parse(*req.ReferredBy)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_referral",
				"message": "Invalid referral UUID format",
			})
			return
		}
		referredByUUID = &parsedUUID
	}

	// Set default values for optional fields
	language := "en"
	if req.Language != nil && *req.Language != "" {
		language = *req.Language
	}

	timezone := "UTC"
	if req.Timezone != nil && *req.Timezone != "" {
		timezone = *req.Timezone
	}

	// Create user in database
	user := &models.User{
		ID:           uuid.New(),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Username:     strings.ToLower(strings.TrimSpace(req.Username)),
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		Password:     hashedPassword,
		EmailStatus:  false, // Email verification required
		PhoneNumber:  req.PhoneNumber,
		PhoneStatus:  false, // Phone verification required
		ReferredBy:   referredByUUID,
		Address:      req.Address,
		City:         req.City,
		Country:      req.Country,
		Role:         "user",          // Default role
		Status:       "pending",       // Default status
		KYCStatus:    "not_submitted", // Default KYC status
		TwoFAEnabled: false,           // 2FA disabled by default
		Language:     language,
		Timezone:     timezone,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Insert user into database
	if err := h.createUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to create user",
		})
		return
	}

	// Create response (excluding sensitive data)
	userResponse := models.UserResponse{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     user.Username,
		Email:        user.Email,
		EmailStatus:  user.EmailStatus,
		PhoneNumber:  user.PhoneNumber,
		PhoneStatus:  user.PhoneStatus,
		Address:      user.Address,
		City:         user.City,
		Country:      user.Country,
		Role:         user.Role,
		Status:       user.Status,
		KYCStatus:    user.KYCStatus,
		TwoFAEnabled: user.TwoFAEnabled,
		Language:     user.Language,
		Timezone:     user.Timezone,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    userResponse,
	})
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password, returns JWT tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "User login credentials"
// @Success 200 {object} models.LoginResponse "Login successful"
// @Failure 400 {object} map[string]interface{} "Bad request - validation errors"
// @Failure 401 {object} map[string]interface{} "Unauthorized - invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_failed",
			"message": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Normalize email (lowercase and trim) to match registration behavior
	normalizedEmail := strings.ToLower(strings.TrimSpace(req.Email))

	// Find user by email
	user, err := h.getUserByEmail(normalizedEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_credentials",
				"message": "Invalid email or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to authenticate user",
		})
		return
	}

	// Check if user account is allowed to login
	// Allow "active" and "pending" status (pending = newly registered users)
	// Block suspended, banned, locked, frozen, inactive accounts
	if user.Status != "active" && user.Status != "pending" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "account_inactive",
			"message": "Account is not active. Please contact support.",
		})
		return
	}

	// Verify password
	isValidPassword, err := models.VerifyPassword(req.Password, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "password_verification_error",
			"message": "Failed to verify password",
		})
		return
	}

	if !isValidPassword {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_credentials",
			"message": "Invalid email or password",
		})
		return
	}

	// Generate JWT tokens
	tokens, err := models.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "token_generation_error",
			"message": "Failed to generate authentication tokens",
		})
		return
	}

	// Update last login information
	if err := h.updateLastLogin(user.ID, c.ClientIP()); err != nil {
		// Log error but don't fail the login
		// In production, you might want to log this properly
	}

	// Create user response (excluding sensitive data)
	userResponse := models.UserResponse{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     user.Username,
		Email:        user.Email,
		EmailStatus:  user.EmailStatus,
		PhoneNumber:  user.PhoneNumber,
		PhoneStatus:  user.PhoneStatus,
		Address:      user.Address,
		City:         user.City,
		Country:      user.Country,
		Role:         user.Role,
		Status:       user.Status,
		KYCStatus:    user.KYCStatus,
		TwoFAEnabled: user.TwoFAEnabled,
		LastLoginAt:  user.LastLoginAt,
		Language:     user.Language,
		Timezone:     user.Timezone,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	// Create login response with redirect info
	loginResponse := models.LoginResponse{
		Message:        "Login successful",
		User:           userResponse,
		Tokens:         *tokens,
		RequiresVerify: !user.EmailStatus, // True if email is not verified
		RedirectTo:     "",                // Will be set if verification needed
	}

	// Set redirect if email verification is required
	if !user.EmailStatus {
		loginResponse.RedirectTo = "/verify-email"
	}

	c.JSON(http.StatusOK, loginResponse)
}

// RefreshToken godoc
// @Summary Refresh JWT tokens
// @Description Generate new JWT tokens using a valid refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param refresh body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} models.JWTTokens "Tokens refreshed successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - validation errors"
// @Failure 401 {object} map[string]interface{} "Unauthorized - invalid refresh token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_failed",
			"message": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Validate refresh token and get claims
	claims, err := models.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_refresh_token",
			"message": "Invalid or expired refresh token",
		})
		return
	}

	// Get user from database
	user, err := h.getUserByID(claims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "user_not_found",
				"message": "User associated with token not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to retrieve user",
		})
		return
	}

	// Check if user is still active
	if user.Status != "active" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "account_inactive",
			"message": "Account is no longer active",
		})
		return
	}

	// Generate new tokens
	newTokens, err := models.RefreshTokens(req.RefreshToken, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "token_generation_error",
			"message": "Failed to generate new tokens",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tokens refreshed successfully",
		"tokens":  *newTokens,
	})
}

// GetCurrentUser godoc
// @Summary Get current authenticated user
// @Description Get current user information based on JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse "Current user data"
// @Failure 401 {object} map[string]interface{} "Unauthorized - invalid or missing token"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// Get authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "missing_authorization",
			"message": "Authorization header is required",
		})
		return
	}

	// Check if the header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_authorization_format",
			"message": "Authorization header must be in format 'Bearer <token>'",
		})
		return
	}

	// Extract the token
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Validate token and get claims
	claims, err := models.ValidateAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_token",
			"message": "Invalid or expired access token",
		})
		return
	}

	// Get user from database
	user, err := h.getUserByID(claims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "user_not_found",
				"message": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to retrieve user data",
		})
		return
	}

	// Check if user account is allowed
	// Allow "active" and "pending" status (pending = newly registered users)
	// Block suspended, banned, locked, frozen, inactive accounts
	if user.Status != "active" && user.Status != "pending" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "account_inactive",
			"message": "Account is no longer active",
		})
		return
	}

	// Create response
	userResponse := models.UserResponse{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     user.Username,
		Email:        user.Email,
		EmailStatus:  user.EmailStatus,
		PhoneNumber:  user.PhoneNumber,
		PhoneStatus:  user.PhoneStatus,
		Address:      user.Address,
		City:         user.City,
		Country:      user.Country,
		Role:         user.Role,
		Status:       user.Status,
		KYCStatus:    user.KYCStatus,
		TwoFAEnabled: user.TwoFAEnabled,
		LastLoginAt:  user.LastLoginAt,
		Language:     user.Language,
		Timezone:     user.Timezone,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User data retrieved successfully",
		"user":    userResponse,
	})
}

// checkUserExists checks if a user with the given username or email already exists
func (h *AuthHandler) checkUserExists(username, email string) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM users 
		WHERE LOWER(username) = LOWER($1) OR LOWER(email) = LOWER($2)
		AND deleted_at IS NULL
	`

	err := h.DB.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// createUser inserts a new user into the database
func (h *AuthHandler) createUser(user *models.User) error {
	query := `
		INSERT INTO users (
			id, first_name, last_name, username, email, password,
			email_status, phone_number, phone_status, referred_by,
			address, city, country, role, status, kyc_status,
			twofa_enabled, language, timezone, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		)
	`

	_, err := h.DB.Exec(query,
		user.ID, user.FirstName, user.LastName, user.Username, user.Email, user.Password,
		user.EmailStatus, user.PhoneNumber, user.PhoneStatus, user.ReferredBy,
		user.Address, user.City, user.Country, user.Role, user.Status, user.KYCStatus,
		user.TwoFAEnabled, user.Language, user.Timezone, user.CreatedAt, user.UpdatedAt,
	)

	return err
}

// getUserByEmail retrieves a user by email address
func (h *AuthHandler) getUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, first_name, last_name, username, email, password,
			   email_status, phone_number, phone_status, referred_by,
			   address, city, country, role, status, kyc_status,
			   twofa_enabled, last_login_at, last_login_ip, device_info,
			   language, timezone, created_at, updated_at, deleted_at
		FROM users 
		WHERE LOWER(email) = LOWER($1) AND deleted_at IS NULL
	`

	err := h.DB.QueryRow(query, email).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password,
		&user.EmailStatus, &user.PhoneNumber, &user.PhoneStatus, &user.ReferredBy,
		&user.Address, &user.City, &user.Country, &user.Role, &user.Status, &user.KYCStatus,
		&user.TwoFAEnabled, &user.LastLoginAt, &user.LastLoginIP, &user.DeviceInfo,
		&user.Language, &user.Timezone, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// getUserByID retrieves a user by ID
func (h *AuthHandler) getUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, first_name, last_name, username, email, password,
			   email_status, phone_number, phone_status, referred_by,
			   address, city, country, role, status, kyc_status,
			   twofa_enabled, last_login_at, last_login_ip, device_info,
			   language, timezone, created_at, updated_at, deleted_at
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := h.DB.QueryRow(query, userID).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password,
		&user.EmailStatus, &user.PhoneNumber, &user.PhoneStatus, &user.ReferredBy,
		&user.Address, &user.City, &user.Country, &user.Role, &user.Status, &user.KYCStatus,
		&user.TwoFAEnabled, &user.LastLoginAt, &user.LastLoginIP, &user.DeviceInfo,
		&user.Language, &user.Timezone, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// updateLastLogin updates the user's last login timestamp and IP
func (h *AuthHandler) updateLastLogin(userID uuid.UUID, clientIP string) error {
	query := `
		UPDATE users 
		SET last_login_at = NOW(), last_login_ip = $2, updated_at = NOW()
		WHERE id = $1
	`

	_, err := h.DB.Exec(query, userID, clientIP)
	return err
}

// RequestOTP godoc
// @Summary Request an OTP code
// @Description Generate and send an OTP code for email verification or other purposes
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RequestOTPRequest true "OTP request data"
// @Success 200 {object} map[string]interface{} "OTP sent successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - validation errors"
// @Failure 401 {object} map[string]interface{} "Unauthorized - invalid or missing token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/otp/request [post]
func (h *AuthHandler) RequestOTP(c *gin.Context) {
	// Get user from token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "missing_authorization",
			"message": "Authorization header is required",
		})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := models.ValidateAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_token",
			"message": "Invalid or expired access token",
		})
		return
	}

	var req models.RequestOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_failed",
			"message": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Get user from database
	user, err := h.getUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to retrieve user",
		})
		return
	}

	// For email-verification, check if already verified
	if req.Type == "email-verification" && user.EmailStatus {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "already_verified",
			"message": "Email is already verified",
		})
		return
	}

	// Generate 6-digit OTP code
	otpCode, err := h.generateOTPCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "otp_generation_error",
			"message": "Failed to generate OTP code",
		})
		return
	}

	// Expires in 10 minutes
	expiresAt := time.Now().Add(10 * time.Minute)

	// Mark all previous unused OTPs of the same type as used (only accept latest)
	_, err = h.DB.Exec(`
		UPDATE otps 
		SET used = TRUE, updated_at = NOW()
		WHERE user_id = $1 AND type = $2 AND used = FALSE
	`, user.ID, req.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to invalidate previous OTPs",
		})
		return
	}

	// Insert new OTP
	otpID := uuid.New()
	_, err = h.DB.Exec(`
		INSERT INTO otps (id, user_id, type, code, used, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`, otpID, user.ID, req.Type, otpCode, false, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to save OTP",
		})
		return
	}

	// Send email with OTP code (for all types that use email)
	emailService := services.NewEmailService()

	// Determine if this OTP type should send email
	emailOTPTypes := map[string]bool{
		"email-verification": true,
		"password-reset":     true,
		"2fa":                true,
		"phone-verification": true,
	}

	shouldSendEmail := emailOTPTypes[req.Type]

	if shouldSendEmail && emailService.IsEnabled() {
		// Send email
		userName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		if userName == " " {
			userName = user.Username
		}

		err = emailService.SendOTPEmail(req.Type, user.Email, userName, otpCode)
		if err != nil {
			// Log error but don't fail the request - OTP is already saved
			// In production, you might want to log this to a proper logging service
			fmt.Printf("Failed to send email: %v\n", err)
			c.JSON(http.StatusOK, gin.H{
				"message":    "OTP code generated, but email sending failed. Please contact support.",
				"error":      "email_send_failed",
				"expires_in": 600,
				// Only include OTP in development if email fails
				"otp_code": otpCode,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "OTP code sent to your email",
			"expires_in": 600, // 10 minutes in seconds
		})
	} else if shouldSendEmail && !emailService.IsEnabled() {
		// SMTP not enabled - return OTP in response (development mode)
		c.JSON(http.StatusOK, gin.H{
			"message":    "OTP code generated (SMTP not enabled)",
			"otp_code":   otpCode, // Only in development
			"expires_in": 600,
		})
	} else {
		// OTP type that doesn't use email
		c.JSON(http.StatusOK, gin.H{
			"message":    "OTP code generated successfully",
			"expires_in": 600,
		})
	}
}

// VerifyOTP godoc
// @Summary Verify an OTP code
// @Description Verify an OTP code for email verification or other purposes
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.VerifyOTPRequest true "OTP verification data"
// @Success 200 {object} map[string]interface{} "OTP verified successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - validation errors"
// @Failure 401 {object} map[string]interface{} "Unauthorized - invalid or missing token"
// @Failure 404 {object} map[string]interface{} "OTP not found or expired"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/otp/verify [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	// Get user from token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "missing_authorization",
			"message": "Authorization header is required",
		})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := models.ValidateAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_token",
			"message": "Invalid or expired access token",
		})
		return
	}

	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_failed",
			"message": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Get the latest unused OTP for this user and type
	var otp models.OTP
	err = h.DB.QueryRow(`
		SELECT id, user_id, type, code, used, expires_at, created_at, updated_at
		FROM otps
		WHERE user_id = $1 AND type = $2 AND used = FALSE
		ORDER BY created_at DESC
		LIMIT 1
	`, claims.UserID, req.Type).Scan(
		&otp.ID, &otp.UserID, &otp.Type, &otp.Code, &otp.Used, &otp.ExpiresAt, &otp.CreatedAt, &otp.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "otp_not_found",
			"message": "No valid OTP found. Please request a new code.",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to retrieve OTP",
		})
		return
	}

	// Check if OTP is expired
	if time.Now().After(otp.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "otp_expired",
			"message": "OTP code has expired. Please request a new code.",
		})
		return
	}

	// Verify the code
	if otp.Code != req.Code {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_code",
			"message": "Invalid OTP code",
		})
		return
	}

	// Mark OTP as used
	_, err = h.DB.Exec(`
		UPDATE otps 
		SET used = TRUE, updated_at = NOW()
		WHERE id = $1
	`, otp.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to mark OTP as used",
		})
		return
	}

	// If email verification, update user's email_status and status
	if req.Type == "email-verification" {
		_, err = h.DB.Exec(`
			UPDATE users 
			SET email_status = TRUE, status = 'active', updated_at = NOW()
			WHERE id = $1
		`, claims.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "database_error",
				"message": "Failed to update email status",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "OTP verified successfully",
		"verified": true,
	})
}

// generateOTPCode generates a random 6-digit OTP code
func (h *AuthHandler) generateOTPCode() (string, error) {
	// Generate random 6-digit code (000000 to 999999)
	// Using crypto/rand for secure random generation
	max := big.NewInt(1000000) // 0 to 999999
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Format as 6-digit string with leading zeros
	return fmt.Sprintf("%06d", n.Int64()), nil
}
