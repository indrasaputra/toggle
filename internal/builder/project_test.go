package builder_test

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/internal/builder"
	"github.com/indrasaputra/toggle/internal/config"
)

func TestBuildToggleCommandHandler(t *testing.T) {
	t.Run("success create toggle command handler", func(t *testing.T) {
		dep := &builder.Dependency{
			PgxPool:     &pgxpool.Pool{},
			RedisClient: &goredis.Client{},
			KafkaWriter: &kafka.Writer{},
			Config: &config.Config{
				Redis: config.Redis{},
			},
		}

		handler := builder.BuildToggleCommandHandler(dep)

		assert.NotNil(t, handler)
	})
}

func TestBuildToggleHandler(t *testing.T) {
	t.Run("success create toggle query handler", func(t *testing.T) {
		dep := &builder.Dependency{
			PgxPool:     &pgxpool.Pool{},
			RedisClient: &goredis.Client{},
			Config: &config.Config{
				Redis: config.Redis{},
			},
		}

		handler := builder.BuildToggleQueryHandler(dep)

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
		SSLMode:         "disable",
	}

	t.Run("fail build postgres pgxpool client", func(t *testing.T) {
		client, err := builder.BuildPgxPool(cfg)

		assert.NotNil(t, err)
		assert.Nil(t, client)
	})
}

func TestBuildCockroachPgxPool(t *testing.T) {
	cfg := &config.CockroachDB{
		Host:            "localhost",
		Port:            "5432",
		Name:            "guru",
		User:            "user",
		Password:        "password",
		MaxOpenConns:    "10",
		MaxConnLifetime: "10m",
		MaxIdleLifetime: "5m",
		SSLMode:         "verify-full",
		SSLRootCert:     "/root",
		Options:         "--cluster%3Dtoggle",
	}

	t.Run("fail build cockroachdb pgx pool client", func(t *testing.T) {
		client, err := builder.BuildCockroachPgxPool(cfg)

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

func TestBuildKafkaWriter(t *testing.T) {
	cfg := &config.Kafka{
		Address:      "localhost:9092",
		Topic:        "toggle",
		MaxAttempts:  10,
		BatchSize:    100,
		BatchTimeout: 1,
		WriteTimeout: 10,
	}

	t.Run("success build kafka writer", func(t *testing.T) {
		writer := builder.BuildKafkaWriter(cfg)

		assert.NotNil(t, writer)
	})
}
