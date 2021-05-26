package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	goredis "github.com/go-redis/redis/v8"

	"github.com/indrasaputra/toggle/entity"
)

var (
	attributes        = []string{"key", "is_enabled", "description", "created_at", "updated_at"}
	numberOfAttribute = len(attributes)
)

// Toggle is responsible to connect toggle entity with toggle data structure in Redis.
// It uses https://github.com/go-redis/redis.
type Toggle struct {
	client goredis.Cmdable
	ttl    time.Duration
}

// NewToggle creates an instance of Toggle.
func NewToggle(client goredis.Cmdable, ttl time.Duration) *Toggle {
	return &Toggle{
		client: client,
		ttl:    ttl,
	}
}

// SetClient sets go-redis client.
// This method is specifically created for testing purposes.
func (t *Toggle) SetClient(client goredis.Cmdable) {
	t.client = client
}

// SetIfNotExists sets the toggle in redis using SETNX.
// It only sets the toggle for a certain time. It is set in ttl parameter in constructor.
func (t *Toggle) SetIfNotExists(ctx context.Context, toggle *entity.Toggle) error {
	hash := createToggleHash(toggle)

	pipe := t.client.Pipeline()
	res := pipe.HSet(ctx, toggle.Key, hash)
	pipe.Expire(ctx, toggle.Key, t.ttl)
	_, err := pipe.Exec(ctx)

	if int(res.Val()) != numberOfAttribute {
		return entity.ErrInternal(fmt.Sprintf("only success to save %d out of %d attributes", res.Val(), numberOfAttribute))
	}
	if err != nil {
		return entity.ErrInternal(err.Error())
	}
	return nil
}

func createToggleHash(toggle *entity.Toggle) []string {
	return []string{
		"key",
		toggle.Key,
		"is_enabled",
		strconv.FormatBool(toggle.IsEnabled),
		"description",
		toggle.Description,
		"created_at",
		toggle.CreatedAt.Format(time.RFC3339),
		"updated_at",
		toggle.UpdatedAt.Format(time.RFC3339),
	}
}
