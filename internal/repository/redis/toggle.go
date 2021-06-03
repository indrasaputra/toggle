package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	goredis "github.com/go-redis/redis/v8"

	"github.com/indrasaputra/toggle/entity"
)

const (
	redisNotFound = "redis: nil"
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

// Set sets the toggle in redis using hash (https://redis.io/commands/hset)..
// It only sets the toggle for a certain time. It is set in ttl parameter in constructor.
func (t *Toggle) Set(ctx context.Context, toggle *entity.Toggle) error {
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

// Get gets a toggle in cache.
// It only returns error of there is error in the system or toggle value can't be processed.
// If the data can't be found but the system is fine, it returns nil.
func (t *Toggle) Get(ctx context.Context, key string) (*entity.Toggle, error) {
	res, err := t.client.HGetAll(ctx, key).Result()
	if err != nil && err.Error() == redisNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, entity.ErrInternal(err.Error())
	}
	if len(res) == 0 {
		return nil, entity.ErrNotFound()
	}
	return createToggleFromHash(res)
}

// Delete deletes a toggle from redis.
// It doesn't return error if toggle doesn't exist.
func (t *Toggle) Delete(ctx context.Context, key string) error {
	return t.client.Del(ctx, key).Err()
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

func createToggleFromHash(hash map[string]string) (*entity.Toggle, error) {
	toggle := &entity.Toggle{}
	var err error

	toggle.Key = hash["key"]
	toggle.IsEnabled, err = strconv.ParseBool(hash["is_enabled"])
	if err != nil {
		return nil, entity.ErrInternal(err.Error())
	}
	toggle.Description = hash["description"]
	toggle.CreatedAt, err = time.Parse(time.RFC3339, hash["created_at"])
	if err != nil {
		return nil, entity.ErrInternal(err.Error())
	}
	toggle.UpdatedAt, err = time.Parse(time.RFC3339, hash["updated_at"])
	if err != nil {
		return nil, entity.ErrInternal(err.Error())
	}

	return toggle, nil
}
