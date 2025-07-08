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

// User represents a system user
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"` // Never expose password hash in JSON
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	Version      int       `json:"version" db:"version"`
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
