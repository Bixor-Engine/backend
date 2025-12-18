package models

import (
	"time"

	"github.com/google/uuid"
)

// SchemaMigration tracks database schema versions
type SchemaMigration struct {
	Version     string    `json:"version" db:"version"`
	Description string    `json:"description" db:"description"`
	AppliedAt   time.Time `json:"applied_at" db:"applied_at"`
	Checksum    string    `json:"checksum" db:"checksum"`
	AppliedBy   string    `json:"applied_by" db:"applied_by"`
}

// DatabaseInfo stores database metadata and version information
type DatabaseInfo struct {
	ID            int       `json:"id" db:"id"`
	DBVersion     string    `json:"db_version" db:"db_version"`
	SchemaVersion string    `json:"schema_version" db:"schema_version"`
	EngineVersion string    `json:"engine_version" db:"engine_version"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	Environment   string    `json:"environment" db:"environment"`
}

// User represents a complete system user matching the database schema
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	FirstName    string     `json:"first_name" db:"first_name"`
	LastName     string     `json:"last_name" db:"last_name"`
	Username     string     `json:"username" db:"username"`
	Email        string     `json:"email" db:"email"`
	Password     string     `json:"-" db:"password"` // Never expose password hash in JSON
	EmailStatus  bool       `json:"email_status" db:"email_status"`
	PhoneNumber  *string    `json:"phone_number,omitempty" db:"phone_number"`
	PhoneStatus  bool       `json:"phone_status" db:"phone_status"`
	ReferredBy   *uuid.UUID `json:"referred_by,omitempty" db:"referred_by"`
	Address      *string    `json:"address,omitempty" db:"address"`
	City         *string    `json:"city,omitempty" db:"city"`
	Country      *string    `json:"country,omitempty" db:"country"`
	Role         string     `json:"role" db:"role"`
	Status       string     `json:"status" db:"status"`
	KYCStatus    string     `json:"kyc_status" db:"kyc_status"`
	TwoFAEnabled bool       `json:"twofa_enabled" db:"twofa_enabled"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	LastLoginIP  *string    `json:"last_login_ip,omitempty" db:"last_login_ip"`
	DeviceInfo   *string    `json:"device_info,omitempty" db:"device_info"` // JSONB as string
	Language     string     `json:"language" db:"language"`
	Timezone     string     `json:"timezone" db:"timezone"`
	GlobalBalance float64   `json:"global_balance" db:"global_balance"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// RegisterRequest represents the request payload for user registration
type RegisterRequest struct {
	FirstName   string  `json:"first_name" binding:"required,min=2,max=50"`
	LastName    string  `json:"last_name" binding:"required,min=2,max=50"`
	Username    string  `json:"username" binding:"required,min=3,max=30,alphanum"`
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=8,max=128"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	ReferredBy  *string `json:"referred_by,omitempty"` // UUID as string in request
	Address     *string `json:"address,omitempty"`
	City        *string `json:"city,omitempty"`
	Country     *string `json:"country,omitempty"`
	Language    *string `json:"language,omitempty"`
	Timezone    *string `json:"timezone,omitempty"`
}

// UserResponse represents the response payload for user data (excluding sensitive fields)
type UserResponse struct {
	ID           uuid.UUID  `json:"id"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	EmailStatus  bool       `json:"email_status"`
	PhoneNumber  *string    `json:"phone_number,omitempty"`
	PhoneStatus  bool       `json:"phone_status"`
	Address      *string    `json:"address,omitempty"`
	City         *string    `json:"city,omitempty"`
	Country      *string    `json:"country,omitempty"`
	Role         string     `json:"role"`
	Status       string     `json:"status"`
	KYCStatus    string     `json:"kyc_status"`
	TwoFAEnabled bool       `json:"twofa_enabled"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	Language     string     `json:"language"`
	Timezone     string     `json:"timezone"`
	GlobalBalance float64   `json:"global_balance"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// SystemHealth tracks system health metrics
type SystemHealth struct {
	ID             int       `json:"id" db:"id"`
	ServiceName    string    `json:"service_name" db:"service_name"`
	Status         string    `json:"status" db:"status"`
	ResponseTimeMs *int      `json:"response_time_ms,omitempty" db:"response_time_ms"`
	Timestamp      time.Time `json:"timestamp" db:"timestamp"`
	Details        string    `json:"details,omitempty" db:"details"` // JSON string
	Environment    string    `json:"environment" db:"environment"`
}

// Market represents a trading market/pair
type Market struct {
	ID                uuid.UUID `json:"id" db:"id"`
	Symbol            string    `json:"symbol" db:"symbol"`
	BaseCurrency      string    `json:"base_currency" db:"base_currency"`
	QuoteCurrency     string    `json:"quote_currency" db:"quote_currency"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	MinQuantity       string    `json:"min_quantity" db:"min_quantity"` // Use string for precise decimal handling
	MaxQuantity       *string   `json:"max_quantity,omitempty" db:"max_quantity"`
	PricePrecision    int       `json:"price_precision" db:"price_precision"`
	QuantityPrecision int       `json:"quantity_precision" db:"quantity_precision"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=1"`
}

// LoginResponse represents the response payload for successful login
type LoginResponse struct {
	Message        string       `json:"message"`
	User           UserResponse `json:"user"`
	Tokens         JWTTokens    `json:"tokens"`
	RequiresVerify bool         `json:"requires_verify"`       // True if user needs to verify email
	RedirectTo     string       `json:"redirect_to,omitempty"` // Where to redirect (e.g., "/verify-email")
}

// RefreshTokenRequest represents the request payload for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// OTP represents a one-time password record
type OTP struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Type      string    `json:"type" db:"type"` // e.g., 'email-verification', 'password-reset', '2fa'
	Code      string    `json:"code" db:"code"` // 6-digit code
	Used      bool      `json:"used" db:"used"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// RequestOTPRequest represents the request payload for requesting an OTP
type RequestOTPRequest struct {
	Type string `json:"type" binding:"required,oneof=email-verification password-reset 2fa phone-verification"`
}

// VerifyOTPRequest represents the request payload for verifying an OTP
type VerifyOTPRequest struct {
	Type string `json:"type" binding:"required,oneof=email-verification password-reset 2fa phone-verification"`
	Code string `json:"code" binding:"required,len=6,numeric"`
}

// RequestOTPResponse represents the response for requesting an OTP
type RequestOTPResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in"`         // Expiration time in seconds
	OTPCode   string `json:"otp_code,omitempty"` // Only included in development mode when SMTP is disabled
	Error     string `json:"error,omitempty"`    // Only included if email sending failed
}

