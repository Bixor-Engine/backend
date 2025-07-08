# Bixor Trading Engine

A high-performance trading backend built with Go.

## Architecture

- **Go Service**: REST API, WebSocket handlers, business logic, order matching, risk calculations
- **PostgreSQL**: Primary database for persistent data

## Getting Started

### Prerequisites
- Go 1.21+
- Rust 1.70+
- Docker & Docker Compose
- PostgreSQL 15+

### Quick Start

```bash
# Start all services
docker-compose up -d

# Health check
curl http://localhost:8080/health  # Go service
```

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
