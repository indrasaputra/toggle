package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"

	"github.com/indrasaputra/toggle/internal/config"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
)

// AsynqPublisher is responsible to publish message to Asynq.
type AsynqPublisher struct {
	client *asynq.Client
}

// NewAsynqPublisher creates an instance of AsynqPublisher.
func NewAsynqPublisher(cfg *config.Asynq) *AsynqPublisher {
	client := asynq.NewClient(
		asynq.RedisClientOpt{
			Addr: cfg.Address,
			DB:   cfg.DBSelect,
		},
	)
	return &AsynqPublisher{client: client}
}

// Publish publishes toggle event to Asynq.
// The event will be converted to JSON.
func (ap *AsynqPublisher) Publish(ctx context.Context, event *togglev1.ToggleEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	task := asynq.NewTask(event.GetName().String(), payload)
	_, err = ap.client.Enqueue(task)
	return err
}

// AsynqSubscriber is responsible to subscribe message from Asynq.
type AsynqSubscriber struct {
	server *asynq.Server
}

// NewAsynqSubscriber creates an instance of AsynqSubscriber.
func NewAsynqSubscriber(cfg *config.Asynq) *AsynqSubscriber {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Address},
		asynq.Config{Concurrency: cfg.Concurrency},
	)
	return &AsynqSubscriber{server: server}
}

// Subscribe subscribes to a certain topic and process the incoming message using the fn parameter.
// This method is blocking.
func (as *AsynqSubscriber) Subscribe(ctx context.Context, fn func(*togglev1.ToggleEvent) error) error {
	return as.server.Run(asynq.HandlerFunc(as.handleToggleEvent(fn)))
}

// Stop stops the subscriber.
func (as *AsynqSubscriber) Stop() {
	as.server.Shutdown()
}

func (as *AsynqSubscriber) handleToggleEvent(fn func(*togglev1.ToggleEvent) error) func(ctx context.Context, t *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var event *togglev1.ToggleEvent
		if err := json.Unmarshal(t.Payload(), &event); err != nil {
			log.Printf("error unmarshal message: %v\n", err)
			return err
		}
		return fn(event)
	}
}
