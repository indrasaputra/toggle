package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/service"
	mock_service "github.com/indrasaputra/toggle/test/mock/service"
)

type ToggleGetterExecutor struct {
	getter *service.ToggleGetter
	repo   *mock_service.MockGetToggleRepository
}

func TestNewToggleGetter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of ToggleGetter", func(t *testing.T) {
		exec := createToggleGetterExecutor(ctrl)
		assert.NotNil(t, exec.getter)
	})
}

func TestToggleGetter_GetByKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("repository returns internal error", func(t *testing.T) {
		exec := createToggleGetterExecutor(ctrl)
		exec.repo.EXPECT().GetByKey(testCtx, testKey).Return(nil, entity.ErrInternal(""))

		res, err := exec.getter.GetByKey(testCtx, testKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Nil(t, res)
	})

	t.Run("repository returns not found error", func(t *testing.T) {
		exec := createToggleGetterExecutor(ctrl)
		exec.repo.EXPECT().GetByKey(testCtx, testKey).Return(nil, entity.ErrNotFound())

		res, err := exec.getter.GetByKey(testCtx, testKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
		assert.Nil(t, res)
	})

	t.Run("successfully get a single toggle", func(t *testing.T) {
		exec := createToggleGetterExecutor(ctrl)
		exec.repo.EXPECT().GetByKey(testCtx, testKey).Return(testToggle, nil)

		res, err := exec.getter.GetByKey(testCtx, testKey)

		assert.Nil(t, err)
		assert.Equal(t, testToggle, res)
	})
}

func TestToggleGetter_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("repository returns internal error", func(t *testing.T) {
		exec := createToggleGetterExecutor(ctrl)
		exec.repo.EXPECT().GetAll(testCtx).Return([]*entity.Toggle{}, entity.ErrInternal(""))

		res, err := exec.getter.GetAll(testCtx)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Empty(t, res)
	})

	t.Run("repository returns empty list", func(t *testing.T) {
		exec := createToggleGetterExecutor(ctrl)
		exec.repo.EXPECT().GetAll(testCtx).Return([]*entity.Toggle{}, nil)

		res, err := exec.getter.GetAll(testCtx)

		assert.Nil(t, err)
		assert.Empty(t, res)
	})

	t.Run("successfully get a single toggle", func(t *testing.T) {
		exec := createToggleGetterExecutor(ctrl)
		exec.repo.EXPECT().GetAll(testCtx).Return([]*entity.Toggle{testToggle}, nil)

		res, err := exec.getter.GetAll(testCtx)

		assert.Nil(t, err)
		assert.NotEmpty(t, res)
	})
}

func createToggleGetterExecutor(ctrl *gomock.Controller) *ToggleGetterExecutor {
	r := mock_service.NewMockGetToggleRepository(ctrl)
	c := service.NewToggleGetter(r)
	return &ToggleGetterExecutor{
		getter: c,
		repo:   r,
	}
}
