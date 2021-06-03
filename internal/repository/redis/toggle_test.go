package redis_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/internal/repository/redis"
)

var (
	testCtx               = context.Background()
	testTTL               = 5 * time.Minute
	testToggleKey         = "toggle-1"
	testToggleDescription = "description"
	testToggleCreatedAt   = time.Now()
	testToggleUpdatedAt   = time.Now()
	testToggle            = &entity.Toggle{
		Key:         testToggleKey,
		IsEnabled:   false,
		Description: testToggleDescription,
		CreatedAt:   testToggleCreatedAt,
		UpdatedAt:   testToggleUpdatedAt,
	}
	testHSetInput = []string{
		"key",
		testToggleKey,
		"is_enabled",
		"false",
		"description",
		testToggleDescription,
		"created_at",
		testToggleCreatedAt.Format(time.RFC3339),
		"updated_at",
		testToggleUpdatedAt.Format(time.RFC3339),
	}
	testEmptyMapResult = make(map[string]string)
	testValidMapResult = map[string]string{
		"key":         testToggleKey,
		"is_enabled":  "true",
		"description": testToggleDescription,
		"created_at":  testToggleCreatedAt.Format(time.RFC3339),
		"updated_at":  testToggleUpdatedAt.Format(time.RFC3339),
	}
	testRedisDownMessage = "redis down"
)

type ToggleExecutor struct {
	toggle *redis.Toggle
	mock   redismock.ClientMock
}

func TestNewToggle(t *testing.T) {
	t.Run("successfully create an instance of Toggle", func(t *testing.T) {
		exec := createToggleExecutor()
		assert.NotNil(t, exec.toggle)
	})
}

func TestToggle_Set(t *testing.T) {
	t.Run("not all attributes are saved", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.mock.ExpectHSet(testToggleKey, testHSetInput).SetVal(2)
		exec.mock.ExpectExpire(testToggleKey, testTTL).SetVal(true)

		err := exec.toggle.Set(testCtx, testToggle)

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "only success to save 2 out of 5 attributes")
	})

	t.Run("redis is down", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.mock.ExpectHSet(testToggleKey, testHSetInput).SetVal(5)
		exec.mock.ExpectExpire(testToggleKey, testTTL).SetErr(errors.New(testRedisDownMessage))

		err := exec.toggle.Set(testCtx, testToggle)

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), testRedisDownMessage)
	})

	t.Run("success save res in redis hash", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.mock.ExpectHSet(testToggleKey, testHSetInput).SetVal(5)
		exec.mock.ExpectExpire(testToggleKey, testTTL).SetVal(true)

		err := exec.toggle.Set(testCtx, testToggle)

		assert.Nil(t, err)
	})
}

func TestToggle_Get(t *testing.T) {
	t.Run("redis hgetall returns not found (redis: nil)", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.mock.ExpectHGetAll(testToggleKey).SetErr(errors.New("redis: nil"))

		res, err := exec.toggle.Get(testCtx, testToggleKey)

		assert.Nil(t, err)
		assert.Nil(t, res)
	})

	t.Run("redis hgetall returns error", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.mock.ExpectHGetAll(testToggleKey).SetErr(errors.New(testRedisDownMessage))

		res, err := exec.toggle.Get(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(testRedisDownMessage), err)
		assert.Nil(t, res)
	})

	t.Run("redis hgetall returns empty hash", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.mock.ExpectHGetAll(testToggleKey).SetVal(testEmptyMapResult)

		res, err := exec.toggle.Get(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
		assert.Nil(t, res)
	})

	t.Run("toggle is_enabled is invalid", func(t *testing.T) {
		exec := createToggleExecutor()
		hash := make(map[string]string)
		hash["is_enabled"] = "no-value"
		exec.mock.ExpectHGetAll(testToggleKey).SetVal(hash)

		res, err := exec.toggle.Get(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Nil(t, res)
	})

	t.Run("toggle created_at is invalid", func(t *testing.T) {
		exec := createToggleExecutor()
		hash := make(map[string]string)
		hash["is_enabled"] = "false"
		hash["created_at"] = ""
		exec.mock.ExpectHGetAll(testToggleKey).SetVal(hash)

		res, err := exec.toggle.Get(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Nil(t, res)
	})

	t.Run("toggle updated_at is invalid", func(t *testing.T) {
		exec := createToggleExecutor()
		hash := make(map[string]string)
		hash["is_enabled"] = "false"
		hash["created_at"] = "2021-04-04T12:05:38.728727+07:00"
		hash["updated_at"] = ""
		exec.mock.ExpectHGetAll(testToggleKey).SetVal(hash)

		res, err := exec.toggle.Get(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Nil(t, res)
	})

	t.Run("success get toggle from redis hash", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.mock.ExpectHGetAll(testToggleKey).SetVal(testValidMapResult)

		res, err := exec.toggle.Get(testCtx, testToggleKey)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
}

func TestToggle_SetIsEnabled(t *testing.T) {
	t.Run("set returns error", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.mock.ExpectHSet(testToggleKey, "is_enabled", "false", "updated_at", time.Now().Format(time.RFC3339)).SetErr(errors.New(testRedisDownMessage))

		err := exec.toggle.SetIsEnabled(testCtx, testToggleKey, false)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(testRedisDownMessage), err)
	})

	t.Run("success set is_enabled field", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.mock.ExpectHSet(testToggleKey, "is_enabled", "true", "updated_at", time.Now().Format(time.RFC3339)).SetVal(0)

		err := exec.toggle.SetIsEnabled(testCtx, testToggleKey, true)

		assert.Nil(t, err)
	})
}

func TestToggle_Delete(t *testing.T) {
	t.Run("delete returns error", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.mock.ExpectDel(testToggleKey).SetErr(errors.New(testRedisDownMessage))

		err := exec.toggle.Delete(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(testRedisDownMessage), err)
	})

	t.Run("success delete", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.mock.ExpectDel(testToggleKey).SetVal(1)

		err := exec.toggle.Delete(testCtx, testToggleKey)

		assert.Nil(t, err)
	})
}

func createToggleExecutor() *ToggleExecutor {
	client, mock := redismock.NewClientMock()
	rds := redis.NewToggle(client, testTTL)
	return &ToggleExecutor{
		toggle: rds,
		mock:   mock,
	}
}
