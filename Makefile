.PHONY: help build run stop clean test logs health

# Default target
help:
	@echo "Bixor Trading Engine - Development Commands"
	@echo ""
	@echo "Available commands:"
	@echo "  build     Build all services"
	@echo "  run       Start all services with docker-compose"
	@echo "  stop      Stop all services"
	@echo "  clean     Stop and remove all containers and volumes"
	@echo "  test      Run health checks on all services"
	@echo "  logs      Show logs for all services"
	@echo "  health    Check health of all services"
	@echo ""

# Build all services
build:
	@echo "Building all services..."
	docker-compose build

# Start all services
run:
	@echo "Starting all services..."
	docker-compose up -d
	@echo "Services started. Run 'make health' to check status."

# Stop all services
stop:
	@echo "Stopping all services..."
	docker-compose down

# Clean everything
clean:
	@echo "Cleaning up all containers and volumes..."
	docker-compose down -v --remove-orphans
	docker system prune -f

# Run health checks
health:
	@echo "Checking service health..."
	@echo ""
	@echo "Go Service Health:"
	@curl -s http://localhost:8080/health | jq . || echo "Go service not responding"
	@echo ""
	@echo "Rust Service Health:"  
	@curl -s http://localhost:8081/health | jq . || echo "Rust service not responding"
	@echo ""

# Show logs
logs:
	docker-compose logs -f

# Test endpoints
test:
	@echo "Testing all endpoints..."
	@echo ""
	@echo "=== Go Service ==="
	@echo "Health endpoint:"
	@curl -s http://localhost:8080/health | jq .
	@echo ""
	@echo "Status endpoint:"
	@curl -s http://localhost:8080/api/v1/status | jq .
	@echo ""
	@echo "=== Rust Service ==="
	@echo "Health endpoint:"
	@curl -s http://localhost:8081/health | jq .
	@echo ""
	@echo "Status endpoint:"
	@curl -s http://localhost:8081/api/v1/status | jq .

# Development commands
dev-go:
	@echo "Starting Go service in development mode..."
	cd go-service && go run main.go

dev-rust:
	@echo "Starting Rust service in development mode..."
	cd rust-service && cargo run 