package app_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/internal/app"
	"github.com/indrasaputra/toggle/internal/config"
)

func TestInitTracer(t *testing.T) {
	t.Run("init tracer returns directly if jaeger is disabled", func(t *testing.T) {
		cfg := &config.Config{
			Jaeger: config.Jaeger{
				Enabled: false,
			},
		}

		prov, err := app.InitTracer(cfg)

		assert.Nil(t, err)
		assert.Nil(t, prov)
	})

	t.Run("success init tracer and set it to app", func(t *testing.T) {
		cfg := &config.Config{
			Jaeger: config.Jaeger{
				Enabled: true,
			},
			ServiceName: "svc",
			AppEnv:      "test",
		}
		prov, err := app.InitTracer(cfg)
		defer func() { _ = prov.Shutdown(context.Background()) }()

		assert.Nil(t, err)
		assert.NotNil(t, prov)
	})
}
