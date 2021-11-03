package server_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/internal/server"
	mock_server "github.com/indrasaputra/toggle/test/mock/server"
)

const (
	serverName = "mock server"
	serverPort = "8080"
)

type ManagerExecutor struct {
	manager *server.Manager
	servers []*mock_server.MockServer
}

func TestNewManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Manager", func(t *testing.T) {
		t.Parallel()

		exec := createManagerExecutor(ctrl)
		assert.NotNil(t, exec.manager)
	})
}

func TestManager_Serve(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully run all servers", func(t *testing.T) {
		t.Parallel()

		exec := createManagerExecutor(ctrl)
		for i := range exec.servers {
			exec.servers[i].EXPECT().Name().Return(serverName)
			exec.servers[i].EXPECT().Port().Return(serverPort)
			exec.servers[i].EXPECT().Serve().Return(nil)
		}

		assert.NotPanics(t, func() { exec.manager.Serve() })
	})

	t.Run("all servers are fail", func(t *testing.T) {
		t.Parallel()

		exec := createManagerExecutor(ctrl)
		for i := range exec.servers {
			exec.servers[i].EXPECT().Name().AnyTimes().Return(serverName)
			exec.servers[i].EXPECT().Port().Return(serverPort)
			exec.servers[i].EXPECT().Serve().Return(errors.New("shutdown"))
		}

		assert.NotPanics(t, func() { exec.manager.Serve() })
	})
}

func createManagerExecutor(ctrl *gomock.Controller) *ManagerExecutor {
	s1 := mock_server.NewMockServer(ctrl)
	s2 := mock_server.NewMockServer(ctrl)
	m := server.NewManager([]server.Server{s1, s2})

	return &ManagerExecutor{
		manager: m,
		servers: []*mock_server.MockServer{s1, s2},
	}
}
