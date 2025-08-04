package config

import "os"

// Config holds all configuration for the application
type Config struct {
	ServiceName  string
	CollectorURL string
	Insecure     string
	ServerPort   string
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		ServiceName:  getEnv("SERVICE_NAME", "go-otel-demo"),
		CollectorURL: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
		Insecure:     getEnv("INSECURE_MODE", "true"),
		ServerPort:   getEnv("SERVER_PORT", ":8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
