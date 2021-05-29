package handler_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/internal/grpc/handler"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
	mock_service "github.com/indrasaputra/toggle/test/mock/service"
)

var (
	testCtx                 = context.Background()
	testToggleKey           = "test_key"
	testToggleIsEnabled     = false
	testToggle              = &entity.Toggle{Key: testToggleKey, IsEnabled: testToggleIsEnabled}
	testCreateToggleRequest = &togglev1.CreateToggleRequest{Key: testToggleKey}
)

type ToggleExecutor struct {
	handler *handler.Toggle
	creator *mock_service.MockCreateToggle
}

func TestNewToggle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successful create an instance of Toggle", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		assert.NotNil(t, exec.handler)
	})
}

func TestToggle_CreateToggle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)

		res, err := exec.handler.CreateToggle(testCtx, nil)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyToggle(), err)
		assert.Nil(t, res)
	})

	t.Run("request attributes are invalid", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		errTables := []error{entity.ErrEmptyToggle(), entity.ErrInvalidKey()}

		for _, errTab := range errTables {
			exec.creator.EXPECT().Create(testCtx, testToggle).Return(errTab)

			res, err := exec.handler.CreateToggle(testCtx, testCreateToggleRequest)

			assert.NotNil(t, err)
			assert.Equal(t, errTab, err)
			assert.Nil(t, res)
		}
	})

	t.Run("creator service returns internal error", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		exec.creator.EXPECT().Create(testCtx, testToggle).Return(entity.ErrInternal(""))

		res, err := exec.handler.CreateToggle(testCtx, testCreateToggleRequest)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Nil(t, res)
	})

	t.Run("success create a toggle", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		exec.creator.EXPECT().Create(testCtx, testToggle).Return(nil)

		res, err := exec.handler.CreateToggle(testCtx, testCreateToggleRequest)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
}

func createToggleExecutor(ctrl *gomock.Controller) *ToggleExecutor {
	c := mock_service.NewMockCreateToggle(ctrl)

	h := handler.NewToggle(c)
	return &ToggleExecutor{
		handler: h,
		creator: c,
	}
}
