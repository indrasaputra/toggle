package interceptor_test

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/indrasaputra/toggle/internal/grpc/interceptor"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
	mock_server "github.com/indrasaputra/toggle/test/mock/proto/grpc/server"
)

const (
	buffer = 1024 * 1024
)

var (
	testCtx = context.Background()
)

type ToggleClientExecutor struct {
	client togglev1.ToggleCommandServiceClient
	closer func()
}

func TestOpenTracingUnaryServerInterceptor(t *testing.T) {
	t.Run("success create a new span and finish", func(t *testing.T) {
		exec := createClientExecutor(interceptor.OpenTracingUnaryServerInterceptor())
		defer exec.closer()

		resp, err := exec.client.CreateToggle(testCtx, &togglev1.CreateToggleRequest{})

		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}

func createClientExecutor(intercept grpc.UnaryServerInterceptor) *ToggleClientExecutor {
	listener := bufconn.Listen(buffer)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(intercept))
	toggleServer := &mock_server.MockToggleServiceServer{}
	togglev1.RegisterToggleCommandServiceServer(grpcServer, toggleServer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			panic(err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithInsecure())
	if err != nil {
		grpcServer.Stop()
		panic(err)
	}

	closer := func() {
		_ = listener.Close()
		grpcServer.GracefulStop()
	}

	return &ToggleClientExecutor{
		client: togglev1.NewToggleCommandServiceClient(conn),
		closer: closer,
	}
}
