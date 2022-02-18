package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"

	"github.com/indrasaputra/toggle/internal/config"
	togglev1 "github.com/indrasaputra/toggle/proto/toggle/v1"
)

// RedisPublisher is responsible to publish message to Redis.
// It uses asynq client.
type RedisPublisher struct {
	client *asynq.Client
}

// NewRedisPublisher creates an instance of RedisPublisher.
func NewRedisPublisher(cfg *config.Redis) *RedisPublisher {
	client := asynq.NewClient(
		asynq.RedisClientOpt{
			Addr: cfg.Address,
			DB:   cfg.DBSelect,
		},
	)
	return &RedisPublisher{client: client}
}

// Publish publishes toggle event to Redis.
// The event will be converted to JSON.
func (rp *RedisPublisher) Publish(ctx context.Context, event *togglev1.ToggleEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	task := asynq.NewTask(event.GetName().String(), payload)
	_, err = rp.client.Enqueue(task)
	return err
}

// RedisSubscriber is responsible to subscribe message from Redis.
type RedisSubscriber struct {
	server *asynq.Server
}

// NewRedisSubscriber creates an instance of RedisSubscriber.
func NewRedisSubscriber(cfg *config.Redis) *RedisSubscriber {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Address},
		asynq.Config{Concurrency: cfg.Concurrency},
	)
	return &RedisSubscriber{server: server}
}

// Subscribe subscribes to a certain topic and process the incoming message using the fn parameter.
// This method is blocking.
func (rs *RedisSubscriber) Subscribe(ctx context.Context, fn func(*togglev1.ToggleEvent) error) error {
	return rs.server.Run(asynq.HandlerFunc(rs.handleToggleEvent(fn)))
}

// Stop stops the subscriber.
func (rs *RedisSubscriber) Stop() {
	rs.server.Shutdown()
}

func (rs *RedisSubscriber) handleToggleEvent(fn func(*togglev1.ToggleEvent) error) func(ctx context.Context, t *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var event *togglev1.ToggleEvent
		if err := json.Unmarshal(t.Payload(), &event); err != nil {
			log.Printf("error unmarshal message: %v\n", err)
			return err
		}
		return fn(event)
	}
}
