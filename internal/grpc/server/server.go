package server

import (
	"fmt"
	"net"

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

	"github.com/indrasaputra/toggle/internal/grpc/interceptor"
)

const (
	connProtocol   = "tcp"
	grpcServerName = "grpc server"
)

// GrpcServer is responsible to act as gRPC server.
// It composes grpc.Server.
type GrpcServer struct {
	server      *grpc.Server
	serviceFunc []func(*grpc.Server)
	listener    net.Listener
	port        string
}

// newGrpc creates an instance of GrpcServer.
func newGrpcServer(port string, options ...grpc.ServerOption) *GrpcServer {
	return &GrpcServer{
		server: grpc.NewServer(options...),
		port:   port,
	}
}

// NewGrpcServer creates an instance of GrpcServer for used in development environment.
//
// These are list of interceptors that are attached (from innermost to outermost):
// 	- Metrics, using Prometheus.
// 	- Logging, using zap logger.
// 	- Recoverer, using grpc_recovery.
func NewGrpcServer(port string) *GrpcServer {
	options := grpc_middleware.WithUnaryServerChain(defaultUnaryServerInterceptors()...)
	srv := newGrpcServer(port, options)
	grpc_prometheus.Register(srv.server)
	return srv
}

// Name returns server's name.
func (gs *GrpcServer) Name() string {
	return grpcServerName
}

// Port returns server's port.
func (gs *GrpcServer) Port() string {
	return gs.port
}

// AttachService attaches service to gRPC server.
// It will be called before serve.
func (gs *GrpcServer) AttachService(fn func(*grpc.Server)) {
	gs.serviceFunc = append(gs.serviceFunc, fn)
}

// Serve runs the server.
// It basically runs grpc.Server.Serve and is a blocking.
func (gs *GrpcServer) Serve() error {
	for _, service := range gs.serviceFunc {
		service(gs.server)
	}

	var err error
	gs.listener, err = net.Listen(connProtocol, fmt.Sprintf(":%s", gs.port))
	if err != nil {
		return err
	}
	return gs.server.Serve(gs.listener)
}

// GracefulStop blocks the server and wait for termination signal.
// The termination signal must be one of SIGINT or SIGTERM.
// Once it receives one of those signals, the gRPC server will perform graceful stop and close the listener.
//
// It receives Closer and will perform all closers before closing itself.
func (gs *GrpcServer) GracefulStop() {
	gs.server.GracefulStop()
	if gs.listener != nil {
		_ = gs.listener.Close()
	}
}

// Stop immediately stops the gRPC server.
// Currently, this method exists just for the sake of testing.
// For production purpose, use GracefulStop().
func (gs *GrpcServer) Stop() {
	gs.server.Stop()
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
		interceptor.OpenTracingUnaryServerInterceptor(),
	}
	return options
}

func recoveryHandler(p interface{}) error {
	return status.Errorf(codes.Unknown, "%v", p)
}
