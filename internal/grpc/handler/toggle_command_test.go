package handler_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/internal/grpc/handler"
	togglev1 "github.com/indrasaputra/toggle/proto/toggle/v1"
	mock_service "github.com/indrasaputra/toggle/test/mock/service"
)

var (
	testCtx               = context.Background()
	testToggleKey         = "test_key"
	testToggleIsEnabled   = false
	testToggleDescription = "description"
	testToggleCreatedAt   = time.Now()
	testToggleUpdatedAt   = time.Now()
	testToggle            = &entity.Toggle{
		Key:         testToggleKey,
		IsEnabled:   testToggleIsEnabled,
		Description: testToggleDescription,
	}
	testToggleResult = &entity.Toggle{
		Key:         testToggleKey,
		IsEnabled:   testToggleIsEnabled,
		Description: testToggleDescription,
		CreatedAt:   testToggleCreatedAt,
		UpdatedAt:   testToggleUpdatedAt,
	}
	testToggleProto = &togglev1.Toggle{
		Key:         testToggleKey,
		IsEnabled:   testToggleIsEnabled,
		Description: testToggleDescription,
		CreatedAt:   timestamppb.New(testToggleCreatedAt),
		UpdatedAt:   timestamppb.New(testToggleUpdatedAt),
	}
	testCreateToggleRequest  = &togglev1.CreateToggleRequest{Toggle: testToggleProto}
	testEnableToggleRequest  = &togglev1.EnableToggleRequest{Key: testToggleKey}
	testDisableToggleRequest = &togglev1.DisableToggleRequest{Key: testToggleKey}
	testDeleteToggleRequest  = &togglev1.DeleteToggleRequest{Key: testToggleKey}
)

type ToggleCommandExecutor struct {
	handler *handler.ToggleCommand

	creator  *mock_service.MockCreateToggle
	enabler  *mock_service.MockEnableToggle
	disabler *mock_service.MockDisableToggle
	deleter  *mock_service.MockDeleteToggle
}

func TestNewToggleCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successful create an instance of ToggleCommand", func(t *testing.T) {
		exec := createToggleCommandExecutor(ctrl)
		assert.NotNil(t, exec.handler)
	})
}

func TestToggleCommand_CreateToggle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		exec := createToggleCommandExecutor(ctrl)

		res, err := exec.handler.CreateToggle(testCtx, nil)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyToggle(), err)
		assert.Nil(t, res)
	})

	t.Run("empty toggle is prohibited", func(t *testing.T) {
		exec := createToggleCommandExecutor(ctrl)

		res, err := exec.handler.CreateToggle(testCtx, &togglev1.CreateToggleRequest{})

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyToggle(), err)
		assert.Nil(t, res)
	})

	t.Run("request attributes are invalid", func(t *testing.T) {
		exec := createToggleCommandExecutor(ctrl)
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
		exec := createToggleCommandExecutor(ctrl)
		exec.creator.EXPECT().Create(testCtx, testToggle).Return(entity.ErrInternal(""))

		res, err := exec.handler.CreateToggle(testCtx, testCreateToggleRequest)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Nil(t, res)
	})

	t.Run("success create a toggle", func(t *testing.T) {
		exec := createToggleCommandExecutor(ctrl)
		exec.creator.EXPECT().Create(testCtx, testToggle).Return(nil)

		res, err := exec.handler.CreateToggle(testCtx, testCreateToggleRequest)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
}

func TestToggleCommand_EnableToggle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		exec := createToggleCommandExecutor(ctrl)

		res, err := exec.handler.EnableToggle(testCtx, nil)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyToggle(), err)
		assert.Nil(t, res)
	})

	t.Run("toggle not found", func(t *testing.T) {
		exec := createToggleCommandExecutor(ctrl)
		exec.enabler.EXPECT().Enable(testCtx, testToggleKey).Return(entity.ErrNotFound())

		res, err := exec.handler.EnableToggle(testCtx, testEnableToggleRequest)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
		assert.Nil(t, res)
	})

	t.Run("updater service returns internal error", func(t *testing.T) {
		exec := createToggleCommandExecutor(ctrl)
		exec.enabler.EXPECT().Enable(testCtx, testToggleKey).Return(entity.ErrInternal(""))

		res, err := exec.handler.EnableToggle(testCtx, testEnableToggleRequest)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Nil(t, res)
	})

	t.Run("success enable toggle", func(t *testing.T) {
		exec := createToggleCommandExecutor(ctrl)
		exec.enabler.EXPECT().Enable(testCtx, testToggleKey).Return(nil)

		res, err := exec.handler.EnableToggle(testCtx, testEnableToggleRequest)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
}

func TestToggleCommand_DisableToggle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		exec := createToggleCommandExecutor(ctrl)

		res, err := exec.handler.DisableToggle(testCtx, nil)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyToggle(), err)
		assert.Nil(t, res)
	})

	t.Run("toggle not found", func(t *testing.T) {
		exec := createToggleCommandExecutor(ctrl)
		exec.disabler.EXPECT().Disable(testCtx, testToggleKey).Return(entity.ErrNotFound())

		res, err := exec.handler.DisableToggle(testCtx, testDisableToggleRequest)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
		assert.Nil(t, res)
	})

	t.Run("updater service returns internal error", func(t *testing.T) {
		exec := createToggleCommandExecutor(ctrl)
		exec.disabler.EXPECT().Disable(testCtx, testToggleKey).Return(entity.ErrInternal(""))

		res, err := exec.handler.DisableToggle(testCtx, testDisableToggleRequest)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Nil(t, res)
	})

	t.Run("success disable toggle", func(t *testing.T) {
		exec := createToggleCommandExecutor(ctrl)
		exec.disabler.EXPECT().Disable(testCtx, testToggleKey).Return(nil)

		res, err := exec.handler.DisableToggle(testCtx, testDisableToggleRequest)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
}

func TestToggleCommand_DeleteToggle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		exec := createToggleCommandExecutor(ctrl)

		res, err := exec.handler.DeleteToggle(testCtx, nil)

		assert.NotNil(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Nil(t, res)
	})

	t.Run("deleter service returns internal error", func(t *testing.T) {
		exec := createToggleCommandExecutor(ctrl)
		exec.deleter.EXPECT().DeleteByKey(testCtx, testToggleKey).Return(entity.ErrInternal(""))

		res, err := exec.handler.DeleteToggle(testCtx, testDeleteToggleRequest)

		assert.NotNil(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Nil(t, res)
	})

	t.Run("success delete toggle", func(t *testing.T) {
		exec := createToggleCommandExecutor(ctrl)
		exec.deleter.EXPECT().DeleteByKey(testCtx, testToggleKey).Return(nil)

		res, err := exec.handler.DeleteToggle(testCtx, testDeleteToggleRequest)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
}

func createToggleCommandExecutor(ctrl *gomock.Controller) *ToggleCommandExecutor {
	c := mock_service.NewMockCreateToggle(ctrl)
	e := mock_service.NewMockEnableToggle(ctrl)
	s := mock_service.NewMockDisableToggle(ctrl)
	d := mock_service.NewMockDeleteToggle(ctrl)

	h := handler.NewToggleCommand(c, e, s, d)
	return &ToggleCommandExecutor{
		handler:  h,
		creator:  c,
		enabler:  e,
		disabler: s,
		deleter:  d,
	}
}
