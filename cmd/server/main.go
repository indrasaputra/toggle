package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
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

	grpcServer := server.NewGrpc(cfg.Port.GRPC)
	registerGrpcHandlers(grpcServer.Server, psql, rds, cfg)

	restServer := server.NewRest(cfg.Port.REST)
	registerRestHandlers(context.Background(), restServer.ServeMux, fmt.Sprintf(":%s", cfg.Port.GRPC), grpc.WithInsecure())

	_ = grpcServer.Run()
	_ = restServer.Run()
	_ = grpcServer.AwaitTermination()
}

func health(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func registerGrpcHandlers(server *grpc.Server, psql *pgxpool.Pool, rds *goredis.Client, cfg *config.Config) {
	// start register all module's gRPC handlers
	toggle := builder.BuildToggleHandler(psql, rds, time.Duration(cfg.Redis.TTL)*time.Minute)
	togglev1.RegisterToggleServiceServer(server, toggle)

	health := handler.NewHealth()
	grpc_health_v1.RegisterHealthServer(server, health)
	// end of register all module's gRPC handlers
}

func registerRestHandlers(ctx context.Context, server *runtime.ServeMux, grpcPort string, options ...grpc.DialOption) {
	// start register all module's REST handlers
	err := togglev1.RegisterToggleServiceHandlerFromEndpoint(ctx, server, grpcPort, options)
	checkError(err)
	// end of register all module's REST handlers
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
