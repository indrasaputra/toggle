package server_test

import (
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/internal/grpc-gateway/server"
)

var (
	testGrpcGatewayPort = "8081"
)

func TestNewGrpcGateway(t *testing.T) {
	t.Run("success create rest server", func(t *testing.T) {
		srv := server.NewGrpcGateway(testGrpcGatewayPort)
		assert.NotNil(t, srv)
	})
}

func TestGrpcGateway_EnablePrometheus(t *testing.T) {
	t.Run("success enable prometheus", func(t *testing.T) {
		srv := server.NewGrpcGateway(testGrpcGatewayPort)
		err := srv.EnablePrometheus()
		assert.Nil(t, err)
	})
}

func TestGrpcGateway_EnableHealth(t *testing.T) {
	t.Run("success enable health", func(t *testing.T) {
		srv := server.NewGrpcGateway(testGrpcGatewayPort)
		err := srv.EnableHealth()
		assert.Nil(t, err)
	})
}

func TestNewGrpcGateway_AttachService(t *testing.T) {
	t.Run("success attach service to server", func(t *testing.T) {
		fn := func(s *runtime.ServeMux) error {
			return nil
		}

		srv := server.NewGrpcGateway(testGrpcGatewayPort)

		assert.NotPanics(t, func() { srv.AttachService(fn) })
	})
}

func TestNewGrpcGateway_Name(t *testing.T) {
	t.Run("success get server's name", func(t *testing.T) {
		srv := server.NewGrpcGateway(testGrpcGatewayPort)
		assert.Equal(t, "grpc-gateway server", srv.Name())
	})
}

func TestNewGrpcGateway_Port(t *testing.T) {
	t.Run("success get server's port", func(t *testing.T) {
		srv := server.NewGrpcGateway(testGrpcGatewayPort)
		assert.Equal(t, testGrpcGatewayPort, srv.Port())
	})
}

func TestNewGrpcGateway_GracefulStop(t *testing.T) {
	t.Run("success stop server", func(t *testing.T) {
		srv := server.NewGrpcGateway(testGrpcGatewayPort)
		assert.NotPanics(t, func() { srv.GracefulStop() })
	})
}
