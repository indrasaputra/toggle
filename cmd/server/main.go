package main

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/indrasaputra/toggle/internal/app"
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

	postgrePool, err := builder.BuildPostgrePgxPool(&cfg.Postgres)
	checkError(err)
	redisClient, err := builder.BuildRedisClient(&cfg.Redis)
	checkError(err)
	kafkaWriter := builder.BuildKafkaWriter(&cfg.Kafka)

	tracerProvider, err := app.InitTracer(cfg)
	checkError(err)

	dep := &builder.Dependency{
		PgxPool:     postgrePool,
		RedisClient: redisClient,
		KafkaWriter: kafkaWriter,
		Config:      cfg,
	}

	grpcServer := grpcserver.NewGrpcServer(cfg.Port.Grpc)
	registerGrpcService(grpcServer, dep)

	gatewayServer := gwayserver.NewGrpcGateway(cfg.Port.GrpcGateway)
	registerGrpcGatewayService(context.Background(), gatewayServer, fmt.Sprintf(":%s", cfg.Port.Grpc), grpc.WithInsecure())

	closer := func() {
		_ = tracerProvider.Shutdown(context.Background())
		_ = kafkaWriter.Close()
		_ = redisClient.Close()
		postgrePool.Close()
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

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
