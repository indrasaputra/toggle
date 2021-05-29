package builder

import (
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"

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

	creator := service.NewToggleCreator(inserterRepo)

	return handler.NewToggle(creator)
}
