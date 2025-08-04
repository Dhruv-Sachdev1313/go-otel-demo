# Complete Setup Guide - OTEL Collector with SigNoz

## Quick Start (Recommended Method)

### Step 1: Update OTEL Collector Configuration
Edit `otel-collector-config.yaml` and update your SigNoz details:

```yaml
exporters:
  otlp/signoz:
    endpoint: https://ingest.{your-region}.signoz.cloud:443
    headers:
      "signoz-access-token": "your-actual-access-token"
```

Replace:
- `{your-region}` with your SigNoz region (`us`, `eu`, `in`, etc.)
- `your-actual-access-token` with your actual SigNoz access token

### Step 2: Start OTEL Collector
```bash
docker-compose up otel-collector -d
```

### Step 3: Run the Go Application
```bash
./run-with-collector.sh
```

### Step 4: Generate Test Data
```bash
./test-app.sh
```

## Manual Setup

### 1. Install Dependencies
```bash
go mod tidy
```

### 2. Set Environment Variables
```bash
export SERVICE_NAME="go-otel-demo"
export INSECURE_MODE="true"
export OTEL_EXPORTER_OTLP_ENDPOINT="localhost:4317"
```

### 3. Start OTEL Collector
```bash
docker run -d \
  --name otel-collector \
  -p 4317:4317 \
  -p 4318:4318 \
  -v $(pwd)/otel-collector-config.yaml:/etc/otel-collector-config.yaml \
  otel/opentelemetry-collector-contrib:latest \
  --config=/etc/otel-collector-config.yaml
```

### 4. Run Application
```bash
go run main.go
```

## Verification

### Check OTEL Collector Logs
```bash
docker logs otel-collector
```

### Check Application
Visit: http://localhost:8080/health

### Generate Sample Data
```bash
# Health checks
curl http://localhost:8080/health

# Generate errors
curl http://localhost:8080/error

# Cart operations
curl "http://localhost:8080/cart/add?user_id=user1&item=laptop"
curl "http://localhost:8080/cart/get?user_id=user1"
```

## Data Flow

1. **Go Application** sends traces and metrics to `localhost:4317` (OTEL Collector)
2. **OTEL Collector** receives data and forwards to SigNoz cloud
3. **SigNoz Cloud** displays metrics and traces in dashboards

## Troubleshooting

### OTEL Collector Issues
```bash
# Check if collector is running
docker ps | grep otel-collector

# Check collector logs
docker logs otel-collector

# Restart collector
docker-compose restart otel-collector
```

### No Data in SigNoz?
1. Verify OTEL collector configuration file
2. Check your SigNoz access token
3. Ensure correct region endpoint
4. Check network connectivity

### Connection Issues
```bash
# Test collector connectivity
telnet localhost 4317

# Check application logs for OTEL errors
go run main.go | grep -i otel
```

## What You'll See in SigNoz

### Metrics Dashboard
- `http_errors_total` - Error counts by endpoint
- `http_request_duration_seconds` - Request latency histograms
- `cart_items_count` - Real-time cart sizes

### Traces Dashboard
- Request traces showing latency breakdown
- Error traces with detailed context
- Cart operation traces with user information

## Next Steps

1. Create custom dashboards in SigNoz
2. Set up alerts for high error rates
3. Monitor cart abandonment patterns
4. Add more business-specific metrics and traces
