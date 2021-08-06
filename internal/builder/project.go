package builder

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/segmentio/kafka-go"

	"github.com/indrasaputra/toggle/internal/config"
	decorservice "github.com/indrasaputra/toggle/internal/decorator/service"
	"github.com/indrasaputra/toggle/internal/grpc/handler"
	"github.com/indrasaputra/toggle/internal/messaging"
	"github.com/indrasaputra/toggle/internal/repository"
	"github.com/indrasaputra/toggle/internal/repository/postgres"
	"github.com/indrasaputra/toggle/internal/repository/redis"
	"github.com/indrasaputra/toggle/service"
)

// BuildToggleCommandHandler builds toggle command handler including all of its dependencies.
func BuildToggleCommandHandler(pool *pgxpool.Pool, rdsClient goredis.Cmdable, rdsTTL time.Duration, kafkaWriter *kafka.Writer) *handler.ToggleCommand {
	psql := postgres.NewToggle(pool)
	rds := redis.NewToggle(rdsClient, rdsTTL)
	writer := messaging.NewKafkaPublisher(kafkaWriter)

	inserterRepo := repository.NewToggleInserter(psql, rds)
	updaterRepo := repository.NewToggleUpdater(psql, rds)
	deleterRepo := repository.NewToggleDeleter(psql, rds)

	creator := service.NewToggleCreator(inserterRepo, writer)
	enabler := service.NewToggleEnabler(updaterRepo, writer)
	disabler := service.NewToggleDisabler(updaterRepo, writer)
	deleter := service.NewToggleDeleter(deleterRepo, writer)

	decor := decorservice.NewTracing(creator, nil, enabler, disabler, deleter)

	// this one is not a good example, but I let it be for now since I compose all services in one decorator.
	return handler.NewToggleCommand(decor, decor, decor, decor)
}

// BuildToggleQueryHandler builds toggle query handler including all of its dependencies.
func BuildToggleQueryHandler(pool *pgxpool.Pool, rdsClient goredis.Cmdable, rdsTTL time.Duration) *handler.ToggleQuery {
	psql := postgres.NewToggle(pool)
	rds := redis.NewToggle(rdsClient, rdsTTL)

	getterRepo := repository.NewToggleGetter(psql, rds)

	getter := service.NewToggleGetter(getterRepo)

	decor := decorservice.NewTracing(nil, getter, nil, nil, nil)

	// this one is not a good example, but I let it be for now since I compose all services in one decorator.
	return handler.NewToggleQuery(decor)
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
	client.AddHook(redis.NewHookTracing())

	return client, nil
}

// BuildKafkaWriter builds an instance of kafka writer.
func BuildKafkaWriter(cfg *config.Kafka) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(cfg.Address),
		Topic:        cfg.Topic,
		MaxAttempts:  cfg.MaxAttempts,
		BatchSize:    cfg.BatchSize,
		BatchTimeout: time.Duration(cfg.BatchTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	}
}
