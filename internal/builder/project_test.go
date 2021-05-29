package builder_test

import (
	"testing"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/internal/builder"
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
