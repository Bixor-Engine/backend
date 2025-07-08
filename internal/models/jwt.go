package models

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims represents the custom claims for our JWT tokens
type JWTClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	TokenType string    `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// JWTTokens represents a pair of access and refresh tokens
type JWTTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // seconds until expiration
}

var (
	// ErrInvalidToken represents an error when the token is invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrExpiredToken represents an error when the token has expired
	ErrExpiredToken = errors.New("token has expired")

	// ErrInvalidTokenType represents an error when the token type is invalid
	ErrInvalidTokenType = errors.New("invalid token type")
)

// getJWTSecret retrieves the JWT secret from environment variables
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Fallback secret (should not be used in production)
		secret = "bixor_default_jwt_secret_change_in_production"
	}
	return secret
}

// getJWTExpirationHours retrieves the JWT expiration hours from environment variables
func getJWTExpirationHours() time.Duration {
	hoursStr := os.Getenv("JWT_EXPIRES_HOURS")
	if hoursStr == "" {
		return 24 * time.Hour // Default 24 hours
	}

	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		return 24 * time.Hour // Default 24 hours on error
	}

	return time.Duration(hours) * time.Hour
}

// GenerateTokens generates both access and refresh tokens for a user
func GenerateTokens(user *User) (*JWTTokens, error) {
	// Generate access token
	accessToken, err := generateToken(user, "access", getJWTExpirationHours())
	if err != nil {
		return nil, err
	}

	// Generate refresh token (longer expiration)
	refreshToken, err := generateToken(user, "refresh", getJWTExpirationHours()*7) // 7x longer
	if err != nil {
		return nil, err
	}

	return &JWTTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(getJWTExpirationHours().Seconds()),
	}, nil
}

// generateToken creates a JWT token with the specified type and expiration
func generateToken(user *User, tokenType string, expiration time.Duration) (string, error) {
	now := time.Now()
	expirationTime := now.Add(expiration)

	claims := &JWTClaims{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "bixor-engine",
			Subject:   user.ID.String(),
			ID:        uuid.New().String(), // Unique token ID
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(getJWTSecret()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(getJWTSecret()), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Check if token has expired
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}

// ValidateAccessToken validates specifically an access token
func ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "access" {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

// ValidateRefreshToken validates specifically a refresh token
func ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

// RefreshTokens generates new tokens using a valid refresh token
func RefreshTokens(refreshTokenString string, user *User) (*JWTTokens, error) {
	// Validate the refresh token
	_, err := ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	return GenerateTokens(user)
}
