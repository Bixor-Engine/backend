use axum::{
    extract::State,
    http::StatusCode,
    response::Json,
    routing::get,
    Router,
};
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::{PgPool, Row};
use std::{net::SocketAddr, sync::Arc, time::Instant};
use tower_http::cors::CorsLayer;
use tracing_subscriber;

#[derive(Clone)]
struct AppState {
    db: PgPool,
}

#[derive(Serialize)]
struct HealthResponse {
    status: String,
    service: String,
    timestamp: DateTime<Utc>,
    database: DatabaseHealth,
    details: std::collections::HashMap<String, String>,
}

#[derive(Serialize)]
struct DatabaseHealth {
    status: String,
    response_time_ms: u64,
}

#[derive(Serialize)]
struct StatusResponse {
    message: String,
    timestamp: DateTime<Utc>,
    service: String,
    version: String,
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // Initialize tracing
    tracing_subscriber::init();

    // Load environment variables
    dotenv::dotenv().ok();

    // Get database URL
    let database_url = std::env::var("DATABASE_URL")
        .unwrap_or_else(|_| "postgres://bixor_user:bixor_pass@localhost:5432/bixor".to_string());

    // Create database connection pool
    let db = PgPool::connect(&database_url).await?;
    
    // Test database connection
    sqlx::query("SELECT 1").fetch_one(&db).await?;
    tracing::info!("Database connection established");

    let state = AppState { db };

    // Build our application with routes
    let app = Router::new()
        .route("/health", get(health_check))
        .route("/api/v1/status", get(get_status))
        .layer(CorsLayer::permissive())
        .with_state(state);

    // Get port from environment or default to 8081
    let port = std::env::var("PORT")
        .unwrap_or_else(|_| "8081".to_string())
        .parse::<u16>()
        .unwrap_or(8081);

    let addr = SocketAddr::from(([0, 0, 0, 0], port));
    tracing::info!("Rust service starting on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await?;
    axum::serve(listener, app).await?;

    Ok(())
}

async fn health_check(State(state): State<AppState>) -> Result<Json<HealthResponse>, StatusCode> {
    let start = Instant::now();
    
    // Check database health
    let db_health = match sqlx::query("SELECT 1")
        .fetch_one(&state.db)
        .await
    {
        Ok(_) => DatabaseHealth {
            status: "healthy".to_string(),
            response_time_ms: start.elapsed().as_millis() as u64,
        },
        Err(e) => {
            tracing::error!("Database health check failed: {}", e);
            DatabaseHealth {
                status: "unhealthy".to_string(),
                response_time_ms: start.elapsed().as_millis() as u64,
            }
        }
    };

    let mut details = std::collections::HashMap::new();
    details.insert("version".to_string(), "1.0.0".to_string());
    details.insert("uptime".to_string(), format!("{}ms", start.elapsed().as_millis()));

    let response = HealthResponse {
        status: if db_health.status == "healthy" { "healthy".to_string() } else { "unhealthy".to_string() },
        service: "bixor-rust-service".to_string(),
        timestamp: Utc::now(),
        database: db_health,
        details,
    };

    if response.status == "unhealthy" {
        return Err(StatusCode::SERVICE_UNAVAILABLE);
    }

    Ok(Json(response))
}

async fn get_status() -> Json<StatusResponse> {
    Json(StatusResponse {
        message: "Bixor Rust Service is running".to_string(),
        timestamp: Utc::now(),
        service: "bixor-rust-service".to_string(),
        version: "1.0.0".to_string(),
    })
} 