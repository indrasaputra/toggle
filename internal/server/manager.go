package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Server defines contract to implement server.
type Server interface {
	// Name represents server's name.
	Name() string
	// Port represents server's port.
	Port() string
	// Serve runs the server in a blocking way.
	// It is up to implementor to make it run in a goroutine so that it doesn't block
	// or just let it be a blocking method.
	Serve() error
	// GracefulStop stops the server gracefully.
	GracefulStop()
}

// Manager manages the attached servers.
type Manager struct {
	servers []Server
}

// NewManager creates an instance of Manager.
func NewManager(servers []Server) *Manager {
	return &Manager{
		servers: servers,
	}
}

// Serve runs all attached servers.
// Each server will be run in an independent goroutine
// to make sure that no server blocking each others.
func (m *Manager) Serve() {
	for _, server := range m.servers {
		go func(srv Server) {
			log.Printf("%s is running on port %s\n", srv.Name(), srv.Port())
			if err := srv.Serve(); err != nil {
				log.Printf("%s got error: %v\n", srv.Name(), err)
			}
		}(server)
	}
}

// GracefulStop stops all servers gracefully.
// It waits for signal which currently implemented as signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM).
func (m *Manager) GracefulStop() {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign

	for _, server := range m.servers {
		server.GracefulStop()
		log.Printf("%s has been stopped\n", server.Name())
	}
}
