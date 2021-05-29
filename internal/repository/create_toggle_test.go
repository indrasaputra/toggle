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

var (
	testCtx                = context.Background()
	testToggleKey          = "toggle-1"
	testToggleDesc         = "description"
	testToggle             = &entity.Toggle{Key: testToggleKey, Description: testToggleDesc}
	errPostgresInternalMsg = "database down"
)

type ToggleInserterExecutor struct {
	repo     *repository.ToggleInserter
	database *mock_repository.MockInsertToggleDatabase
	cache    *mock_repository.MockSetToggleCache
}

func TestNewToggleInserter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of ToggleInserter", func(t *testing.T) {
		exec := createToggleInserterExecutor(ctrl)
		assert.NotNil(t, exec.repo)
	})
}

func TestToggleInserter_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("empty toggle is prohibited", func(t *testing.T) {
		exec := createToggleInserterExecutor(ctrl)

		err := exec.repo.Insert(testCtx, nil)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyToggle(), err)
	})

	t.Run("database returns error", func(t *testing.T) {
		exec := createToggleInserterExecutor(ctrl)
		exec.database.EXPECT().Insert(testCtx, testToggle).Return(entity.ErrInternal(errPostgresInternalMsg))

		err := exec.repo.Insert(testCtx, testToggle)

		assert.NotNil(t, err)
	})

	t.Run("cache error is ignored", func(t *testing.T) {
		exec := createToggleInserterExecutor(ctrl)
		exec.database.EXPECT().Insert(testCtx, testToggle).Return(nil)
		exec.cache.EXPECT().Set(testCtx, testToggle).Return(entity.ErrInternal(errPostgresInternalMsg))

		err := exec.repo.Insert(testCtx, testToggle)

		assert.Nil(t, err)
	})

	t.Run("all steps are successful", func(t *testing.T) {
		exec := createToggleInserterExecutor(ctrl)
		exec.database.EXPECT().Insert(testCtx, testToggle).Return(nil)
		exec.cache.EXPECT().Set(testCtx, testToggle).Return(nil)

		err := exec.repo.Insert(testCtx, testToggle)

		assert.Nil(t, err)
	})
}

func createToggleInserterExecutor(ctrl *gomock.Controller) *ToggleInserterExecutor {
	d := mock_repository.NewMockInsertToggleDatabase(ctrl)
	c := mock_repository.NewMockSetToggleCache(ctrl)
	i := repository.NewToggleInserter(d, c)
	return &ToggleInserterExecutor{
		repo:     i,
		database: d,
		cache:    c,
	}
}
