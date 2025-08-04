package telemetry

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc/credentials"

	"go-otel-demo/internal/config"
	"go-otel-demo/internal/models"
)

// MetricsCollector holds all the metrics instruments and cart data
type MetricsCollector struct {
	ErrorCounter    metric.Int64Counter
	LatencyRecorder metric.Float64Histogram
	CartSizeGauge   metric.Int64ObservableGauge
	Carts           map[string]*models.ShoppingCart
}

// Setup initializes both tracing and metrics
func Setup(cfg *config.Config) (*MetricsCollector, func(context.Context) error, error) {
	// Initialize tracer
	tracerCleanup, err := initTracer(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	// Initialize metrics
	metricsCollector, err := initMetrics(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	return metricsCollector, tracerCleanup, nil
}

// initTracer initializes the tracer following SigNoz pattern
func initTracer(cfg *config.Config) (func(context.Context) error, error) {
	var secureOption otlptracegrpc.Option

	if strings.ToLower(cfg.Insecure) == "false" || cfg.Insecure == "0" || strings.ToLower(cfg.Insecure) == "f" {
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(cfg.CollectorURL),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceName),
			attribute.String("library.language", "go"),
			attribute.String("service.version", "1.0.0"),
			attribute.String("environment", "development"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("could not set resources: %w", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown, nil
}

// initMetrics initializes the metrics
func initMetrics(cfg *config.Config) (*MetricsCollector, error) {
	var secureOption otlpmetricgrpc.Option

	if strings.ToLower(cfg.Insecure) == "false" || cfg.Insecure == "0" || strings.ToLower(cfg.Insecure) == "f" {
		secureOption = otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlpmetricgrpc.WithInsecure()
	}

	// Create resource
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("environment", "development"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create OTLP gRPC exporter for metrics
	exporter, err := otlpmetricgrpc.New(context.Background(),
		secureOption,
		otlpmetricgrpc.WithEndpoint(cfg.CollectorURL),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP metrics exporter: %w", err)
	}

	// Create meter provider
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter,
			sdkmetric.WithInterval(10*time.Second),
		)),
	)

	// Set global meter provider
	otel.SetMeterProvider(meterProvider)

	// Get meter
	meter := otel.Meter("go-otel-demo")

	// Initialize metrics collector
	metricsCollector := &MetricsCollector{
		Carts: make(map[string]*models.ShoppingCart),
	}

	// Create Counter for error requests
	metricsCollector.ErrorCounter, err = meter.Int64Counter(
		"http_errors_total",
		metric.WithDescription("Total number of HTTP error requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create error counter: %w", err)
	}

	// Create Histogram for request latency
	metricsCollector.LatencyRecorder, err = meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("HTTP request latency in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create latency histogram: %w", err)
	}

	// Create Gauge for cart size
	metricsCollector.CartSizeGauge, err = meter.Int64ObservableGauge(
		"cart_items_count",
		metric.WithDescription("Number of items currently in user carts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create cart size gauge: %w", err)
	}

	// Register callback for gauge
	_, err = meter.RegisterCallback(
		func(ctx context.Context, observer metric.Observer) error {
			for userID, cart := range metricsCollector.Carts {
				observer.ObserveInt64(metricsCollector.CartSizeGauge, int64(len(cart.Items)),
					metric.WithAttributes(attribute.String("user_id", userID)))
			}
			return nil
		},
		metricsCollector.CartSizeGauge,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register gauge callback: %w", err)
	}

	log.Println("Telemetry initialized successfully")
	return metricsCollector, nil
}
