package middleware

import (
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"go-otel-demo/internal/config"
	"go-otel-demo/internal/telemetry"
)

// ResponseWriter wrapper to capture status code
type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// MetricsMiddleware measures request latency, counts errors, and adds tracing
func MetricsMiddleware(cfg *config.Config, collector *telemetry.MetricsCollector) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Start a new span for tracing
			tracer := otel.Tracer(cfg.ServiceName)
			ctx, span := tracer.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.Path))
			defer span.End()

			// Add some attributes to the span
			span.SetAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.route", r.URL.Path),
			)

			start := time.Now()

			// Create a response writer wrapper to capture status code
			ww := &ResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}

			// Call the next handler with the traced context
			r = r.WithContext(ctx)
			next.ServeHTTP(ww, r)

			// Add status code to span
			span.SetAttributes(attribute.Int("http.status_code", ww.StatusCode))

			// Record latency
			duration := time.Since(start).Seconds()
			collector.LatencyRecorder.Record(ctx, duration,
				metric.WithAttributes(
					attribute.String("method", r.Method),
					attribute.String("endpoint", r.URL.Path),
					attribute.Int("status_code", ww.StatusCode),
				))

			// Count errors (4xx and 5xx status codes)
			if ww.StatusCode >= 400 {
				span.SetAttributes(attribute.Bool("error", true))
				collector.ErrorCounter.Add(ctx, 1,
					metric.WithAttributes(
						attribute.String("method", r.Method),
						attribute.String("endpoint", r.URL.Path),
						attribute.Int("status_code", ww.StatusCode),
					))
			}
		}
	}
}
