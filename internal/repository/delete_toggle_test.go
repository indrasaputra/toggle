package repository_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/internal/repository"
	mock_repository "github.com/indrasaputra/toggle/test/mock/repository"
)

type ToggleDeleterExecutor struct {
	deleter  *repository.ToggleDeleter
	database *mock_repository.MockDeleteToggleDatabase
	cache    *mock_repository.MockDeleteToggleCache
}

func TestNewToggleDeleter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of ToggleDeleter", func(t *testing.T) {
		exec := createToggleDeleterExecutor(ctrl)
		assert.NotNil(t, exec.deleter)
	})
}

func TestToggleDeleter_GetByKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("database returns error", func(t *testing.T) {
		exec := createToggleDeleterExecutor(ctrl)
		exec.database.EXPECT().GetByKey(context.Background(), testToggle.Key).Return(nil, entity.ErrInternal(""))

		res, err := exec.deleter.GetByKey(context.Background(), testToggle.Key)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Nil(t, res)
	})

	t.Run("toggle not found", func(t *testing.T) {
		exec := createToggleDeleterExecutor(ctrl)
		exec.database.EXPECT().GetByKey(context.Background(), testToggle.Key).Return(nil, entity.ErrNotFound())

		res, err := exec.deleter.GetByKey(context.Background(), testToggle.Key)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
		assert.Nil(t, res)
	})

	t.Run("success get toggle from db", func(t *testing.T) {
		exec := createToggleDeleterExecutor(ctrl)
		exec.database.EXPECT().GetByKey(context.Background(), testToggle.Key).Return(testToggle, nil)

		res, err := exec.deleter.GetByKey(context.Background(), testToggle.Key)

		assert.Nil(t, err)
		assert.Equal(t, testToggle, res)
	})
}

func TestToggleDeleter_DeleteByKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("cache returns error", func(t *testing.T) {
		exec := createToggleDeleterExecutor(ctrl)
		exec.cache.EXPECT().Delete(context.Background(), testToggle.Key).Return(entity.ErrInternal(""))

		err := exec.deleter.DeleteByKey(context.Background(), testToggle.Key)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
	})

	t.Run("database returns error", func(t *testing.T) {
		exec := createToggleDeleterExecutor(ctrl)
		exec.cache.EXPECT().Delete(context.Background(), testToggle.Key).Return(nil)
		exec.database.EXPECT().Delete(context.Background(), testToggle.Key).Return(entity.ErrInternal(""))

		err := exec.deleter.DeleteByKey(context.Background(), testToggle.Key)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
	})

	t.Run("success delete toggle from cache and db ", func(t *testing.T) {
		exec := createToggleDeleterExecutor(ctrl)
		exec.cache.EXPECT().Delete(context.Background(), testToggle.Key).Return(nil)
		exec.database.EXPECT().Delete(context.Background(), testToggle.Key).Return(nil)

		err := exec.deleter.DeleteByKey(context.Background(), testToggle.Key)

		assert.Nil(t, err)
	})
}

func createToggleDeleterExecutor(ctrl *gomock.Controller) *ToggleDeleterExecutor {
	d := mock_repository.NewMockDeleteToggleDatabase(ctrl)
	c := mock_repository.NewMockDeleteToggleCache(ctrl)
	r := repository.NewToggleDeleter(d, c)
	return &ToggleDeleterExecutor{
		deleter:  r,
		database: d,
		cache:    c,
	}
}
