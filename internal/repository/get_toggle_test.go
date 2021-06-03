package repository_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/internal/repository"
	mock_repository "github.com/indrasaputra/toggle/test/mock/repository"
)

type ToggleGetterExecutor struct {
	getter   *repository.ToggleGetter
	database *mock_repository.MockGetToggleDatabase
	cache    *mock_repository.MockGetToggleCache
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

	t.Run("cache returns error", func(t *testing.T) {
		exec := createToggleGetterExecutor(ctrl)
		exec.cache.EXPECT().Get(testCtx, testToggle.Key).Return(nil, entity.ErrInternal(""))

		res, err := exec.getter.GetByKey(testCtx, testToggle.Key)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Nil(t, res)
	})

	t.Run("toggle found in cache", func(t *testing.T) {
		exec := createToggleGetterExecutor(ctrl)
		exec.cache.EXPECT().Get(testCtx, testToggle.Key).Return(testToggle, nil)

		res, err := exec.getter.GetByKey(testCtx, testToggle.Key)

		assert.Nil(t, err)
		assert.Equal(t, testToggle, res)
	})

	t.Run("database returns error", func(t *testing.T) {
		exec := createToggleGetterExecutor(ctrl)
		exec.cache.EXPECT().Get(testCtx, testToggle.Key).Return(nil, nil)
		exec.database.EXPECT().GetByKey(testCtx, testToggle.Key).Return(nil, entity.ErrInternal(""))

		res, err := exec.getter.GetByKey(testCtx, testToggle.Key)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Nil(t, res)
	})

	t.Run("success get toggle from db and save to cache", func(t *testing.T) {
		exec := createToggleGetterExecutor(ctrl)
		exec.cache.EXPECT().Get(testCtx, testToggle.Key).Return(nil, nil)
		exec.database.EXPECT().GetByKey(testCtx, testToggle.Key).Return(testToggle, nil)
		exec.cache.EXPECT().Set(testCtx, testToggle).Return(nil)

		res, err := exec.getter.GetByKey(testCtx, testToggle.Key)

		assert.Nil(t, err)
		assert.Equal(t, testToggle, res)
	})
}

func TestToggleGetter_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("database returns error", func(t *testing.T) {
		exec := createToggleGetterExecutor(ctrl)
		exec.database.EXPECT().GetAll(testCtx, repository.DefaultToggleLimit).Return([]*entity.Toggle{}, entity.ErrInternal(""))

		res, err := exec.getter.GetAll(testCtx)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Empty(t, res)
	})

	t.Run("database returns empty list and nil error", func(t *testing.T) {
		exec := createToggleGetterExecutor(ctrl)
		exec.database.EXPECT().GetAll(testCtx, repository.DefaultToggleLimit).Return([]*entity.Toggle{}, nil)

		res, err := exec.getter.GetAll(testCtx)

		assert.Nil(t, err)
		assert.Empty(t, res)
	})

	t.Run("success get toggle from db", func(t *testing.T) {
		exec := createToggleGetterExecutor(ctrl)
		exec.database.EXPECT().GetAll(testCtx, repository.DefaultToggleLimit).Return([]*entity.Toggle{testToggle}, nil)

		res, err := exec.getter.GetAll(testCtx)

		assert.Nil(t, err)
		assert.NotEmpty(t, testToggle, res)
	})
}

func createToggleGetterExecutor(ctrl *gomock.Controller) *ToggleGetterExecutor {
	d := mock_repository.NewMockGetToggleDatabase(ctrl)
	c := mock_repository.NewMockGetToggleCache(ctrl)
	i := repository.NewToggleGetter(d, c)
	return &ToggleGetterExecutor{
		getter:   i,
		database: d,
		cache:    c,
	}
}
