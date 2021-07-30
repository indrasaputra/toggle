package handler_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/internal/grpc/handler"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
	mock_service "github.com/indrasaputra/toggle/test/mock/service"
)

var (
	testGetToggleByKeyRequest  = &togglev1.GetToggleByKeyRequest{Key: testToggleKey}
	testGetToggleByKeyResponse = &togglev1.GetToggleByKeyResponse{Toggle: testToggleProto}
	testGetAllTogglesRequest   = &togglev1.GetAllTogglesRequest{}
	testGetAllTogglesResponse  = &togglev1.GetAllTogglesResponse{Toggles: []*togglev1.Toggle{testToggleProto}}
)

type ToggleQueryExecutor struct {
	handler *handler.ToggleQuery

	getter *mock_service.MockGetToggle
}

func TestNewToggleQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successful create an instance of ToggleQuery", func(t *testing.T) {
		exec := createToggleQueryExecutor(ctrl)
		assert.NotNil(t, exec.handler)
	})
}

func TestToggleQuery_GetToggleByKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		exec := createToggleQueryExecutor(ctrl)

		res, err := exec.handler.GetToggleByKey(testCtx, nil)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyToggle(), err)
		assert.Nil(t, res)
	})

	t.Run("toggle not found", func(t *testing.T) {
		exec := createToggleQueryExecutor(ctrl)
		exec.getter.EXPECT().GetByKey(testCtx, testToggleKey).Return(nil, entity.ErrNotFound())

		res, err := exec.handler.GetToggleByKey(testCtx, testGetToggleByKeyRequest)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
		assert.Nil(t, res)
	})

	t.Run("getter service returns error", func(t *testing.T) {
		exec := createToggleQueryExecutor(ctrl)
		exec.getter.EXPECT().GetByKey(testCtx, testToggleKey).Return(nil, entity.ErrInternal(""))

		res, err := exec.handler.GetToggleByKey(testCtx, testGetToggleByKeyRequest)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Nil(t, res)
	})

	t.Run("success get a single toggle", func(t *testing.T) {
		exec := createToggleQueryExecutor(ctrl)
		exec.getter.EXPECT().GetByKey(testCtx, testToggleKey).Return(testToggleResult, nil)

		res, err := exec.handler.GetToggleByKey(testCtx, testGetToggleByKeyRequest)

		assert.Nil(t, err)
		assert.Equal(t, testGetToggleByKeyResponse, res)
	})
}

func TestToggleQuery_GetAllToggles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		exec := createToggleQueryExecutor(ctrl)

		res, err := exec.handler.GetAllToggles(testCtx, nil)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyToggle(), err)
		assert.Nil(t, res)
	})

	t.Run("getter service returns error", func(t *testing.T) {
		exec := createToggleQueryExecutor(ctrl)
		exec.getter.EXPECT().GetAll(testCtx).Return([]*entity.Toggle{}, entity.ErrInternal(""))

		res, err := exec.handler.GetAllToggles(testCtx, testGetAllTogglesRequest)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Nil(t, res)
	})

	t.Run("success get all toggles", func(t *testing.T) {
		exec := createToggleQueryExecutor(ctrl)
		exec.getter.EXPECT().GetAll(testCtx).Return([]*entity.Toggle{testToggleResult}, nil)

		res, err := exec.handler.GetAllToggles(testCtx, testGetAllTogglesRequest)

		assert.Nil(t, err)
		assert.Equal(t, testGetAllTogglesResponse, res)
	})
}

func createToggleQueryExecutor(ctrl *gomock.Controller) *ToggleQueryExecutor {
	g := mock_service.NewMockGetToggle(ctrl)

	h := handler.NewToggleQuery(g)
	return &ToggleQueryExecutor{
		handler: h,
		getter:  g,
	}
}
