package messaging

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/indrasaputra/toggle/entity"

	"github.com/hibiken/asynq"

	"github.com/alicebob/miniredis/v2"
	"github.com/indrasaputra/toggle/internal/config"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
	"github.com/stretchr/testify/assert"
)

type RedisPublisherExecutor struct {
	publisher *RedisPublisher
	server    *miniredis.Miniredis
}

func TestNewRedisPublisher(t *testing.T) {
	t.Run("successfully create redis publisher", func(t *testing.T) {
		exec := createRedisPublisherExecutor()
		defer exec.server.Close()
		assert.NotNil(t, exec.publisher)
	})
}

func TestRedisPublisher_Publish(t *testing.T) {
	ctx := context.Background()

	t.Run("success publish event to redis queue", func(t *testing.T) {
		exec := createRedisPublisherExecutor()
		defer exec.server.Close()

		err := exec.publisher.Publish(ctx, &togglev1.ToggleEvent{})

		assert.Nil(t, err)
	})
}

type RedisSubscriberExecutor struct {
	subscriber *RedisSubscriber
	server     *miniredis.Miniredis
}

func TestNewRedisSubscriber(t *testing.T) {
	t.Run("successfully create redis subscriber", func(t *testing.T) {
		exec := createRedisSubscriberExecutor()
		defer exec.server.Close()
		assert.NotNil(t, exec.subscriber)
	})
}

func TestRedisSubscriber_Subscribe(t *testing.T) {
	ctx := context.Background()
	exec := createRedisSubscriberExecutor()
	defer exec.server.Close()

	t.Run("success execute fn", func(t *testing.T) {
		fn := func(*togglev1.ToggleEvent) error {
			return nil
		}
		runCh := make(chan struct{})
		go func() {
			close(runCh)
			_ = exec.subscriber.Subscribe(ctx, fn)
		}()
		<-runCh
		exec.subscriber.Stop()
	})
}

func TestRedisSubscriber_handleToggleEvent(t *testing.T) {
	ctx := context.Background()
	exec := createRedisSubscriberExecutor()
	defer exec.server.Close()

	t.Run("unmarshal task fail", func(t *testing.T) {
		fn := func(*togglev1.ToggleEvent) error {
			return nil
		}
		retFn := exec.subscriber.handleToggleEvent(fn)
		err := retFn(ctx, &asynq.Task{})

		assert.NotNil(t, err)
	})

	t.Run("unmarshal task success", func(t *testing.T) {
		fn := func(*togglev1.ToggleEvent) error {
			return nil
		}
		payload, _ := json.Marshal(entity.EventToggleCreated(&entity.Toggle{Key: "test"}))
		task := asynq.NewTask("test", payload)

		retFn := exec.subscriber.handleToggleEvent(fn)
		err := retFn(ctx, task)

		assert.Nil(t, err)
	})
}

func createRedisPublisherExecutor() *RedisPublisherExecutor {
	mr, _ := miniredis.Run()
	cfg := &config.Redis{
		Address: mr.Addr(),
	}

	return &RedisPublisherExecutor{
		publisher: NewRedisPublisher(cfg),
		server:    mr,
	}
}

func createRedisSubscriberExecutor() *RedisSubscriberExecutor {
	mr, _ := miniredis.Run()
	cfg := &config.Redis{
		Address: mr.Addr(),
	}

	return &RedisSubscriberExecutor{
		subscriber: NewRedisSubscriber(cfg),
		server:     mr,
	}
}
