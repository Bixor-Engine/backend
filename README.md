# Bixor Trading Engine

A high-performance cryptocurrency trading backend built with Go and PostgreSQL.

## Project Status

**Current State**: Early development phase with basic authentication system implemented.

This project is actively under development. The authentication system is functional, but most trading-specific features are not yet implemented.

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher
- Git

## Quick Start

### 1. Database Setup

Start PostgreSQL (example with Docker):
```bash
docker run --name bixor-postgres \
  -e POSTGRES_USER=bixor \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=bixor \
  -p 5432:5432 \
  -d postgres:15
```

### 2. Environment Configuration

Create a `.env` file in the project root:
```bash
cp .env.example .env
```

Edit `.env` with your database credentials:
```env
DATABASE_URL=postgres://bixor:password@localhost:5432/bixor?sslmode=disable
PORT=8080
API_VERSION=1.0.0
JWT_SECRET=your_secure_jwt_secret_key_here
JWT_EXPIRES_HOURS=24
BACKEND_SECRET=your_backend_secret_key_here_change_in_production

# SMTP Configuration (optional, for email sending)
SMTP_ENABLED=false
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=your_smtp_username
SMTP_PASSWORD=your_smtp_password
SMTP_FROM_EMAIL=noreply@example.com
SMTP_FROM_NAME=Bixor Engine
```

### 3. Install and Setup

```bash
# Download Go dependencies
go mod download

# Install CLI tool
go install ./cmd/bixor

# Setup database (creates tables and initial data)
bixor database setup

# Start the API server
bixor server start
```

The API server will be available at `http://localhost:8080`

## Available Features

### Public API Endpoints (No Authentication Required)
- **GET /api/v1/health** - Health check for API and database
- **GET /api/v1/status** - Service status information
- **GET /api/v1/info** - API information and available endpoints
- **GET /api/v1/currency** - List all supported cryptocurrencies
- **GET /api/v1/currency/:ticker** - Get coin information by ticker

### Private API Endpoints (Backend Secret Required)
- **POST /api/v1/auth/register** - User registration
- **POST /api/v1/auth/login** - User login with JWT tokens
- **POST /api/v1/auth/refresh** - JWT token refresh
- **GET /api/v1/auth/me** - Get current authenticated user
- **POST /api/v1/auth/logout** - Logout user
- **POST /api/v1/auth/otp/request** - Request OTP code
- **POST /api/v1/auth/otp/verify** - Verify OTP code

### API Documentation

The API documentation is organized into separate sections based on access level:

- **Documentation Landing Page**: http://localhost:8080/docs
  - Navigate to different API sections from here

- **Public API**: http://localhost:8080/docs/public
  - System health, status, and currency information endpoints
  - No authentication required

- **Private API**: http://localhost:8080/docs/private
  - Authentication, registration, and user management endpoints
  - Requires backend secret (X-Backend-Secret header)

- **Personal API**: http://localhost:8080/docs/personal
  - User-specific endpoints for personal data and operations
  - Requires personal token (X-Personal-Token header)

## Database Schema

Currently implemented tables:

### Users Table
Complete user management with the following fields:
- User identification (id, username, email)
- Personal information (first_name, last_name, phone_number, address)
- Authentication (password with Argon2i hashing)
- Account status (role, status, kyc_status)
- Security features (twofa_enabled, last_login tracking)
- Referral system support
- Localization (language, timezone)

### OTPs Table
One-Time Password management for email verification and other verification purposes:
- OTP codes with expiration
- Support for multiple OTP types (email-verification, password-reset, 2fa, phone-verification)
- Automatic invalidation of old unused codes

### Coins Table
Cryptocurrency management:
- Coin information (name, ticker, decimals, price)
- Gateway support (deposit/withdraw gateways as arrays)
- Fee configuration (deposit/withdraw fees and types)
- Status management (active/inactive, deposit/withdraw status)
- Explorer links (website, explorer, transaction, address)

**Note**: Trading-specific tables (orders, trades, markets, wallets) are not yet implemented.

## CLI Tool

The `bixor` CLI tool provides database and server management:

```bash
# Database operations
bixor database setup    # Initial database setup
bixor database migrate  # Run pending migrations
bixor database status   # Check migration status
bixor database seed     # Populate with sample data
bixor database reset    # Reset database (destructive)

# Server operations
bixor server start      # Start the API server
```

## Development Tools

Interactive testing tools are available in the `tools/` directory:

```bash
# Authentication testing tool
go run tools/auth/main.go

# Password hashing tool
go run tools/hash/main.go
```

See `tools/README.md` for detailed usage instructions.

## Project Structure

```
bixor-engine/
├── cmd/
│   ├── bixor/           # CLI tool for management
│   └── server/          # Main API server
├── internal/
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Data models and business logic
│   └── routes/          # Route definitions
├── database/
│   ├── migrations/      # Database schema migrations
│   └── seeders/         # Sample data for development
├── tools/               # Development and testing tools
└── docs/                # Auto-generated API documentation
```

## Security Features

- **Password Hashing**: Argon2i with secure parameters
- **JWT Authentication**: Access and refresh token system
- **API Route Protection**: 
  - Public routes: No authentication required
  - Private routes: Protected with backend secret (X-Backend-Secret header)
  - Personal routes: Protected with personal token (X-Personal-Token header) - Future implementation
- **Input Validation**: Request validation with detailed error messages
- **Database Constraints**: Enforced data integrity at database level
- **Email Verification**: OTP-based email verification system

## Testing

```bash
# Test API endpoints
curl http://localhost:8080/api/v1/health

# Interactive testing via Swagger UI
# Visit: http://localhost:8080/docs to access all API documentation sections
```

## Planned Features

The following features are planned but not yet implemented:

### Core Trading Infrastructure
- WebSocket connections for real-time data
- Order management system
- Trade execution engine
- Market data feeds
- Portfolio and wallet management

### Advanced Features
- Risk management system
- Margin trading support
- Advanced order types
- Performance optimization
- Comprehensive monitoring

## API Testing

Example API usage:

```bash
# Test public endpoint (no authentication required)
curl http://localhost:8080/api/v1/health

# Get list of supported coins (public)
curl http://localhost:8080/api/v1/currency

# Register a new user (requires backend secret)
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "X-Backend-Secret: your_backend_secret_key_here" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "securepassword123",
    "first_name": "Test",
    "last_name": "User"
  }'

# Login (requires backend secret)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Backend-Secret: your_backend_secret_key_here" \
  -d '{
    "email": "test@example.com",
    "password": "securepassword123"
  }'

# Get current user (requires backend secret + JWT token)
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "X-Backend-Secret: your_backend_secret_key_here" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

## Contributing

This project is in active development. Contributions are welcome, particularly in the following areas:

1. WebSocket implementation for real-time features
2. Trading infrastructure (orders, trades, markets)
3. Performance optimization
4. Documentation improvements

## License

MIT License

## Support

For questions or issues, please check the documentation or create an issue in the project repository. 
