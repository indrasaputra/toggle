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
	testToggleIsEnabledFalse = false
)

type ToggleDisablerExecutor struct {
	updater *service.ToggleDisabler
	repo    *mock_service.MockDisableToggleRepository
}

func TestNewToggleDisabler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of ToggleDisabler", func(t *testing.T) {
		exec := createToggleDisablerExecutor(ctrl)
		assert.NotNil(t, exec.updater)
	})
}

func TestToggleDisabler_Disable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("repository returns internal error", func(t *testing.T) {
		exec := createToggleDisablerExecutor(ctrl)
		exec.repo.EXPECT().Disable(testCtx, testToggleKey, testToggleIsEnabledFalse).Return(entity.ErrInternal(""))

		err := exec.updater.Disable(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
	})

	t.Run("repository returns not found error", func(t *testing.T) {
		exec := createToggleDisablerExecutor(ctrl)
		exec.repo.EXPECT().Disable(testCtx, testToggleKey, testToggleIsEnabledFalse).Return(entity.ErrNotFound())

		err := exec.updater.Disable(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
	})

	t.Run("successfully disable a toggle", func(t *testing.T) {
		exec := createToggleDisablerExecutor(ctrl)
		exec.repo.EXPECT().Disable(testCtx, testToggleKey, testToggleIsEnabledFalse).Return(nil)

		err := exec.updater.Disable(testCtx, testToggleKey)

		assert.Nil(t, err)
	})
}

func createToggleDisablerExecutor(ctrl *gomock.Controller) *ToggleDisablerExecutor {
	r := mock_service.NewMockDisableToggleRepository(ctrl)
	u := service.NewToggleDisabler(r)
	return &ToggleDisablerExecutor{
		updater: u,
		repo:    r,
	}
}
