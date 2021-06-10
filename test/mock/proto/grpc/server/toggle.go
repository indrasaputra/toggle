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
	togglev1.UnimplementedToggleServiceServer
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
		return &togglev1.GetToggleByKeyResponse{Toggle: &togglev1.Toggle{}}, nil
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

func (MockToggleServiceServer) Enable(ctx context.Context, _ *togglev1.EnableRequest) (*togglev1.EnableResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md[keyErr]) > 0 && md[keyErr][0] != "" {
		return nil, errInternal
	}
	return &togglev1.EnableResponse{}, nil
}

func (MockToggleServiceServer) Disable(ctx context.Context, _ *togglev1.DisableRequest) (*togglev1.DisableResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md[keyErr]) > 0 && md[keyErr][0] != "" {
		return nil, errInternal
	}
	return &togglev1.DisableResponse{}, nil
}

func (MockToggleServiceServer) DeleteToggle(ctx context.Context, _ *togglev1.DeleteToggleRequest) (*togglev1.DeleteToggleResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md[keyErr]) > 0 && md[keyErr][0] != "" {
		return nil, errInternal
	}
	return &togglev1.DeleteToggleResponse{}, nil
}
