package server_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/indrasaputra/toggle/internal/grpc/server"
)

var (
	testGrpcPort = "8080"
)

func TestNewGrpcServer(t *testing.T) {
	t.Run("successfully create a development gRPC server", func(t *testing.T) {
		srv := server.NewGrpcServer(testGrpcPort)
		assert.NotNil(t, srv)
	})
}

func TestGrpcServer_Serve(t *testing.T) {
	t.Run("success run", func(t *testing.T) {
		srv := server.NewGrpcServer(testGrpcPort)

		stopServerAfterOneSecond(srv)
		err := srv.Serve()

		assert.Nil(t, err)
	})

	t.Run("success run with attached service", func(t *testing.T) {
		srv := server.NewGrpcServer(testGrpcPort)
		srv.AttachService(func(s *grpc.Server) {})

		stopServerAfterOneSecond(srv)
		err := srv.Serve()

		assert.Nil(t, err)
	})
}

func TestGrpcServer_GracefulStop(t *testing.T) {
	t.Run("stop server without listener", func(t *testing.T) {
		srv := server.NewGrpcServer(testGrpcPort)
		assert.NotPanics(t, func() { srv.GracefulStop() })
	})

	t.Run("stop server after serve", func(t *testing.T) {
		srv := server.NewGrpcServer(testGrpcPort)

		stopServerGracefullyAfterOneSecond(srv)
		err := srv.Serve()

		assert.Nil(t, err)
	})
}

func TestGrpcServer_AttachService(t *testing.T) {
	t.Run("success attach service to server", func(t *testing.T) {
		fn := func(s *grpc.Server) {}
		srv := server.NewGrpcServer(testGrpcPort)
		assert.NotPanics(t, func() { srv.AttachService(fn) })
	})
}

func TestGrpcServer_Name(t *testing.T) {
	t.Run("success get server's name", func(t *testing.T) {
		srv := server.NewGrpcServer(testGrpcPort)
		assert.Equal(t, "grpc server", srv.Name())
	})
}

func TestGrpcServer_Port(t *testing.T) {
	t.Run("success get server's port", func(t *testing.T) {
		srv := server.NewGrpcServer(testGrpcPort)
		assert.Equal(t, testGrpcPort, srv.Port())
	})
}

func stopServerAfterOneSecond(srv *server.GrpcServer) {
	go func() {
		time.Sleep(1 * time.Second)
		srv.Stop()
	}()
}

func stopServerGracefullyAfterOneSecond(srv *server.GrpcServer) {
	go func() {
		time.Sleep(1 * time.Second)
		srv.GracefulStop()
	}()
}
