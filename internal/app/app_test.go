package app_test

import (
	"testing"

	"github.com/indrasaputra/toggle/internal/app"
	"github.com/indrasaputra/toggle/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestInitTracer(t *testing.T) {
	t.Run("init tracer returns directly if jaeger is disabled", func(t *testing.T) {
		cfg := &config.Config{
			Jaeger: config.Jaeger{
				Enabled: false,
			},
		}
		err := app.InitTracer(cfg)
		assert.Nil(t, err)
	})

	t.Run("success init tracer and set it to app", func(t *testing.T) {
		cfg := &config.Config{
			Jaeger: config.Jaeger{
				Enabled: true,
			},
			ServiceName: "svc",
			AppEnv:      "test",
		}
		err := app.InitTracer(cfg)
		assert.Nil(t, err)
	})
}
