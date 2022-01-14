package builder

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/extra/redisotel"
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

var (
	postgresConnFormat  = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable pool_max_conns=%s pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s sslmode=%s"
	cockroachConnFormat = postgresConnFormat + " sslrootcert=%s options=%s"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	PgxPool     *pgxpool.Pool
	RedisClient goredis.Cmdable
	KafkaWriter *kafka.Writer
	Config      *config.Config
}

// BuildToggleCommandHandler builds toggle command handler including all of its dependencies.
func BuildToggleCommandHandler(dep *Dependency) *handler.ToggleCommand {
	psql := postgres.NewToggle(dep.PgxPool)
	rds := redis.NewToggle(dep.RedisClient, time.Duration(dep.Config.Redis.TTL)*time.Minute)
	publisher := messaging.NewRedisPublisher(&dep.Config.Redis)

	inserterRepo := repository.NewToggleInserter(psql, rds)
	updaterRepo := repository.NewToggleUpdater(psql, rds)
	deleterRepo := repository.NewToggleDeleter(psql, rds)

	creator := service.NewToggleCreator(inserterRepo, publisher)
	enabler := service.NewToggleEnabler(updaterRepo, publisher)
	disabler := service.NewToggleDisabler(updaterRepo, publisher)
	deleter := service.NewToggleDeleter(deleterRepo, publisher)

	decor := decorservice.NewTracing(creator, nil, enabler, disabler, deleter)

	// this one is not a good example, but I let it be for now since I compose all services in one decorator.
	return handler.NewToggleCommand(decor, decor, decor, decor)
}

// BuildToggleQueryHandler builds toggle query handler including all of its dependencies.
func BuildToggleQueryHandler(dep *Dependency) *handler.ToggleQuery {
	psql := postgres.NewToggle(dep.PgxPool)
	rds := redis.NewToggle(dep.RedisClient, time.Duration(dep.Config.Redis.TTL)*time.Minute)

	getterRepo := repository.NewToggleGetter(psql, rds)

	getter := service.NewToggleGetter(getterRepo)

	decor := decorservice.NewTracing(nil, getter, nil, nil, nil)

	// this one is not a good example, but I let it be for now since I compose all services in one decorator.
	return handler.NewToggleQuery(decor)
}

// BuildPostgrePgxPool builds a pool of pgx client.
func BuildPostgrePgxPool(cfg *config.Postgres) (*pgxpool.Pool, error) {
	connCfg := fmt.Sprintf(postgresConnFormat,
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.MaxOpenConns,
		cfg.MaxConnLifetime,
		cfg.MaxIdleLifetime,
		cfg.SSLMode,
	)
	return pgxpool.Connect(context.Background(), connCfg)
}

// BuildCockroachPgxPool builds a pool of cockroachdb client using pgx.
func BuildCockroachPgxPool(cfg *config.CockroachDB) (*pgxpool.Pool, error) {
	connCfg := fmt.Sprintf(cockroachConnFormat,
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.MaxOpenConns,
		cfg.MaxConnLifetime,
		cfg.MaxIdleLifetime,
		cfg.SSLMode,
		cfg.SSLRootCert,
		cfg.Options,
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

	client.AddHook(redisotel.TracingHook{})

	return client, nil
}

// BuildKafkaWriter builds an instance of kafka writer.
//
// Currently, the writer is having issue on auto-creating topic.
// Follow more on https://github.com/segmentio/kafka-go/issues/683 and https://github.com/segmentio/kafka-go/pull/700.
// Once the PR is submitted and the package is updated, this configuration should work.
func BuildKafkaWriter(cfg *config.Kafka) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(cfg.Address),
		Topic:        cfg.Topic,
		MaxAttempts:  cfg.MaxAttempts,
		BatchSize:    cfg.BatchSize,
		BatchTimeout: time.Duration(cfg.BatchTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		Async:        cfg.WriteAsync,
	}
}
