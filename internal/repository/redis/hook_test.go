package redis_test

import (
	"testing"

	goredis "github.com/go-redis/redis/v8"
	"github.com/indrasaputra/toggle/internal/repository/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

func TestNewHookTracing(t *testing.T) {
	t.Run("success create an instance of HookTracing", func(t *testing.T) {
		hook := redis.NewHookTracing()
		assert.NotNil(t, hook)
	})
}

func TestHookTracing_BeforeProcess(t *testing.T) {
	t.Run("success start span before process", func(t *testing.T) {
		cmd := &goredis.StringCmd{}
		span, _ := opentracing.StartSpanFromContext(testCtx, cmd.Name())
		defer span.Finish()
		hook := redis.NewHookTracing()

		ctx, err := hook.BeforeProcess(testCtx, cmd)

		assert.Nil(t, err)
		assert.Equal(t, opentracing.ContextWithSpan(testCtx, span), ctx)
	})
}

func TestHookTracing_AfterProcess(t *testing.T) {
	t.Run("after process does nothing", func(t *testing.T) {
		cmd := &goredis.StringCmd{}
		hook := redis.NewHookTracing()

		err := hook.AfterProcess(testCtx, cmd)

		assert.Nil(t, err)
	})
}

func TestHookTracing_BeforeProcessPipeline(t *testing.T) {
	t.Run("success start span before process pipeline", func(t *testing.T) {
		cmd := &goredis.StringCmd{}
		span, _ := opentracing.StartSpanFromContext(testCtx, cmd.Name())
		defer span.Finish()
		hook := redis.NewHookTracing()

		ctx, err := hook.BeforeProcessPipeline(testCtx, []goredis.Cmder{cmd})

		assert.Nil(t, err)
		assert.Equal(t, opentracing.ContextWithSpan(testCtx, span), ctx)
	})
}

func TestHookTracing_AfterProcessPipeline(t *testing.T) {
	t.Run("after process does nothing", func(t *testing.T) {
		cmd := &goredis.StringCmd{}
		hook := redis.NewHookTracing()

		err := hook.AfterProcessPipeline(testCtx, []goredis.Cmder{cmd})

		assert.Nil(t, err)
	})
}
