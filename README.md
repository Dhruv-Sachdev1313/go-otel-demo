# Go OpenTelemetry Demo with Clean Architecture

This Go application demonstrates OpenTelemetry metrics and traces integration using an OTEL Collector that forwards data to SigNoz cloud. The code is organized with clean architecture principles and implements all three metric types plus distributed tracing.

## 📁 Project Structure

```
├── main.go                          # Application entry point
├── internal/                        # Internal packages
│   ├── config/                      # Configuration management
│   │   └── config.go               # Environment variable handling
│   ├── handlers/                    # HTTP request handlers
│   │   └── handlers.go             # All endpoint handlers
│   ├── middleware/                  # HTTP middleware
│   │   └── metrics.go              # Metrics and tracing middleware
│   ├── models/                      # Data models
│   │   └── cart.go                 # Shopping cart model
│   └── telemetry/                   # OpenTelemetry setup
│       └── telemetry.go            # Metrics and tracing initialization
├── otel-collector-config.yaml       # OTEL Collector configuration
├── docker-compose.yml               # Docker setup
└── README.md                        # This file
```

## 🎯 Features Implemented

### Metrics (All 3 Types)
1. **Counter** - `http_errors_total`: Counts HTTP error requests (4xx/5xx)
2. **Histogram** - `http_request_duration_seconds`: Records request latency distribution  
3. **Gauge** - `cart_items_count`: Measures current items in user shopping carts

### Distributed Tracing
- Request-level tracing with span attributes
- Error tracking and performance analysis
- Cart operations instrumentation

## 🚀 Architecture Flow

```
HTTP Request → Middleware → Handler → Telemetry → OTEL Collector → SigNoz Cloud
```

## 🔧 Configuration

### Environment Variables
- **SERVICE_NAME**: Service name for tracing (default: "go-otel-demo")
- **OTEL_EXPORTER_OTLP_ENDPOINT**: OTEL Collector endpoint (default: "localhost:4317")
- **INSECURE_MODE**: Use insecure connection (default: "true")
- **SERVER_PORT**: Server port (default: ":8080")

## 🏃 Quick Start

### 1. Start OTEL Collector
```bash
# Using your existing OTEL collector setup
# Make sure it's running on localhost:4317
```

### 2. Run the Application
```bash
./run-with-collector.sh
```

Or manually:
```bash
export SERVICE_NAME="go-otel-demo"
export INSECURE_MODE="true" 
export OTEL_EXPORTER_OTLP_ENDPOINT="localhost:4317"
export SERVER_PORT=":8080"
go run main.go
```

### 3. Test the Endpoints
```bash
./test-app.sh
```

## 📊 API Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```

### Generate Errors (for Counter Metrics)
```bash
curl http://localhost:8080/error
```

### Cart Operations (for Gauge Metrics)
```bash
# Add items
curl "http://localhost:8080/cart/add?user_id=user1&item=laptop"
curl "http://localhost:8080/cart/add?user_id=user1&item=mouse"

# View cart
curl "http://localhost:8080/cart/get?user_id=user1"

# Remove item
curl "http://localhost:8080/cart/remove?user_id=user1&index=0"
```
## 🎛️ What You'll See in SigNoz

### Metrics Dashboard
- **http_errors_total**: Error counts by endpoint and status code
- **http_request_duration_seconds**: Request latency histograms with percentiles
- **cart_items_count**: Real-time cart sizes per user

### Traces Dashboard  
- Request traces showing complete request flow
- Error traces with detailed context and stack information
- Cart operation traces with user and cart size information

## 🔄 Development Workflow

1. **Add New Endpoints**: Create handler methods in `internal/handlers/`
2. **Add Configuration**: Extend `internal/config/config.go`
3. **Add Models**: Create new types in `internal/models/`
4. **Add Middleware**: Extend `internal/middleware/`
5. **Modify Telemetry**: Update `internal/telemetry/telemetry.go`

This structure makes the codebase much more maintainable and follows Go best practices for project organization!

## API Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```

### Simulate Errors
```bash
curl http://localhost:8080/error
```
This endpoint randomly returns errors (30% chance of 500, 20% chance of 400).

### Shopping Cart Operations

#### Add Item to Cart
```bash
curl "http://localhost:8080/cart/add?user_id=user1&item=laptop"
```

#### Remove Item from Cart
```bash
curl "http://localhost:8080/cart/remove?user_id=user1&index=0"
```

#### Get Cart Contents
```bash
curl "http://localhost:8080/cart/get?user_id=user1"
```

## Testing the Metrics

### 1. Generate Traffic
Use the provided endpoints to generate different types of metrics:

```bash
# Generate successful requests
for i in {1..20}; do curl http://localhost:8080/health; done

# Generate errors
for i in {1..10}; do curl http://localhost:8080/error; done

# Add items to carts
curl "http://localhost:8080/cart/add?user_id=user1&item=laptop"
curl "http://localhost:8080/cart/add?user_id=user1&item=mouse"
curl "http://localhost:8080/cart/add?user_id=user2&item=keyboard"
```

### 2. Background Activity
The application automatically generates background activity:
- Randomly adds/removes items from user carts every 5 seconds
- This demonstrates the gauge metric changing over time

## Viewing Metrics in SigNoz

Once the application is running and sending metrics to SigNoz, you can:

1. **Error Rate Dashboard**: Monitor `http_errors_total` to track error rates
2. **Latency Analysis**: Use `http_request_duration_seconds` histogram for latency percentiles
3. **Cart Monitoring**: Track `cart_items_count` gauge for real-time cart sizes

## Metric Types Explained

### Counter (`http_errors_total`)
- Monotonically increasing value
- Perfect for counting events like errors, requests, etc.
- Can only go up (never decreases)

### Histogram (`http_request_duration_seconds`)
- Records observations in buckets
- Provides count, sum, and bucket distributions
- Ideal for latency, request sizes, etc.

### Gauge (`cart_items_count`)
- Can go up and down
- Represents a current value at a point in time
- Perfect for things like queue lengths, active connections, cart sizes

## Architecture

The application uses:
- **OTLP HTTP Exporter**: Sends metrics to SigNoz cloud
- **Periodic Reader**: Exports metrics every 10 seconds
- **Middleware Pattern**: Automatically captures HTTP metrics
- **Observable Gauge**: Uses callback pattern for real-time cart monitoring

## Troubleshooting

1. **No Metrics in SigNoz**: Check environment variables and network connectivity
2. **Authentication Issues**: Verify your SigNoz access token
3. **Endpoint Issues**: Ensure the correct region in the OTLP endpoint URL
