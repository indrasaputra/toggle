package messaging

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/internal/config"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
)

type AsynqPublisherExecutor struct {
	publisher *AsynqPublisher
	server    *miniredis.Miniredis
}

func TestNewAsynqPublisher(t *testing.T) {
	t.Run("successfully create redis publisher", func(t *testing.T) {
		exec := createAsynqPublisherExecutor()
		defer exec.server.Close()
		assert.NotNil(t, exec.publisher)
	})
}

func TestAsynqPublisher_Publish(t *testing.T) {
	ctx := context.Background()

	t.Run("success publish event to redis queue", func(t *testing.T) {
		exec := createAsynqPublisherExecutor()
		defer exec.server.Close()

		err := exec.publisher.Publish(ctx, &togglev1.ToggleEvent{})

		assert.Nil(t, err)
	})
}

type AsynqSubscriberExecutor struct {
	subscriber *AsynqSubscriber
	server     *miniredis.Miniredis
}

func TestNewAsynqSubscriber(t *testing.T) {
	t.Run("successfully create redis subscriber", func(t *testing.T) {
		exec := createAsynqSubscriberExecutor()
		defer exec.server.Close()
		assert.NotNil(t, exec.subscriber)
	})
}

func TestAsynqSubscriber_Subscribe(t *testing.T) {
	ctx := context.Background()
	exec := createAsynqSubscriberExecutor()
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

func TestAsynqSubscriber_handleToggleEvent(t *testing.T) {
	ctx := context.Background()
	exec := createAsynqSubscriberExecutor()
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

func createAsynqPublisherExecutor() *AsynqPublisherExecutor {
	mr, _ := miniredis.Run()
	cfg := &config.Asynq{
		Address: mr.Addr(),
	}

	return &AsynqPublisherExecutor{
		publisher: NewAsynqPublisher(cfg),
		server:    mr,
	}
}

func createAsynqSubscriberExecutor() *AsynqSubscriberExecutor {
	mr, _ := miniredis.Run()
	cfg := &config.Asynq{
		Address: mr.Addr(),
	}

	return &AsynqSubscriberExecutor{
		subscriber: NewAsynqSubscriber(cfg),
		server:     mr,
	}
}
