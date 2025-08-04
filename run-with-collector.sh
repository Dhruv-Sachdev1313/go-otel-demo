#!/bin/bash

echo "Setting up Go OpenTelemetry Demo with OTEL Collector"
echo "===================================================="

# Set environment variables
export SERVICE_NAME="go-otel-demo"
export INSECURE_MODE="true"
export OTEL_EXPORTER_OTLP_ENDPOINT="localhost:4317"
export SERVER_PORT=":8080"

echo "Environment Variables Set:"
echo "SERVICE_NAME=$SERVICE_NAME"
echo "INSECURE_MODE=$INSECURE_MODE"
echo "OTEL_EXPORTER_OTLP_ENDPOINT=$OTEL_EXPORTER_OTLP_ENDPOINT"
echo "SERVER_PORT=$SERVER_PORT"
echo ""

echo "Starting the Go application..."
go run main.go
