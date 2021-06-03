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
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
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
	testCreateToggleRequest    = &togglev1.CreateToggleRequest{Key: testToggleKey, Description: testToggleDescription}
	testGetToggleByKeyRequest  = &togglev1.GetToggleByKeyRequest{Key: testToggleKey}
	testGetToggleByKeyResponse = &togglev1.GetToggleByKeyResponse{Toggle: testToggleProto}
	testGetAllTogglesRequest   = &togglev1.GetAllTogglesRequest{}
	testGetAllTogglesResponse  = &togglev1.GetAllTogglesResponse{Toggles: []*togglev1.Toggle{testToggleProto}}
	testEnableRequest          = &togglev1.EnableRequest{Key: testToggleKey}
	testDisableRequest         = &togglev1.DisableRequest{Key: testToggleKey}
	testDeleteToggleRequest    = &togglev1.DeleteToggleRequest{Key: testToggleKey}
)

type ToggleExecutor struct {
	handler *handler.Toggle

	creator *mock_service.MockCreateToggle
	getter  *mock_service.MockGetToggle
	updater *mock_service.MockUpdateToggle
	deleter *mock_service.MockDeleteToggle
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

func TestToggle_GetToggleByKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)

		res, err := exec.handler.GetToggleByKey(testCtx, nil)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyToggle(), err)
		assert.Nil(t, res)
	})

	t.Run("toggle not found", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		exec.getter.EXPECT().GetByKey(testCtx, testToggleKey).Return(nil, entity.ErrNotFound())

		res, err := exec.handler.GetToggleByKey(testCtx, testGetToggleByKeyRequest)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
		assert.Nil(t, res)
	})

	t.Run("getter service returns error", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		exec.getter.EXPECT().GetByKey(testCtx, testToggleKey).Return(nil, entity.ErrInternal(""))

		res, err := exec.handler.GetToggleByKey(testCtx, testGetToggleByKeyRequest)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Nil(t, res)
	})

	t.Run("success get a single toggle", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		exec.getter.EXPECT().GetByKey(testCtx, testToggleKey).Return(testToggleResult, nil)

		res, err := exec.handler.GetToggleByKey(testCtx, testGetToggleByKeyRequest)

		assert.Nil(t, err)
		assert.Equal(t, testGetToggleByKeyResponse, res)
	})
}

func TestToggle_GetAllToggles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)

		res, err := exec.handler.GetAllToggles(testCtx, nil)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyToggle(), err)
		assert.Nil(t, res)
	})

	t.Run("getter service returns error", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		exec.getter.EXPECT().GetAll(testCtx).Return([]*entity.Toggle{}, entity.ErrInternal(""))

		res, err := exec.handler.GetAllToggles(testCtx, testGetAllTogglesRequest)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Nil(t, res)
	})

	t.Run("success get all toggles", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		exec.getter.EXPECT().GetAll(testCtx).Return([]*entity.Toggle{testToggleResult}, nil)

		res, err := exec.handler.GetAllToggles(testCtx, testGetAllTogglesRequest)

		assert.Nil(t, err)
		assert.Equal(t, testGetAllTogglesResponse, res)
	})
}

func TestToggle_Enable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)

		res, err := exec.handler.Enable(testCtx, nil)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyToggle(), err)
		assert.Nil(t, res)
	})

	t.Run("toggle not found", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		exec.updater.EXPECT().Enable(testCtx, testToggleKey).Return(entity.ErrNotFound())

		res, err := exec.handler.Enable(testCtx, testEnableRequest)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
		assert.Nil(t, res)
	})

	t.Run("updater service returns internal error", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		exec.updater.EXPECT().Enable(testCtx, testToggleKey).Return(entity.ErrInternal(""))

		res, err := exec.handler.Enable(testCtx, testEnableRequest)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Nil(t, res)
	})

	t.Run("success enable toggle", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		exec.updater.EXPECT().Enable(testCtx, testToggleKey).Return(nil)

		res, err := exec.handler.Enable(testCtx, testEnableRequest)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
}

func TestToggle_Disable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)

		res, err := exec.handler.Disable(testCtx, nil)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyToggle(), err)
		assert.Nil(t, res)
	})

	t.Run("toggle not found", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		exec.updater.EXPECT().Disable(testCtx, testToggleKey).Return(entity.ErrNotFound())

		res, err := exec.handler.Disable(testCtx, testDisableRequest)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
		assert.Nil(t, res)
	})

	t.Run("updater service returns internal error", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		exec.updater.EXPECT().Disable(testCtx, testToggleKey).Return(entity.ErrInternal(""))

		res, err := exec.handler.Disable(testCtx, testDisableRequest)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Nil(t, res)
	})

	t.Run("success disable toggle", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		exec.updater.EXPECT().Disable(testCtx, testToggleKey).Return(nil)

		res, err := exec.handler.Disable(testCtx, testDisableRequest)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
}

func TestToggle_DeleteToggle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)

		res, err := exec.handler.DeleteToggle(testCtx, nil)

		assert.NotNil(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Nil(t, res)
	})

	t.Run("deleter service returns internal error", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		exec.deleter.EXPECT().DeleteByKey(testCtx, testToggleKey).Return(entity.ErrInternal(""))

		res, err := exec.handler.DeleteToggle(testCtx, testDeleteToggleRequest)

		assert.NotNil(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Nil(t, res)
	})

	t.Run("success delete toggle", func(t *testing.T) {
		exec := createToggleExecutor(ctrl)
		exec.deleter.EXPECT().DeleteByKey(testCtx, testToggleKey).Return(nil)

		res, err := exec.handler.DeleteToggle(testCtx, testDeleteToggleRequest)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
}

func createToggleExecutor(ctrl *gomock.Controller) *ToggleExecutor {
	c := mock_service.NewMockCreateToggle(ctrl)
	g := mock_service.NewMockGetToggle(ctrl)
	u := mock_service.NewMockUpdateToggle(ctrl)
	d := mock_service.NewMockDeleteToggle(ctrl)

	h := handler.NewToggle(c, g, u, d)
	return &ToggleExecutor{
		handler: h,
		creator: c,
		getter:  g,
		updater: u,
		deleter: d,
	}
}
