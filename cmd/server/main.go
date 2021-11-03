package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerconf "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/indrasaputra/toggle/internal/builder"
	"github.com/indrasaputra/toggle/internal/config"
	gwayserver "github.com/indrasaputra/toggle/internal/grpc-gateway/server"
	"github.com/indrasaputra/toggle/internal/grpc/handler"
	grpcserver "github.com/indrasaputra/toggle/internal/grpc/server"
	manserver "github.com/indrasaputra/toggle/internal/server"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
)

func main() {
	cfg, err := config.NewConfig(".env")
	checkError(err)

	dbpool, err := builder.BuildPgxPool(&cfg.Postgres)
	checkError(err)
	rds, err := builder.BuildRedisClient(&cfg.Redis)
	checkError(err)
	kfk := builder.BuildKafkaWriter(&cfg.Kafka)
	trc := initTracing(cfg)

	dep := &builder.Dependency{
		PgxPool:     dbpool,
		RedisClient: rds,
		KafkaWriter: kfk,
		Config:      cfg,
	}

	grpcServer := grpcserver.NewGrpcServer(cfg.Port.Grpc)
	registerGrpcService(grpcServer, dep)

	gatewayServer := gwayserver.NewGrpcGateway(cfg.Port.GrpcGateway)
	registerGrpcGatewayService(context.Background(), gatewayServer, fmt.Sprintf(":%s", cfg.Port.Grpc), grpc.WithInsecure())

	closer := func() {
		_ = trc.Close()
		_ = kfk.Close()
		_ = rds.Close()
		dbpool.Close()
	}
	defer closer()

	man := manserver.NewManager([]manserver.Server{grpcServer, gatewayServer})
	man.Serve()
	man.GracefulStop()
}

func registerGrpcService(grpcServer *grpcserver.GrpcServer, dep *builder.Dependency) {
	// start register all module's gRPC handlers
	command := builder.BuildToggleCommandHandler(dep)
	query := builder.BuildToggleQueryHandler(dep)
	health := handler.NewHealth()

	grpcServer.AttachService(func(server *grpc.Server) {
		togglev1.RegisterToggleCommandServiceServer(server, command)
		togglev1.RegisterToggleQueryServiceServer(server, query)
		grpc_health_v1.RegisterHealthServer(server, health)
	})
	// end of register all module's gRPC handlers
}

func registerGrpcGatewayService(ctx context.Context, gatewayServer *gwayserver.GrpcGateway, grpcPort string, options ...grpc.DialOption) {
	gatewayServer.AttachService(func(server *runtime.ServeMux) error {
		if err := togglev1.RegisterToggleCommandServiceHandlerFromEndpoint(ctx, server, grpcPort, options); err != nil {
			return err
		}
		if err := togglev1.RegisterToggleQueryServiceHandlerFromEndpoint(ctx, server, grpcPort, options); err != nil {
			return err
		}
		return nil
	})
}

func initTracing(cfg *config.Config) io.Closer {
	if !cfg.Jaeger.Enabled {
		return nopCloser{}
	}

	jaegerCfg := &jaegerconf.Configuration{
		ServiceName: cfg.ServiceName,
		Sampler: &jaegerconf.SamplerConfig{
			Type:  cfg.Jaeger.SamplingType,
			Param: cfg.Jaeger.SamplingParam,
		},
		Reporter: &jaegerconf.ReporterConfig{
			LogSpans:            cfg.Jaeger.LogSpans,
			LocalAgentHostPort:  fmt.Sprintf("%s:%s", cfg.Jaeger.Host, cfg.Jaeger.Port),
			BufferFlushInterval: time.Duration(cfg.Jaeger.FlushInterval) * time.Second,
		},
	}
	tracer, closer, err := jaegerCfg.NewTracer(jaegerconf.Logger(jaeger.StdLogger))
	checkError(err)

	opentracing.SetGlobalTracer(tracer)
	return closer
}

type nopCloser struct{}

// Closer closes nothing.
func (nopCloser) Close() error { return nil }

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
