package service

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/service"
)

const (
	tagService = "service"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "Create")
	defer span.Finish()

	err := t.creator.Create(ctx, toggle)

	span.SetTag(tagService, "Create")
	return err
}

// DeleteByKey decorates DeleteByKey method.
func (t *Tracing) DeleteByKey(ctx context.Context, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DeleteByKey")
	defer span.Finish()

	err := t.deleter.DeleteByKey(ctx, key)

	span.SetTag(tagService, "DeleteByKey")
	return err
}

// GetByKey decorates GetByKey method.
func (t *Tracing) GetByKey(ctx context.Context, key string) (*entity.Toggle, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetByKey")
	defer span.Finish()

	resp, err := t.getter.GetByKey(ctx, key)

	span.SetTag(tagService, "GetByKey")
	return resp, err
}

// GetAll decorates GetAll method.
func (t *Tracing) GetAll(ctx context.Context) ([]*entity.Toggle, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetAll")
	defer span.Finish()

	resp, err := t.getter.GetAll(ctx)

	span.SetTag(tagService, "GetAll")
	return resp, err
}

// Enable decorates Enable method.
func (t *Tracing) Enable(ctx context.Context, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Enable")
	defer span.Finish()

	err := t.enabler.Enable(ctx, key)

	span.SetTag(tagService, "Enable")
	return err
}

// Disable decorates Disable method.
func (t *Tracing) Disable(ctx context.Context, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Disable")
	defer span.Finish()

	err := t.disabler.Disable(ctx, key)

	span.SetTag(tagService, "Disable")
	return err
}
