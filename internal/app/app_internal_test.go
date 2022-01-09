package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func TestTracer(t *testing.T) {
	t.Run("success set tracer", func(t *testing.T) {
		tracer := createTracer()
		Tracer(tracer)
		assert.NotNil(t, appTracer)
	})
}

func TestGetTracer(t *testing.T) {
	expTracer := createTracer()

	t.Run("tracer is default noop tracer", func(t *testing.T) {
		appTracer = nil
		tracer := GetTracer()
		assert.NotEqual(t, expTracer, tracer)
	})

	t.Run("tracer is expected tracer", func(t *testing.T) {
		Tracer(expTracer)
		tracer := GetTracer()
		assert.Equal(t, expTracer, tracer)
	})
}

func createTracer() oteltrace.Tracer {
	tracer := tracesdk.NewTracerProvider()
	return tracer.Tracer("test")
}
