package repository_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/internal/repository"
	mock_repository "github.com/indrasaputra/toggle/test/mock/repository"
)

var (
	testToggleIsEnabledTrue = true
)

type ToggleUpdaterExecutor struct {
	updater  *repository.ToggleUpdater
	database *mock_repository.MockUpdateToggleDatabase
	cache    *mock_repository.MockUpdateToggleCache
}

func TestNewToggleUpdater(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of ToggleUpdater", func(t *testing.T) {
		exec := createToggleUpdaterExecutor(ctrl)
		assert.NotNil(t, exec.updater)
	})
}

func TestToggleUpdater_UpdateIsEnabled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("database returns error", func(t *testing.T) {
		exec := createToggleUpdaterExecutor(ctrl)
		exec.database.EXPECT().UpdateIsEnabled(testCtx, testToggleKey, testToggleIsEnabledTrue).Return(entity.ErrInternal(""))

		err := exec.updater.UpdateIsEnabled(testCtx, testToggleKey, testToggleIsEnabledTrue)

		assert.NotNil(t, err)
	})

	t.Run("cache error is ignored", func(t *testing.T) {
		exec := createToggleUpdaterExecutor(ctrl)
		exec.database.EXPECT().UpdateIsEnabled(testCtx, testToggleKey, testToggleIsEnabledTrue).Return(nil)
		exec.cache.EXPECT().SetIsEnabled(testCtx, testToggleKey, testToggleIsEnabledTrue).Return(entity.ErrInternal(""))

		err := exec.updater.UpdateIsEnabled(testCtx, testToggleKey, testToggleIsEnabledTrue)

		assert.Nil(t, err)
	})

	t.Run("all steps are successful", func(t *testing.T) {
		exec := createToggleUpdaterExecutor(ctrl)
		exec.database.EXPECT().UpdateIsEnabled(testCtx, testToggleKey, testToggleIsEnabledTrue).Return(nil)
		exec.cache.EXPECT().SetIsEnabled(testCtx, testToggleKey, testToggleIsEnabledTrue).Return(nil)

		err := exec.updater.UpdateIsEnabled(testCtx, testToggleKey, testToggleIsEnabledTrue)

		assert.Nil(t, err)
	})
}

func createToggleUpdaterExecutor(ctrl *gomock.Controller) *ToggleUpdaterExecutor {
	d := mock_repository.NewMockUpdateToggleDatabase(ctrl)
	c := mock_repository.NewMockUpdateToggleCache(ctrl)
	i := repository.NewToggleUpdater(d, c)
	return &ToggleUpdaterExecutor{
		updater:  i,
		database: d,
		cache:    c,
	}
}
