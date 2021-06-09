package redis

import (
	"context"
	"fmt"

	goredis "github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	dbType = "redis"
)

// HookTracing traces redis command.
type HookTracing struct {
}

// NewHookTracing creates an instance of HookTracing.
func NewHookTracing() *HookTracing {
	return &HookTracing{}
}

// BeforeProcess runs the tracing before the command executed.
func (ht *HookTracing) BeforeProcess(ctx context.Context, cmd goredis.Cmder) (context.Context, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, cmd.Name())
	defer span.Finish()

	ext.DBType.Set(span, dbType)
	ext.DBStatement.Set(span, fmt.Sprintf("%v", cmd.Args()))

	return ctx, nil
}

// AfterProcess runs the tracing after the command executed.
func (ht *HookTracing) AfterProcess(ctx context.Context, cmd goredis.Cmder) error {
	return nil
}

// BeforeProcessPipeline runs the tracing before the pipeline executed.
func (ht *HookTracing) BeforeProcessPipeline(ctx context.Context, cmds []goredis.Cmder) (context.Context, error) {
	pipelineSpan, ctx := opentracing.StartSpanFromContext(ctx, "redis-pipeline")
	defer pipelineSpan.Finish()

	ext.DBType.Set(pipelineSpan, dbType)

	for _, cmd := range cmds {
		span, _ := opentracing.StartSpanFromContext(ctx, cmd.Name())
		defer span.Finish()

		ext.DBType.Set(pipelineSpan, dbType)
		ext.DBStatement.Set(span, fmt.Sprintf("%v", cmd.Args()))
	}

	return ctx, nil
}

// AfterProcessPipelineruns the tracing after the pipeline executed.
func (ht *HookTracing) AfterProcessPipeline(ctx context.Context, cmds []goredis.Cmder) error {
	return nil
}
