# Bixor Trading Engine

A high-performance trading backend built with Go and Rust microservices.

## Architecture

- **Go Service**: REST API, WebSocket handlers, business logic
- **Rust Service**: High-performance order matching, risk calculations
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

# Health checks
curl http://localhost:8080/health  # Go service
curl http://localhost:8081/health  # Rust service
```

## Services

### Go Service (Port 8080)
- REST API endpoints
- WebSocket connections
- User management
- Portfolio tracking

### Rust Service (Port 8081)  
- Order matching engine
- Risk calculations
- Market data processing
- Performance-critical operations

### Database
- PostgreSQL on port 5432
- Contains user data, orders, trades, market data 