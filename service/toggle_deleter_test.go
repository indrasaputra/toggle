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
	testToggleKey      = "toggle-1"
	testToggle         = &entity.Toggle{Key: testToggleKey, IsEnabled: testToggleIsEnabledTrue}
	testToggleDisabled = &entity.Toggle{Key: testToggleKey, IsEnabled: testToggleIsEnabledFalse}
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
		exec.repo.EXPECT().GetByKey(testCtx, testToggleKey).Return(nil, entity.ErrInternal(""))

		err := exec.deleter.DeleteByKey(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
	})

	t.Run("enabled toggle can't be deleted", func(t *testing.T) {
		exec := createToggleDeleterExecutor(ctrl)
		exec.repo.EXPECT().GetByKey(testCtx, testToggleKey).Return(testToggle, nil)

		err := exec.deleter.DeleteByKey(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrProhibitedToDelete(), err)
	})

	t.Run("repository returns error for delete toggle", func(t *testing.T) {
		exec := createToggleDeleterExecutor(ctrl)
		exec.repo.EXPECT().GetByKey(testCtx, testToggleKey).Return(testToggleDisabled, nil)
		exec.repo.EXPECT().DeleteByKey(testCtx, testToggleKey).Return(entity.ErrInternal(""))

		err := exec.deleter.DeleteByKey(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
	})

	t.Run("successfully delete a single toggle", func(t *testing.T) {
		exec := createToggleDeleterExecutor(ctrl)
		exec.repo.EXPECT().GetByKey(testCtx, testToggleKey).Return(testToggleDisabled, nil)
		exec.repo.EXPECT().DeleteByKey(testCtx, testToggleKey).Return(nil)

		err := exec.deleter.DeleteByKey(testCtx, testToggleKey)

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
