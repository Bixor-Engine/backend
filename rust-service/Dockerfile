# Build stage
FROM rust:latest as builder

WORKDIR /app

# Copy Cargo files
COPY Cargo.toml ./

# Create a dummy main.rs to build dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs

# Build dependencies (this will be cached as long as Cargo.toml doesn't change)
RUN cargo build --release
RUN rm src/main.rs

# Copy source code
COPY src ./src

# Build application
RUN cargo build --release

# Runtime stage
FROM debian:bookworm-slim

# Install required packages
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/target/release/bixor-rust-service /app/

# Expose port
EXPOSE 8081

# Run the application
CMD ["./bixor-rust-service"] 