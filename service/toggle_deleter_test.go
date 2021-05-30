package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/service"
	mock_service "github.com/indrasaputra/toggle/test/mock/service"
)

var (
	testKey            = "toggle-1"
	testToggle         = &entity.Toggle{Key: testKey, IsEnabled: true}
	testToggleDisabled = &entity.Toggle{Key: testKey, IsEnabled: false}
)

type ToggleDeleterExecutor struct {
	deleter *service.ToggleDeleter
	repo    *mock_service.MockDeleteToggleRepository
}

func TestNewToggleDeleter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of ToggleDeleter", func(t *testing.T) {
		exec := createToggleDeleterExecutor(ctrl)
		assert.NotNil(t, exec.deleter)
	})
}

func TestToggleDeleter_DeleteByKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("repository returns error for find toggle", func(t *testing.T) {
		exec := createToggleDeleterExecutor(ctrl)
		exec.repo.EXPECT().GetByKey(testCtx, testKey).Return(nil, entity.ErrInternal(""))

		err := exec.deleter.DeleteByKey(testCtx, testKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
	})

	t.Run("enabled toggle can't be deleted", func(t *testing.T) {
		exec := createToggleDeleterExecutor(ctrl)
		exec.repo.EXPECT().GetByKey(testCtx, testKey).Return(testToggle, nil)

		err := exec.deleter.DeleteByKey(testCtx, testKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrProhibitedToDelete(), err)
	})

	t.Run("repository returns error for delete toggle", func(t *testing.T) {
		exec := createToggleDeleterExecutor(ctrl)
		exec.repo.EXPECT().GetByKey(testCtx, testKey).Return(testToggleDisabled, nil)
		exec.repo.EXPECT().DeleteByKey(testCtx, testKey).Return(entity.ErrInternal(""))

		err := exec.deleter.DeleteByKey(testCtx, testKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
	})

	t.Run("successfully delete a single toggle", func(t *testing.T) {
		exec := createToggleDeleterExecutor(ctrl)
		exec.repo.EXPECT().GetByKey(testCtx, testKey).Return(testToggleDisabled, nil)
		exec.repo.EXPECT().DeleteByKey(testCtx, testKey).Return(nil)

		err := exec.deleter.DeleteByKey(testCtx, testKey)

		assert.Nil(t, err)
	})
}

func createToggleDeleterExecutor(ctrl *gomock.Controller) *ToggleDeleterExecutor {
	r := mock_service.NewMockDeleteToggleRepository(ctrl)
	d := service.NewToggleDeleter(r)
	return &ToggleDeleterExecutor{
		deleter: d,
		repo:    r,
	}
}
