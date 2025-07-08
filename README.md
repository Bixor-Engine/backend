# Bixor Trading Engine

A high-performance trading backend built with Go.

## Architecture

- **Go Service**: REST API, WebSocket handlers, business logic, order matching, risk calculations
- **PostgreSQL**: Primary database for persistent data

## Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Quick Start

#### Using Docker Compose
```bash
# Start all services
docker-compose up -d

# Health check
curl http://localhost:8080/api/v1/health
```

#### Manual Setup
```bash
# 1. Install the CLI tool
go install ./cmd/bixor

# 2. Setup database (make sure PostgreSQL is running)
bixor database setup

# 3. Start the API service
bixor server start
```

## CLI Tool

The project includes a professional CLI tool built with Cobra for managing all aspects of the Bixor Engine:

### Installation
```bash
# Install globally
go install ./cmd/bixor

# Verify installation
bixor --version
```

### Database Management

| Command | Description |
|---------|-------------|
| `bixor database setup` | Complete first-time database setup (create + migrate + seed) |
| `bixor database migrate` | Run pending database migrations |
| `bixor database seed` | Populate database with sample data |
| `bixor database reset` | Drop all tables and rebuild (destructive) |
| `bixor database status` | Show migration status and database info |

### Server Management

| Command | Description |
|---------|-------------|
| `bixor server start` | Start the API server |

### Examples
```bash
# First-time setup
bixor database setup

# Run new migrations
bixor database migrate

# Check database status
bixor database status

# Start the server
bixor server start

# Reset database (WARNING: destroys all data)
bixor database reset
```

For detailed CLI documentation, see [CLI README](cmd/bixor/README.md).

## API Documentation

### Interactive Documentation
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Quick Access**: http://localhost:8080/docs

### API Endpoints

#### General
- `GET /` - Landing page
- `GET /api/v1/info` - API information and service details

#### Monitoring
- `GET /api/v1/health` - Health check for all services
- `GET /api/v1/status` - Status check for all services

## Services

### Go Service (Port 8080)
- REST API endpoints
- WebSocket connections
- User management
- Portfolio tracking
- Order matching engine
- Risk calculations
- Market data processing

### Database
- PostgreSQL on port 5432
- Contains user data, orders, trades, market data
- Managed via migrations and seeders

## Development

### Project Structure
```
bixor-engine/
├── cmd/
│   ├── bixor/           # CLI tool for database and server management
│   └── server/          # Main API server entry point
├── internal/
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Database models
│   └── routes/          # Route definitions
├── database/
│   ├── migrations/      # Database schema migrations
│   ├── seeders/         # Database seed data
│   └── init.sql         # Basic database initialization
├── docs/                # Auto-generated API documentation
└── configs/             # Configuration files
```

### Adding New Database Changes

1. Create a new migration file:
```sql
-- database/migrations/007_add_new_table.sql
CREATE TABLE new_table (...);
INSERT INTO schema_migrations (version, description, checksum) 
VALUES ('007', 'Add new table', 'migration_007_checksum');
```

2. Run the migration:
```bash
bixor database migrate
```

### Environment Variables

Create a `.env` file in the project root:
```env
# Database Configuration
DATABASE_URL=postgres://username:password@localhost:5432/bixor?sslmode=disable

# Go Service Configuration
PORT=8080
API_VERSION=1.0.0
```

## API Testing

### Using curl
```bash
# Health check
curl http://localhost:8080/api/v1/health

# Service status
curl http://localhost:8080/api/v1/status

# API information
curl http://localhost:8080/api/v1/info
```

### Using the Interactive Documentation
Visit http://localhost:8080/swagger/index.html to test endpoints directly in your browser.

## License

MIT License - Build your exchange freely, safely, and securely. 
