package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/internal/app"
	"github.com/indrasaputra/toggle/internal/decorator/service"
	mock_service "github.com/indrasaputra/toggle/test/mock/service"
)

var (
	testCtx               = context.Background()
	testToggleKey         = "test_key"
	testToggleIsEnabled   = false
	testToggleDescription = "description"
	testToggle            = &entity.Toggle{
		Key:         testToggleKey,
		IsEnabled:   testToggleIsEnabled,
		Description: testToggleDescription,
	}
)

type TracingExecutor struct {
	tracing *service.Tracing

	creator  *mock_service.MockCreateToggle
	getter   *mock_service.MockGetToggle
	enabler  *mock_service.MockEnableToggle
	disabler *mock_service.MockDisableToggle
	deleter  *mock_service.MockDeleteToggle
}

func TestTracing_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success decorate Create method", func(t *testing.T) {
		ctx, span := app.GetTracer().Start(testCtx, "Create")
		defer span.End()

		exec := createTracingExecutor(ctrl)
		exec.creator.EXPECT().Create(ctx, testToggle).Return(nil)

		err := exec.tracing.Create(testCtx, testToggle)

		assert.Nil(t, err)
	})
}

func TestTracing_DeleteByKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success decorate DeleteByKey method", func(t *testing.T) {
		ctx, span := app.GetTracer().Start(testCtx, "DeleteByKey")
		defer span.End()

		exec := createTracingExecutor(ctrl)
		exec.deleter.EXPECT().DeleteByKey(ctx, testToggleKey).Return(nil)

		err := exec.tracing.DeleteByKey(testCtx, testToggleKey)

		assert.Nil(t, err)
	})
}

func TestTracing_Enable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success decorate Enable method", func(t *testing.T) {
		ctx, span := app.GetTracer().Start(testCtx, "Enable")
		defer span.End()

		exec := createTracingExecutor(ctrl)
		exec.enabler.EXPECT().Enable(ctx, testToggleKey).Return(nil)

		err := exec.tracing.Enable(testCtx, testToggleKey)

		assert.Nil(t, err)
	})
}

func TestTracing_Disable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success decorate Disable method", func(t *testing.T) {
		ctx, span := app.GetTracer().Start(testCtx, "Disable")
		defer span.End()

		exec := createTracingExecutor(ctrl)
		exec.disabler.EXPECT().Disable(ctx, testToggleKey).Return(nil)

		err := exec.tracing.Disable(testCtx, testToggleKey)

		assert.Nil(t, err)
	})
}

func TestTracing_GetByKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success decorate GetByKey method", func(t *testing.T) {
		ctx, span := app.GetTracer().Start(testCtx, "GetByKey")
		defer span.End()

		exec := createTracingExecutor(ctrl)
		exec.getter.EXPECT().GetByKey(ctx, testToggleKey).Return(nil, nil)

		resp, err := exec.tracing.GetByKey(testCtx, testToggleKey)

		assert.Nil(t, err)
		assert.Nil(t, resp)
	})
}

func TestTracing_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success decorate GetAll method", func(t *testing.T) {
		ctx, span := app.GetTracer().Start(testCtx, "GetAll")
		defer span.End()

		exec := createTracingExecutor(ctrl)
		exec.getter.EXPECT().GetAll(ctx).Return(nil, nil)

		resp, err := exec.tracing.GetAll(testCtx)

		assert.Nil(t, err)
		assert.Nil(t, resp)
	})
}

func createTracingExecutor(ctrl *gomock.Controller) *TracingExecutor {
	c := mock_service.NewMockCreateToggle(ctrl)
	g := mock_service.NewMockGetToggle(ctrl)
	e := mock_service.NewMockEnableToggle(ctrl)
	s := mock_service.NewMockDisableToggle(ctrl)
	d := mock_service.NewMockDeleteToggle(ctrl)

	t := service.NewTracing(c, g, e, s, d)
	return &TracingExecutor{
		tracing:  t,
		creator:  c,
		getter:   g,
		enabler:  e,
		disabler: s,
		deleter:  d,
	}
}
