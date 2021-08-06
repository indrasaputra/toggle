package messaging_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/internal/messaging"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
	mock_messaging "github.com/indrasaputra/toggle/test/mock/messaging"
)

var (
	testCtx = context.Background()
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
		err := exec.publisher.Publish(testCtx, &togglev1.EventToggle{})
		assert.NotNil(t, err)
	})

	t.Run("fail write message", func(t *testing.T) {
		exec := createKafkaPublisherExecutor(ctrl)
		exec.writer.EXPECT().WriteMessages(testCtx, gomock.Any()).Return(errors.New("error"))

		err := exec.publisher.Publish(testCtx, &togglev1.EventToggle{Toggle: &togglev1.Toggle{}})

		assert.NotNil(t, err)
	})

	t.Run("success write message", func(t *testing.T) {
		exec := createKafkaPublisherExecutor(ctrl)
		exec.writer.EXPECT().WriteMessages(testCtx, gomock.Any()).Return(nil)

		err := exec.publisher.Publish(testCtx, &togglev1.EventToggle{Toggle: &togglev1.Toggle{}})

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
