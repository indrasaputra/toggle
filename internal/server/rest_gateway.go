package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Rest is responsible to act as HTTP/1.1 REST server.
// It composes grpc-gateway runtime.ServeMux.
type Rest struct {
	*runtime.ServeMux
	port string
}

// NewRest creates an instance of Rest with default production options attached.
// It enables Prometheus metrics by default.
func NewRest(port string) *Rest {
	srv := &Rest{
		ServeMux: runtime.NewServeMux(),
		port:     port,
	}
	_ = srv.EnablePrometheus() // error is impossible, hence ignored.
	return srv
}

// EnablePrometheus enables prometheus endpoint.
// It can be accessed via /metrics.
func (r *Rest) EnablePrometheus() error {
	return r.ServeMux.HandlePath(http.MethodGet, "/metrics", prometheusHandler())
}

// Run runs HTTP/1.1 runtime.ServeMux.
// It runs inside a goroutine.
func (r *Rest) Run() error {
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", r.port), allowCORS(r.ServeMux)); err != nil {
			panic(err)
		}
	}()
	return nil
}

func prometheusHandler() runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		promhttp.Handler().ServeHTTP(w, r)
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
