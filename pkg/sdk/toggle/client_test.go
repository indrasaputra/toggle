package toggle_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/pkg/sdk/toggle"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
	mock_server "github.com/indrasaputra/toggle/test/mock/proto/grpc/server"
)

const (
	buffer = 1024 * 1024
)

var (
	executor       *ClientExecutor
	testCtx        = context.Background()
	testCtxError   = metadata.NewOutgoingContext(testCtx, metadata.Pairs("has-error", "true"))
	testCtxReturn  = metadata.NewOutgoingContext(testCtx, metadata.Pairs("complete-return", "true"))
	testToggleKey  = "toggle-1"
	testToggleDesc = "description-1"
	testToggle     = &entity.Toggle{Key: testToggleKey, Description: testToggleDesc}
)

type ClientExecutor struct {
	client *toggle.Client
	closer func()
}

func TestMain(m *testing.M) {
	executor = createClientExecutor()

	code := m.Run()

	executor.closer()
	os.Exit(code)
}

func TestNewClient(t *testing.T) {
	t.Run("successfully create a new Client", func(t *testing.T) {
		assert.NotNil(t, executor.client)
	})
}

func TestClient_Create(t *testing.T) {
	t.Run("server returns error", func(t *testing.T) {
		err := executor.client.Create(testCtxError, testToggle)

		assert.NotNil(t, err)
	})

	t.Run("success create a new toggle", func(t *testing.T) {
		err := executor.client.Create(testCtx, testToggle)

		assert.Nil(t, err)
	})
}

func TestClient_Get(t *testing.T) {
	t.Run("server returns error", func(t *testing.T) {
		resp, err := executor.client.Get(testCtxError, testToggleKey)

		assert.NotNil(t, err)
		assert.Nil(t, resp)
	})

	t.Run("toggle not found", func(t *testing.T) {
		resp, err := executor.client.Get(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Nil(t, resp)
	})

	t.Run("success get a toggle", func(t *testing.T) {
		resp, err := executor.client.Get(testCtxReturn, testToggleKey)
		fmt.Println(err)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}

func TestClient_Enable(t *testing.T) {
	t.Run("server returns error", func(t *testing.T) {
		err := executor.client.Enable(testCtxError, testToggleKey)

		assert.NotNil(t, err)
	})

	t.Run("success enable a toggle", func(t *testing.T) {
		err := executor.client.Enable(testCtx, testToggleKey)

		assert.Nil(t, err)
	})
}

func TestClient_Disable(t *testing.T) {
	t.Run("server returns error", func(t *testing.T) {
		err := executor.client.Disable(testCtxError, testToggleKey)

		assert.NotNil(t, err)
	})

	t.Run("success disable a toggle", func(t *testing.T) {
		err := executor.client.Disable(testCtx, testToggleKey)

		assert.Nil(t, err)
	})
}

func TestClient_Delete(t *testing.T) {
	t.Run("server returns error", func(t *testing.T) {
		err := executor.client.Delete(testCtxError, testToggleKey)

		assert.NotNil(t, err)
	})

	t.Run("success delete a toggle", func(t *testing.T) {
		err := executor.client.Delete(testCtx, testToggleKey)

		assert.Nil(t, err)
	})
}

func createClientExecutor() *ClientExecutor {
	listener := bufconn.Listen(buffer)

	grpcServer := grpc.NewServer()
	toggleServer := &mock_server.MockToggleServiceServer{}
	togglev1.RegisterToggleServiceServer(grpcServer, toggleServer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			panic(err)
		}
	}()

	client, err := toggle.NewClient(testCtxError,
		"",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithInsecure(),
	)

	if err != nil {
		grpcServer.Stop()
		panic(err)
	}

	closer := func() {
		_ = listener.Close()
		grpcServer.Stop()
	}

	return &ClientExecutor{
		client: client,
		closer: closer,
	}
}