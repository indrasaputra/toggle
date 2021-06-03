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
	testToggleIsEnabledTrue  = true
	testToggleIsEnabledFalse = false
)

type ToggleUpdaterExecutor struct {
	updater *service.ToggleUpdater
	repo    *mock_service.MockUpdateToggleRepository
}

func TestNewToggleUpdater(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of ToggleUpdater", func(t *testing.T) {
		exec := createToggleUpdaterExecutor(ctrl)
		assert.NotNil(t, exec.updater)
	})
}

func TestToggleUpdater_Enable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("repository returns internal error", func(t *testing.T) {
		exec := createToggleUpdaterExecutor(ctrl)
		exec.repo.EXPECT().UpdateIsEnabled(testCtx, testToggleKey, testToggleIsEnabledTrue).Return(entity.ErrInternal(""))

		err := exec.updater.Enable(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
	})

	t.Run("repository returns not found error", func(t *testing.T) {
		exec := createToggleUpdaterExecutor(ctrl)
		exec.repo.EXPECT().UpdateIsEnabled(testCtx, testToggleKey, testToggleIsEnabledTrue).Return(entity.ErrNotFound())

		err := exec.updater.Enable(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
	})

	t.Run("successfully enable a toggle", func(t *testing.T) {
		exec := createToggleUpdaterExecutor(ctrl)
		exec.repo.EXPECT().UpdateIsEnabled(testCtx, testToggleKey, testToggleIsEnabledTrue).Return(nil)

		err := exec.updater.Enable(testCtx, testToggleKey)

		assert.Nil(t, err)
	})
}

func TestToggleUpdater_Disable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("repository returns internal error", func(t *testing.T) {
		exec := createToggleUpdaterExecutor(ctrl)
		exec.repo.EXPECT().UpdateIsEnabled(testCtx, testToggleKey, testToggleIsEnabledFalse).Return(entity.ErrInternal(""))

		err := exec.updater.Disable(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
	})

	t.Run("repository returns not found error", func(t *testing.T) {
		exec := createToggleUpdaterExecutor(ctrl)
		exec.repo.EXPECT().UpdateIsEnabled(testCtx, testToggleKey, testToggleIsEnabledFalse).Return(entity.ErrNotFound())

		err := exec.updater.Disable(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
	})

	t.Run("successfully disable a toggle", func(t *testing.T) {
		exec := createToggleUpdaterExecutor(ctrl)
		exec.repo.EXPECT().UpdateIsEnabled(testCtx, testToggleKey, testToggleIsEnabledFalse).Return(nil)

		err := exec.updater.Disable(testCtx, testToggleKey)

		assert.Nil(t, err)
	})
}

func createToggleUpdaterExecutor(ctrl *gomock.Controller) *ToggleUpdaterExecutor {
	r := mock_service.NewMockUpdateToggleRepository(ctrl)
	u := service.NewToggleUpdater(r)
	return &ToggleUpdaterExecutor{
		updater: u,
		repo:    r,
	}
}
