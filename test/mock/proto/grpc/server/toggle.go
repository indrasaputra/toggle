// DO NOT EDIT.
// Umm.. actually you can edit this file :D
// The first sentence was only to avoid golint due to package name that uses underscore.

package mock_server

import (
	"context"

	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
)

// MockToggleServiceServer must be embedded to have forward compatible implementations.
type MockToggleServiceServer struct {
	togglev1.UnimplementedToggleServiceServer
}

func (MockToggleServiceServer) CreateToggle(context.Context, *togglev1.CreateToggleRequest) (*togglev1.CreateToggleResponse, error) {
	return &togglev1.CreateToggleResponse{}, nil
}

func (MockToggleServiceServer) GetToggleByKey(context.Context, *togglev1.GetToggleByKeyRequest) (*togglev1.GetToggleByKeyResponse, error) {
	return &togglev1.GetToggleByKeyResponse{}, nil
}

func (MockToggleServiceServer) GetAllToggles(context.Context, *togglev1.GetAllTogglesRequest) (*togglev1.GetAllTogglesResponse, error) {
	return &togglev1.GetAllTogglesResponse{}, nil
}

func (MockToggleServiceServer) Enable(context.Context, *togglev1.EnableRequest) (*togglev1.EnableResponse, error) {
	return &togglev1.EnableResponse{}, nil
}

func (MockToggleServiceServer) Disable(context.Context, *togglev1.DisableRequest) (*togglev1.DisableResponse, error) {
	return &togglev1.DisableResponse{}, nil
}

func (MockToggleServiceServer) DeleteToggle(context.Context, *togglev1.DeleteToggleRequest) (*togglev1.DeleteToggleResponse, error) {
	return &togglev1.DeleteToggleResponse{}, nil
}
