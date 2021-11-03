package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	grpcGatewayServerName = "grpc-gateway server"
)

// GrpcGateway is responsible to act as HTTP/1.1 server.
// It composes grpc-gateway runtime.ServeMux.
type GrpcGateway struct {
	mux         *runtime.ServeMux
	serviceFunc []func(*runtime.ServeMux) error
	port        string
}

// NewGrpcGateway creates an instance of GrpcGateway with default production options attached.
// It enables Prometheus metrics by default.
func NewGrpcGateway(port string) *GrpcGateway {
	srv := &GrpcGateway{
		mux:  runtime.NewServeMux(),
		port: port,
	}
	_ = srv.EnablePrometheus() // error is impossible, hence ignored.
	_ = srv.EnableHealth()     // error is impossible, hence ignored.
	return srv
}

// Name returns server's name.
func (gg *GrpcGateway) Name() string {
	return grpcGatewayServerName
}

// Port returns server's port.
func (gg *GrpcGateway) Port() string {
	return gg.port
}

// EnablePrometheus enables prometheus endpoint.
// It can be accessed via /metrics.
func (gg *GrpcGateway) EnablePrometheus() error {
	return gg.mux.HandlePath(http.MethodGet, "/metrics", prometheusHandler())
}

// EnableHealth enables health endpoint.
// It can be accessed via /health.
func (gg *GrpcGateway) EnableHealth() error {
	return gg.mux.HandlePath(http.MethodGet, "/health", healthHandler())
}

// Serve runs HTTP/1.1 runtime.ServeMux.
// It is a blocking method.
func (gg *GrpcGateway) Serve() error {
	for _, service := range gg.serviceFunc {
		if err := service(gg.mux); err != nil {
			return err
		}
	}
	return http.ListenAndServe(fmt.Sprintf(":%s", gg.port), allowCORS(gg.mux))
}

// AttachService attaches service to gRPC Gateway server.
// It will be called before serve.
func (gg *GrpcGateway) AttachService(fn func(*runtime.ServeMux) error) {
	gg.serviceFunc = append(gg.serviceFunc, fn)
}

// GracefulStop exists just for the sake implementing server interface.
// It does nothing.
func (gg *GrpcGateway) GracefulStop() {
	// It does nothing. This comment exists due to the "Functions should not be empty" rule.
}

func prometheusHandler() runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		promhttp.Handler().ServeHTTP(w, r)
	}
}

func healthHandler() runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		w.WriteHeader(http.StatusOK)
	}
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				headers := []string{"Content-Type", "Accept"}
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
				methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
