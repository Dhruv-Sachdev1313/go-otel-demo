package main

import (
	"context"
	"log"
	"net/http"

	"go-otel-demo/internal/config"
	"go-otel-demo/internal/handlers"
	"go-otel-demo/internal/middleware"
	"go-otel-demo/internal/telemetry"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize telemetry (tracing and metrics)
	metricsCollector, cleanup, err := telemetry.Setup(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}
	defer cleanup(context.Background())

	// Initialize handlers
	handler := handlers.New(cfg, metricsCollector)

	// Create middleware
	metricsMiddleware := middleware.MetricsMiddleware(cfg, metricsCollector)

	// Setup HTTP handlers with metrics middleware
	http.HandleFunc("/health", metricsMiddleware(handler.Health))
	http.HandleFunc("/error", metricsMiddleware(handler.Error))
	http.HandleFunc("/cart/add", metricsMiddleware(handler.AddToCart))
	http.HandleFunc("/cart/remove", metricsMiddleware(handler.RemoveFromCart))
	http.HandleFunc("/cart/get", metricsMiddleware(handler.GetCart))

	log.Printf("Starting HTTP server on %s", cfg.ServerPort)
	log.Printf("Service Name: %s", cfg.ServiceName)
	log.Printf("OTEL Collector Endpoint: %s", cfg.CollectorURL)
	log.Printf("Insecure Mode: %s", cfg.Insecure)
	log.Println("")
	log.Println("Available endpoints:")
	log.Println("  GET /health - Health check")
	log.Println("  GET /error - Simulates random errors")
	log.Println("  GET /cart/add?user_id=USER&item=ITEM - Add item to cart")
	log.Println("  GET /cart/remove?user_id=USER&index=INDEX - Remove item from cart")
	log.Println("  GET /cart/get?user_id=USER - Get cart contents")
	log.Println("")
	log.Println("Environment variables:")
	log.Println("  SERVICE_NAME - Service name for tracing")
	log.Println("  OTEL_EXPORTER_OTLP_ENDPOINT - OTEL Collector endpoint")
	log.Println("  INSECURE_MODE - Use insecure connection (true/false)")
	log.Println("  SERVER_PORT - Server port (default: :8080)")
	log.Println("")
	log.Println("Metrics and Traces being sent:")
	log.Println("  - http_errors_total (Counter)")
	log.Println("  - http_request_duration_seconds (Histogram)")
	log.Println("  - cart_items_count (Gauge)")
	log.Println("  - Distributed traces for all requests")

	if err := http.ListenAndServe(cfg.ServerPort, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
