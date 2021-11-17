package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/indrasaputra/toggle/entity"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
)

// KafkaPublisher is responsible to publish message to Kafka.
type KafkaPublisher struct {
	writer Writer
}

// NewKafkaPublisher creates an instance of KafkaPublisher.
func NewKafkaPublisher(writer Writer) *KafkaPublisher {
	return &KafkaPublisher{writer: writer}
}

// Publish publishes toggle event to Kafka.
// The event will be converted to JSON.
func (kp *KafkaPublisher) Publish(ctx context.Context, event *togglev1.ToggleEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if event.GetToggle() == nil {
		return entity.ErrEmptyToggle()
	}

	msg := kafka.Message{
		Key:   []byte(event.GetToggle().Key),
		Value: data,
	}
	if err := kp.writer.WriteMessages(ctx, msg); err != nil {
		return entity.ErrInternal(err.Error())
	}
	return nil
}

// KafkaSubscriber is responsible to subscribe message from Kafka.
type KafkaSubscriber struct {
	reader Reader
}

// NewKafkaSubscriber creates an instance of KafkaSubscriber.
func NewKafkaSubscriber(reader Reader) *KafkaSubscriber {
	return &KafkaSubscriber{reader: reader}
}

// Subscribe subscribes to a certain topic and process the incoming message using the fn parameter.
// This method is blocking.
func (ks *KafkaSubscriber) Subscribe(ctx context.Context, fn func(*togglev1.ToggleEvent) error) error {
	for {
		msg, err := ks.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		var event *togglev1.ToggleEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("error unmarshal message: %v\n", err)
			continue
		}
		if err := fn(event); err != nil {
			log.Printf("error process event in fn: %v\n", err)
		}
	}
}