// VerifyOTPResponse represents the response for verifying an OTP
type VerifyOTPResponse struct {
	Message  string `json:"message"`
	Verified bool   `json:"verified"`
}

// Coin represents a cryptocurrency coin
type Coin struct {
	ID              int       `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Ticker          string    `json:"ticker" db:"ticker"`
	Decimal         int       `json:"decimal" db:"decimal"`
	PriceDecimal    int       `json:"price_decimal" db:"price_decimal"`
	Logo            *string   `json:"logo,omitempty" db:"logo"`
	Price           string    `json:"price" db:"price"` // Using string to preserve precision
	DepositGateway  []string  `json:"deposit_gateway" db:"deposit_gateway"`
	WithdrawGateway []string  `json:"withdraw_gateway" db:"withdraw_gateway"`
	DepositFee      *string   `json:"deposit_fee,omitempty" db:"deposit_fee"`
	WithdrawFee     *string   `json:"withdraw_fee,omitempty" db:"withdraw_fee"`
	DepositFeeType  *int      `json:"deposit_fee_type,omitempty" db:"deposit_fee_type"`
	WithdrawFeeType *int      `json:"withdraw_fee_type,omitempty" db:"withdraw_fee_type"`
	Confirmation    *int      `json:"confirmation,omitempty" db:"confirmation"`
	Status          int       `json:"status" db:"status"`
	WithdrawStatus  *int      `json:"withdraw_status,omitempty" db:"withdraw_status"`
	DepositStatus   *int      `json:"deposit_status,omitempty" db:"deposit_status"`
	Website         *string   `json:"website,omitempty" db:"website"`
	Explorer        *string   `json:"explorer,omitempty" db:"explorer"`
	ExplorerTx      *string   `json:"explorer_tx,omitempty" db:"explorer_tx"`
	ExplorerAddress *string   `json:"explorer_address,omitempty" db:"explorer_address"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// CoinListResponse represents the response for listing coins
type CoinListResponse struct {
	Coins []Coin `json:"coins"`
	Total int    `json:"total"`
}

// CoinResponse represents a single coin response
type CoinResponse struct {
	Coin Coin `json:"coin"`
}
