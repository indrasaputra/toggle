package interceptor

import (
	"context"
	"path"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	tagMethod = "grpc.method"
	tagCode   = "grpc.code"
)

// OpenTracingUnaryServerInterceptor intercepts the request and creates a span from the incoming context.
// It names the span using the method that is being called.
func OpenTracingUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		method := path.Base(info.FullMethod)
		span, ctx := opentracing.StartSpanFromContext(ctx, method)
		defer span.Finish()

		resp, err := handler(ctx, req)

		span.SetTag(tagMethod, method)
		span.SetTag(tagCode, status.Code(err))

		return resp, err
	}
}
