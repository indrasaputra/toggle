package builder

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/indrasaputra/toggle/internal/config"
	"github.com/indrasaputra/toggle/internal/grpc/handler"
	"github.com/indrasaputra/toggle/internal/repository"
	"github.com/indrasaputra/toggle/internal/repository/postgres"
	"github.com/indrasaputra/toggle/internal/repository/redis"
	"github.com/indrasaputra/toggle/service"
)

// BuildToggleHandler builds toggle handler including all of its dependencies.
func BuildToggleHandler(pool *pgxpool.Pool, rdsClient goredis.Cmdable, rdsTTL time.Duration) *handler.Toggle {
	psql := postgres.NewToggle(pool)
	rds := redis.NewToggle(rdsClient, rdsTTL)

	inserterRepo := repository.NewToggleInserter(psql, rds)
	deleterRepo := repository.NewToggleDeleter(psql, rds)

	creator := service.NewToggleCreator(inserterRepo)
	deleter := service.NewToggleDeleter(deleterRepo)

	return handler.NewToggle(creator, deleter)
}

// BuildPgxPool builds a pool of pgx client.
func BuildPgxPool(cfg *config.Postgres) (*pgxpool.Pool, error) {
	connCfg := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable pool_max_conns=%s pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.MaxOpenConns,
		cfg.MaxConnLifetime,
		cfg.MaxIdleLifetime,
	)
	return pgxpool.Connect(context.Background(), connCfg)
}

// BuildRedisClient builds an instance of redis client.
func BuildRedisClient(cfg *config.Redis) (*goredis.Client, error) {
	opt := &goredis.Options{
		Addr: cfg.Address,
	}

	client := goredis.NewClient(opt)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
