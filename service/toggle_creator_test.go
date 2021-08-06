package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/service"
	mock_service "github.com/indrasaputra/toggle/test/mock/service"
)

var (
	testCtx        = context.Background()
	testToggleKeys = []string{"toggle-1", "toggle-2", "toggle-3"}
)

type ToggleCreatorExecutor struct {
	creator   *service.ToggleCreator
	repo      *mock_service.MockCreateToggleRepository
	publisher *mock_service.MockTogglePublisher
}

func TestNewToggleCreator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of ToggleCreator", func(t *testing.T) {
		exec := createToggleCreatorExecutor(ctrl)
		assert.NotNil(t, exec.creator)
	})
}

func TestToggleCreator_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("empty toggle is prohibited", func(t *testing.T) {
		exec := createToggleCreatorExecutor(ctrl)

		err := exec.creator.Create(testCtx, nil)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyToggle(), err)
	})

	t.Run("toggle's key contains character outside of alphanumeric and dash", func(t *testing.T) {
		exec := createToggleCreatorExecutor(ctrl)

		keys := []string{"____", "!@#_", "~)(*&@##!{}:", "key123!!!"}
		for _, key := range keys {
			toggle := &entity.Toggle{Key: key}
			err := exec.creator.Create(testCtx, toggle)

			assert.NotNil(t, err)
			assert.Equal(t, entity.ErrInvalidKey(), err)
		}
	})

	t.Run("repository returns error", func(t *testing.T) {
		exec := createToggleCreatorExecutor(ctrl)

		for _, key := range testToggleKeys {
			toggle := &entity.Toggle{Key: key}
			exec.repo.EXPECT().Insert(testCtx, toggle).Return(entity.ErrInternal(""))

			err := exec.creator.Create(testCtx, toggle)

			assert.NotNil(t, err)
			assert.Equal(t, entity.ErrInternal(""), err)
		}
	})

	t.Run("toggles already exitst", func(t *testing.T) {
		exec := createToggleCreatorExecutor(ctrl)

		for _, key := range testToggleKeys {
			toggle := &entity.Toggle{Key: key}
			exec.repo.EXPECT().Insert(testCtx, toggle).Return(entity.ErrAlreadyExists())

			err := exec.creator.Create(testCtx, toggle)

			assert.NotNil(t, err)
			assert.Equal(t, entity.ErrAlreadyExists(), err)
		}
	})

	t.Run("successfully save a new toggle, but fail to publish event", func(t *testing.T) {
		exec := createToggleCreatorExecutor(ctrl)

		for _, key := range testToggleKeys {
			toggle := &entity.Toggle{Key: key}
			exec.repo.EXPECT().Insert(testCtx, toggle).Return(nil)
			exec.publisher.EXPECT().Publish(testCtx, gomock.Any()).Return(errors.New("error"))

			err := exec.creator.Create(testCtx, toggle)

			assert.Nil(t, err)
		}
	})

	t.Run("successfully save and publish a new toggle", func(t *testing.T) {
		exec := createToggleCreatorExecutor(ctrl)

		for _, key := range testToggleKeys {
			toggle := &entity.Toggle{Key: key}
			exec.repo.EXPECT().Insert(testCtx, toggle).Return(nil)
			exec.publisher.EXPECT().Publish(testCtx, gomock.Any()).Return(nil)

			err := exec.creator.Create(testCtx, toggle)

			assert.Nil(t, err)
		}
	})
}

func createToggleCreatorExecutor(ctrl *gomock.Controller) *ToggleCreatorExecutor {
	r := mock_service.NewMockCreateToggleRepository(ctrl)
	p := mock_service.NewMockTogglePublisher(ctrl)
	c := service.NewToggleCreator(r, p)
	return &ToggleCreatorExecutor{
		creator:   c,
		repo:      r,
		publisher: p,
	}
}
