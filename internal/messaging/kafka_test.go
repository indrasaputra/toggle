package messaging_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/internal/messaging"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
	mock_messaging "github.com/indrasaputra/toggle/test/mock/messaging"
)

var (
	testCtx   = context.Background()
	errReturn = errors.New("error")
)

type KafkaPublisherExecutor struct {
	publisher *messaging.KafkaPublisher
	writer    *mock_messaging.MockWriter
}

func TestNewKafkaPublisher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of KafkaPublisher", func(t *testing.T) {
		exec := createKafkaPublisherExecutor(ctrl)
		assert.NotNil(t, exec.publisher)
	})
}

func TestKafkaPublisher_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("fail to encode event", func(t *testing.T) {
		exec := createKafkaPublisherExecutor(ctrl)
		err := exec.publisher.Publish(testCtx, nil)
		assert.NotNil(t, err)
	})

	t.Run("event doesn't have toggle", func(t *testing.T) {
		exec := createKafkaPublisherExecutor(ctrl)
		err := exec.publisher.Publish(testCtx, &togglev1.ToggleEvent{})
		assert.NotNil(t, err)
	})

	t.Run("fail write message", func(t *testing.T) {
		exec := createKafkaPublisherExecutor(ctrl)
		exec.writer.EXPECT().WriteMessages(testCtx, gomock.Any()).Return(errReturn)

		err := exec.publisher.Publish(testCtx, &togglev1.ToggleEvent{Toggle: &togglev1.Toggle{}})

		assert.NotNil(t, err)
	})

	t.Run("success write message", func(t *testing.T) {
		exec := createKafkaPublisherExecutor(ctrl)
		exec.writer.EXPECT().WriteMessages(testCtx, gomock.Any()).Return(nil)

		err := exec.publisher.Publish(testCtx, &togglev1.ToggleEvent{Toggle: &togglev1.Toggle{}})

		assert.Nil(t, err)
	})
}

func createKafkaPublisherExecutor(ctrl *gomock.Controller) *KafkaPublisherExecutor {
	w := mock_messaging.NewMockWriter(ctrl)
	p := messaging.NewKafkaPublisher(w)
	return &KafkaPublisherExecutor{
		publisher: p,
		writer:    w,
	}
}

type KafkaSubscriberExecutor struct {
	subscriber *messaging.KafkaSubscriber
	reader     *mock_messaging.MockReader
}

func TestNewKafkaSubscriber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of KafkaSubscriber", func(t *testing.T) {
		exec := createKafkaSubscriberExecutor(ctrl)
		assert.NotNil(t, exec.subscriber)
	})
}

func TestKafkaSubscriber_Subscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("read message returns error directly", func(t *testing.T) {
		exec := createKafkaSubscriberExecutor(ctrl)
		exec.reader.EXPECT().ReadMessage(testCtx).Return(kafka.Message{}, errReturn)

		err := exec.subscriber.Subscribe(testCtx, processor(nil))

		assert.NotNil(t, err)
	})

	t.Run("error unmarshal value", func(t *testing.T) {
		exec := createKafkaSubscriberExecutor(ctrl)
		b, _ := json.Marshal("unknown")
		msg := kafka.Message{Value: b}
		exec.reader.EXPECT().ReadMessage(testCtx).Return(msg, nil)
		exec.reader.EXPECT().ReadMessage(testCtx).Return(kafka.Message{}, errReturn)

		err := exec.subscriber.Subscribe(testCtx, processor(nil))

		assert.NotNil(t, err)
	})

	t.Run("error processing message", func(t *testing.T) {
		exec := createKafkaSubscriberExecutor(ctrl)
		b, _ := json.Marshal(&togglev1.ToggleEvent{})
		msg := kafka.Message{Value: b}
		exec.reader.EXPECT().ReadMessage(testCtx).Return(msg, nil)
		exec.reader.EXPECT().ReadMessage(testCtx).Return(kafka.Message{}, errReturn)

		err := exec.subscriber.Subscribe(testCtx, processor(errReturn))

		assert.NotNil(t, err)
	})

	t.Run("success processing message", func(t *testing.T) {
		exec := createKafkaSubscriberExecutor(ctrl)
		b, _ := json.Marshal(&togglev1.ToggleEvent{})
		msg := kafka.Message{Value: b}
		exec.reader.EXPECT().ReadMessage(testCtx).Return(msg, nil)
		exec.reader.EXPECT().ReadMessage(testCtx).Return(kafka.Message{}, errReturn)

		err := exec.subscriber.Subscribe(testCtx, processor(nil))

		assert.NotNil(t, err)
	})
}

func processor(err error) func(*togglev1.ToggleEvent) error {
	if err != nil {
		return func(*togglev1.ToggleEvent) error { return err }
	}
	return func(*togglev1.ToggleEvent) error { return nil }
}

func createKafkaSubscriberExecutor(ctrl *gomock.Controller) *KafkaSubscriberExecutor {
	r := mock_messaging.NewMockReader(ctrl)
	s := messaging.NewKafkaSubscriber(r)
	return &KafkaSubscriberExecutor{
		subscriber: s,
		reader:     r,
	}
}
