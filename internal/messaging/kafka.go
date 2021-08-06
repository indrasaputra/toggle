package messaging

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"

	"github.com/indrasaputra/toggle/entity"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
)

// KafkaPublisher is responsible to publish message to Kafka.
type KafkaPublisher struct {
	writer Writer
}

// NewKafkaPublisher creates an instance of KafkaPublisher.
// It sets 10 seconds as default deadline.
func NewKafkaPublisher(writer Writer) *KafkaPublisher {
	return &KafkaPublisher{writer: writer}
}

// Publish publishes toggle event to Kafka.
// The event will be converted to JSON.
func (kp *KafkaPublisher) Publish(ctx context.Context, event *togglev1.EventToggle) error {
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
