// DO NOT EDIT.
// Umm.. actually you can edit this file :D
// The first sentence was only to avoid golint due to package name that uses underscore.

package mock_server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
)

var (
	keyErr      = "has-error"
	keyReturn   = "complete-return"
	errInternal = status.New(codes.Internal, "").Err()
)

// MockToggleServiceServer must be embedded to have forward compatible implementations.
type MockToggleServiceServer struct {
	togglev1.UnimplementedToggleCommandServiceServer
	togglev1.UnimplementedToggleQueryServiceServer
}

func (MockToggleServiceServer) CreateToggle(ctx context.Context, _ *togglev1.CreateToggleRequest) (*togglev1.CreateToggleResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md[keyErr]) > 0 && md[keyErr][0] != "" {
		return nil, errInternal
	}
	return &togglev1.CreateToggleResponse{}, nil
}

func (MockToggleServiceServer) GetToggleByKey(ctx context.Context, _ *togglev1.GetToggleByKeyRequest) (*togglev1.GetToggleByKeyResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md[keyErr]) > 0 && md[keyErr][0] != "" {
		return nil, errInternal
	}
	if len(md[keyReturn]) > 0 && md[keyReturn][0] != "" {
		return &togglev1.GetToggleByKeyResponse{Toggle: &togglev1.Toggle{Key: md[keyReturn][0], IsEnabled: true}}, nil
	}
	return &togglev1.GetToggleByKeyResponse{}, nil
}

func (MockToggleServiceServer) GetAllToggles(ctx context.Context, _ *togglev1.GetAllTogglesRequest) (*togglev1.GetAllTogglesResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md[keyErr]) > 0 && md[keyErr][0] != "" {
		return nil, errInternal
	}
	return &togglev1.GetAllTogglesResponse{}, nil
}

func (MockToggleServiceServer) EnableToggle(ctx context.Context, _ *togglev1.EnableToggleRequest) (*togglev1.EnableToggleResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md[keyErr]) > 0 && md[keyErr][0] != "" {
		return nil, errInternal
	}
	return &togglev1.EnableToggleResponse{}, nil
}

func (MockToggleServiceServer) DisableToggle(ctx context.Context, _ *togglev1.DisableToggleRequest) (*togglev1.DisableToggleResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md[keyErr]) > 0 && md[keyErr][0] != "" {
		return nil, errInternal
	}
	return &togglev1.DisableToggleResponse{}, nil
}

func (MockToggleServiceServer) DeleteToggle(ctx context.Context, _ *togglev1.DeleteToggleRequest) (*togglev1.DeleteToggleResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md[keyErr]) > 0 && md[keyErr][0] != "" {
		return nil, errInternal
	}
	return &togglev1.DeleteToggleResponse{}, nil
}
