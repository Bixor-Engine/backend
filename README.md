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

### Authentication API
- **POST /api/v1/auth/register** - User registration
- **POST /api/v1/auth/login** - User login with JWT tokens
- **POST /api/v1/auth/refresh** - JWT token refresh

### System Monitoring
- **GET /api/v1/health** - Health check for API and database
- **GET /api/v1/status** - Service status information
- **GET /api/v1/info** - API information and available endpoints

### API Documentation
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Quick Access**: http://localhost:8080/docs

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
- **Input Validation**: Request validation with detailed error messages
- **Database Constraints**: Enforced data integrity at database level

## Testing

```bash
# Test API endpoints
curl http://localhost:8080/api/v1/health

# Interactive testing via Swagger UI
# Visit: http://localhost:8080/swagger/index.html
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

Example authentication flow:

```bash
# Register a new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "securepassword123",
    "first_name": "Test",
    "last_name": "User"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "securepassword123"
  }'
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
