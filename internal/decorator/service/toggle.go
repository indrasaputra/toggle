package service

import (
	"context"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/internal/app"
	"github.com/indrasaputra/toggle/service"
)

// Tracing decorates toggle service and imbues it with tracing.
type Tracing struct {
	creator  service.CreateToggle
	getter   service.GetToggle
	enabler  service.EnableToggle
	disabler service.DisableToggle
	deleter  service.DeleteToggle
}

// NewTracing creates an instance of Tracing.
func NewTracing(creator service.CreateToggle, getter service.GetToggle, enabler service.EnableToggle, disabler service.DisableToggle, deleter service.DeleteToggle) *Tracing {
	return &Tracing{
		creator:  creator,
		getter:   getter,
		enabler:  enabler,
		disabler: disabler,
		deleter:  deleter,
	}
}

// Create decorates Create method.
func (t *Tracing) Create(ctx context.Context, toggle *entity.Toggle) error {
	ctx, span := app.GetTracer().Start(ctx, "Create")
	defer span.End()

	return t.creator.Create(ctx, toggle)
}

// DeleteByKey decorates DeleteByKey method.
func (t *Tracing) DeleteByKey(ctx context.Context, key string) error {
	ctx, span := app.GetTracer().Start(ctx, "DeleteByKey")
	defer span.End()

	return t.deleter.DeleteByKey(ctx, key)
}

// GetByKey decorates GetByKey method.
func (t *Tracing) GetByKey(ctx context.Context, key string) (*entity.Toggle, error) {
	ctx, span := app.GetTracer().Start(ctx, "GetByKey")
	defer span.End()

	resp, err := t.getter.GetByKey(ctx, key)

	return resp, err
}

// GetAll decorates GetAll method.
func (t *Tracing) GetAll(ctx context.Context) ([]*entity.Toggle, error) {
	ctx, span := app.GetTracer().Start(ctx, "GetAll")
	defer span.End()

	resp, err := t.getter.GetAll(ctx)

	return resp, err
}

// Enable decorates Enable method.
func (t *Tracing) Enable(ctx context.Context, key string) error {
	ctx, span := app.GetTracer().Start(ctx, "Enable")
	defer span.End()

	return t.enabler.Enable(ctx, key)
}

// Disable decorates Disable method.
func (t *Tracing) Disable(ctx context.Context, key string) error {
	ctx, span := app.GetTracer().Start(ctx, "Disable")
	defer span.End()

	return t.disabler.Disable(ctx, key)
}
