# Bixor CLI

A cross-platform command-line tool for managing the Bixor Trading Engine.

## Installation

### Automatic Installation
```bash
# Install globally and verify
go install ./cmd/bixor
bixor --version
```

### Manual Installation
```bash
# Build and install manually
cd cmd/bixor
go build -o bixor
# Move to your PATH or use ./bixor
```

### Development Installation
```bash
# Run directly without installing
go run ./cmd/bixor [command]
```

## Quick Start

```bash
# Install the CLI
go install ./cmd/bixor

# Setup database (first time)
bixor database setup

# Start the server
bixor server start
```

## Commands

### Database Management

#### `bixor database setup`
Sets up the database for first-time use:
- Creates database if it doesn't exist
- Runs all pending migrations
- Seeds initial data

```bash
bixor database setup
```

#### `bixor database migrate`
Runs pending database migrations with safety checks:
- Checks which migrations are already applied
- Only runs new migrations
- Asks for confirmation before proceeding

```bash
bixor database migrate
```

#### `bixor database seed`
Populates the database with sample data:

```bash
bixor database seed
```

#### `bixor database reset`
**⚠️ DESTRUCTIVE OPERATION**
Drops all tables and rebuilds the database:
- Drops all existing tables
- Runs fresh migrations
- Seeds initial data

```bash
bixor database reset
```

#### `bixor database status`
Shows current migration status:

```bash
bixor database status
```

### Server Management

#### `bixor server start`
Starts the Bixor API server:

```bash
bixor server start
```

## Configuration

The CLI reads configuration from `.env` file in your project root:

```env
# Database Configuration
DATABASE_URL=postgres://username:password@localhost:5432/bixor?sslmode=disable

# Server Configuration
PORT=8080
API_VERSION=1.0.0
```

## Database Setup

### Option 1: Docker PostgreSQL (Recommended)
```bash
# Start PostgreSQL container
docker run --name bixor-postgres \
  -e POSTGRES_USER=bixor \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=bixor \
  -p 5432:5432 \
  -d postgres:15

# Update .env file
echo "DATABASE_URL=postgres://bixor:password@localhost:5432/bixor?sslmode=disable" > .env

# Setup database
bixor database setup
```

### Default Admin Credentials

After running `bixor database setup`, a default superadmin account is created:

```
Username: superadmin
Email: superadmin@bixor.io
Password: H5w8xLLz}0drN&ey{0?!
```

**⚠️ Important**: Change this password immediately in production environments!

### Option 2: Local PostgreSQL
```bash
# Create database and user
psql -U postgres -c "CREATE DATABASE bixor;"
psql -U postgres -c "CREATE USER bixor WITH PASSWORD 'password';"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE bixor TO bixor;"

# Update .env file
echo "DATABASE_URL=postgres://bixor:password@localhost:5432/bixor?sslmode=disable" > .env

# Setup database
bixor database setup
```

## Usage Examples

### First-Time Setup
```bash
# 1. Install CLI
go install ./cmd/bixor

# 2. Setup environment
echo "DATABASE_URL=postgres://bixor:password@localhost:5432/bixor?sslmode=disable" > .env

# 3. Setup database
bixor database setup

# 4. Start server
bixor server start
```

### Development Workflow
```bash
# Add new migration
# Create: database/migrations/007_add_feature.sql

# Run migration
bixor database migrate

# Reset database for testing
bixor database reset

# Check status
bixor database status
```

## Safety Features

### Migration Safety
- **Version Tracking**: Tracks which migrations have been applied
- **Duplicate Prevention**: Won't run the same migration twice
- **User Confirmation**: Asks for confirmation before destructive operations
- **Connection Testing**: Verifies database connection before operations

### Error Handling
- **Graceful Failures**: Clear error messages with suggested solutions
- **Environment Validation**: Checks for required environment variables
- **Database Validation**: Ensures database is accessible before operations

## Troubleshooting

### Common Issues

#### "Cannot connect to database"
```bash
# Check if PostgreSQL is running
docker ps  # for Docker setup
# or
pg_isready -h localhost -p 5432

# Verify .env file exists and has correct DATABASE_URL
cat .env
```

#### "Database already exists"
```bash
# Use migrate instead of setup for existing databases
bixor database migrate
```

#### "CLI not found"
```bash
# Ensure GOPATH/bin is in your PATH
go env GOPATH

# Or run directly
go run ./cmd/bixor [command]
```

#### "Migration already applied"
```bash
# Check current status
bixor database status

# Reset if needed (WARNING: destroys data)
bixor database reset
```

## Contributing

When adding new CLI commands:

1. Add command definition to `cmd/bixor/commands/`
2. Update this README with new command documentation
3. Add appropriate error handling and user confirmation for destructive operations
4. Test on both Windows and Unix systems 