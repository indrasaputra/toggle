package messaging

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// Writer defines a little interface for Kafka writer/publisher functionality.
// Since in the real implementation we can use kafka.Writer,
// this interface exists mostly for testing purpose.
type Writer interface {
	WriteMessages(ctx context.Context, messages ...kafka.Message) error
	Close() error
}
