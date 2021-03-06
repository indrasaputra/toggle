package service_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/service"
	mock_service "github.com/indrasaputra/toggle/test/mock/service"
)

var (
	testToggleIsEnabledTrue = true
)

type ToggleEnablerExecutor struct {
	enabler   *service.ToggleEnabler
	repo      *mock_service.MockEnableToggleRepository
	publisher *mock_service.MockTogglePublisher
}

func TestNewToggleEnabler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of ToggleEnabler", func(t *testing.T) {
		exec := createToggleEnablerExecutor(ctrl)
		assert.NotNil(t, exec.enabler)
	})
}

func TestToggleEnabler_Enable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("repository returns internal error", func(t *testing.T) {
		exec := createToggleEnablerExecutor(ctrl)
		exec.repo.EXPECT().Enable(testCtx, testToggleKey, testToggleIsEnabledTrue).Return(entity.ErrInternal(""))

		err := exec.enabler.Enable(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
	})

	t.Run("repository returns not found error", func(t *testing.T) {
		exec := createToggleEnablerExecutor(ctrl)
		exec.repo.EXPECT().Enable(testCtx, testToggleKey, testToggleIsEnabledTrue).Return(entity.ErrNotFound())

		err := exec.enabler.Enable(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
	})

	t.Run("successfully enable a toggle, but fail to publish", func(t *testing.T) {
		exec := createToggleEnablerExecutor(ctrl)
		exec.repo.EXPECT().Enable(testCtx, testToggleKey, testToggleIsEnabledTrue).Return(nil)
		exec.publisher.EXPECT().Publish(testCtx, gomock.Any()).Return(errors.New("error"))

		err := exec.enabler.Enable(testCtx, testToggleKey)

		assert.Nil(t, err)
	})

	t.Run("successfully enable and publish a toggle", func(t *testing.T) {
		exec := createToggleEnablerExecutor(ctrl)
		exec.repo.EXPECT().Enable(testCtx, testToggleKey, testToggleIsEnabledTrue).Return(nil)
		exec.publisher.EXPECT().Publish(testCtx, gomock.Any()).Return(nil)

		err := exec.enabler.Enable(testCtx, testToggleKey)

		assert.Nil(t, err)
	})
}

func createToggleEnablerExecutor(ctrl *gomock.Controller) *ToggleEnablerExecutor {
	r := mock_service.NewMockEnableToggleRepository(ctrl)
	p := mock_service.NewMockTogglePublisher(ctrl)
	u := service.NewToggleEnabler(r, p)
	return &ToggleEnablerExecutor{
		enabler:   u,
		repo:      r,
		publisher: p,
	}
}
