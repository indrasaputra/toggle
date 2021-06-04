package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// Health handles HTTP/2 gRPC request for health checking.
type Health struct {
	grpc_health_v1.UnimplementedHealthServer
}

// NewHealth creates an instance of Health.
func NewHealth() *Health {
	return &Health{}
}

// Check checks the entire system health.
func (hh *Health) Check(ctx context.Context, request *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	if request == nil {
		st := status.New(codes.InvalidArgument, "health check request is nil")
		return createHealthCheckResponse(grpc_health_v1.HealthCheckResponse_UNKNOWN), st.Err()
	}
	return createHealthCheckResponse(grpc_health_v1.HealthCheckResponse_SERVING), nil
}

func createHealthCheckResponse(status grpc_health_v1.HealthCheckResponse_ServingStatus) *grpc_health_v1.HealthCheckResponse {
	return &grpc_health_v1.HealthCheckResponse{
		Status: status,
	}
}
