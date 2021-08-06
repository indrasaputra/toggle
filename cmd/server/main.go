package main

import (
	"context"
	"fmt"
	"io"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"github.com/uber/jaeger-client-go"
	jaegerconf "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/indrasaputra/toggle/internal/builder"
	"github.com/indrasaputra/toggle/internal/config"
	"github.com/indrasaputra/toggle/internal/grpc/handler"
	"github.com/indrasaputra/toggle/internal/server"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
)

func main() {
	cfg, cerr := config.NewConfig(".env")
	checkError(cerr)

	psql, perr := builder.BuildPgxPool(&cfg.Postgres)
	checkError(perr)
	rds, rerr := builder.BuildRedisClient(&cfg.Redis)
	checkError(rerr)
	kfk := builder.BuildKafkaWriter(&cfg.Kafka)
	trc := initTracing(cfg)

	grpcServer := server.NewGrpc(cfg.Port.GRPC)
	registerGrpcHandlers(grpcServer.Server, psql, rds, kfk, cfg)

	restServer := server.NewRest(cfg.Port.REST)
	registerRestHandlers(context.Background(), restServer.ServeMux, fmt.Sprintf(":%s", cfg.Port.GRPC), grpc.WithInsecure())

	closer := func() {
		_ = trc.Close()
		_ = kfk.Close()
		_ = rds.Close()
		psql.Close()
	}

	_ = grpcServer.Run()
	_ = restServer.Run()
	_ = grpcServer.AwaitTermination(closer)
}

func registerGrpcHandlers(server *grpc.Server, psql *pgxpool.Pool, rds *goredis.Client, kfk *kafka.Writer, cfg *config.Config) {
	// start register all module's gRPC handlers
	command := builder.BuildToggleCommandHandler(psql, rds, time.Duration(cfg.Redis.TTL)*time.Minute, kfk)
	togglev1.RegisterToggleCommandServiceServer(server, command)
	query := builder.BuildToggleQueryHandler(psql, rds, time.Duration(cfg.Redis.TTL)*time.Minute)
	togglev1.RegisterToggleQueryServiceServer(server, query)

	health := handler.NewHealth()
	grpc_health_v1.RegisterHealthServer(server, health)
	// end of register all module's gRPC handlers
}

func registerRestHandlers(ctx context.Context, server *runtime.ServeMux, grpcPort string, options ...grpc.DialOption) {
	// start register all module's REST handlers
	err := togglev1.RegisterToggleCommandServiceHandlerFromEndpoint(ctx, server, grpcPort, options)
	checkError(err)
	err = togglev1.RegisterToggleQueryServiceHandlerFromEndpoint(ctx, server, grpcPort, options)
	checkError(err)
	// end of register all module's REST handlers
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
