package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logsettable "github.com/grpc-ecosystem/go-grpc-middleware/logging/settable"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	connProtocol = "tcp"
)

// Grpc is responsible to act as gRPC server.
// It composes grpc.Server.
type Grpc struct {
	*grpc.Server
	listener net.Listener
	port     string
}

// newGrpc creates an instance of Grpc.
func newGrpc(port string, options ...grpc.ServerOption) *Grpc {
	srv := grpc.NewServer(options...)
	return &Grpc{
		Server: srv,
		port:   port,
	}
}

// NewGrpc creates an instance of Grpc for used in development environment.
//
// These are list of interceptors that are attached (from innermost to outermost):
// 	- Metrics, using Prometheus.
// 	- Logging, using zap logger.
// 	- Recoverer, using grpc_recovery.
func NewGrpc(port string) *Grpc {
	options := grpc_middleware.WithUnaryServerChain(defaultUnaryServerInterceptors()...)
	srv := newGrpc(port, options)
	grpc_prometheus.Register(srv.Server)
	return srv
}

// Run runs the server.
// It basically runs grpc.Server.Serve and is a blocking.
func (g *Grpc) Run() error {
	var err error
	g.listener, err = net.Listen(connProtocol, fmt.Sprintf(":%s", g.port))
	if err != nil {
		return err
	}

	go g.serve()
	log.Printf("grpc server is running on port %s\n", g.port)
	return nil
}

// AwaitTermination blocks the server and wait for termination signal.
// The termination signal must be one of SIGINT or SIGTERM.
// Once it receives one of those signals, the gRPC server will perform graceful stop and close the listener.
func (g *Grpc) AwaitTermination() error {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign

	g.GracefulStop()
	return g.listener.Close()
}

func (g *Grpc) serve() {
	if err := g.Serve(g.listener); err != nil {
		panic(err)
	}
}

func defaultUnaryServerInterceptors() []grpc.UnaryServerInterceptor {
	logger, _ := zap.NewProduction() // error is impossible, hence ignored.
	grpc_zap.SetGrpcLoggerV2(grpc_logsettable.ReplaceGrpcLoggerV2(), logger)
	grpc_prometheus.EnableHandlingTimeHistogram()

	options := []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(recoveryHandler)),
		grpc_zap.UnaryServerInterceptor(logger),
		grpc_prometheus.UnaryServerInterceptor,
		otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
	}
	return options
}

func recoveryHandler(p interface{}) error {
	return status.Errorf(codes.Unknown, "%v", p)
}
