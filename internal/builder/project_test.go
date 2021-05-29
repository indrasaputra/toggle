package builder_test

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/internal/builder"
	"github.com/indrasaputra/toggle/internal/config"
)

func TestBuildToggleHandler(t *testing.T) {
	t.Run("success create toggle handler", func(t *testing.T) {
		psql := &pgxpool.Pool{}
		rds := &goredis.Client{}
		ttl := 5 * time.Minute

		handler := builder.BuildToggleHandler(psql, rds, ttl)

		assert.NotNil(t, handler)
	})
}

func TestBuildPgxPool(t *testing.T) {
	cfg := &config.Postgres{
		Host:            "localhost",
		Port:            "5432",
		Name:            "guru",
		User:            "user",
		Password:        "password",
		MaxOpenConns:    "10",
		MaxConnLifetime: "10m",
		MaxIdleLifetime: "5m",
	}

	t.Run("fail build sql client", func(t *testing.T) {
		client, err := builder.BuildPgxPool(cfg)

		assert.NotNil(t, err)
		assert.Nil(t, client)
	})
}

func TestBuildRedisClient(t *testing.T) {
	t.Run("fail create redis client", func(t *testing.T) {
		server, _ := miniredis.Run()

		cfg := &config.Redis{
			Address: server.Addr(),
		}

		server.Close()
		client, err := builder.BuildRedisClient(cfg)

		assert.NotNil(t, err)
		assert.Nil(t, client)
	})

	t.Run("success create redis client", func(t *testing.T) {
		server, _ := miniredis.Run()
		defer server.Close()

		cfg := &config.Redis{
			Address: server.Addr(),
		}

		client, err := builder.BuildRedisClient(cfg)

		assert.Nil(t, err)
		assert.NotNil(t, client)
	})
}
