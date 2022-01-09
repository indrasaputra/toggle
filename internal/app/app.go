package app

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/indrasaputra/toggle/internal/config"
)

const (
	appNoopTracer = "noop-tracer"
)

var (
	appTracer oteltrace.Tracer
)

// InitTracer initializes tracer and set it as default tracer.
// The created tracer can be accessed using GetTracer method.
// Tracer will only be created if jaeger config is enabled.
func InitTracer(cfg *config.Config) (*tracesdk.TracerProvider, error) {
	if !cfg.Jaeger.Enabled {
		return nil, nil
	}

	exporter, err := jaeger.New(jaeger.WithAgentEndpoint(
		jaeger.WithAgentHost(cfg.Jaeger.Host),
		jaeger.WithAgentPort(cfg.Jaeger.Port),
	))
	if err != nil {
		return nil, err
	}

	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			attribute.String("environment", cfg.AppEnv),
		)),
	)

	otel.SetTracerProvider(tracerProvider)
	Tracer(tracerProvider.Tracer(cfg.ServiceName))
	return tracerProvider, nil
}

// Tracer sets the global app tracer.
func Tracer(tracer oteltrace.Tracer) {
	appTracer = tracer
}

// GetTracer gets the global app tracer.
func GetTracer() oteltrace.Tracer {
	if appTracer == nil {
		return otel.GetTracerProvider().Tracer(appNoopTracer)
	}
	return appTracer
}
