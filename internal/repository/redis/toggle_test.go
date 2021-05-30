package redis_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"

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

	t.Run("success save url in redis hash", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.mock.ExpectHSet(testToggleKey, testHSetInput).SetVal(5)
		exec.mock.ExpectExpire(testToggleKey, testTTL).SetVal(true)

		err := exec.toggle.Set(testCtx, testToggle)

		assert.Nil(t, err)
	})
}

func TestToggle_Delete(t *testing.T) {
	t.Run("delete returns error", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.mock.ExpectDel(testToggleKey).SetErr(errors.New("redis: nil"))

		err := exec.toggle.Delete(testCtx, testToggleKey)

		assert.NotNil(t, err)
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
