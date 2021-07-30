package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/indrasaputra/toggle/entity"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
	"github.com/indrasaputra/toggle/service"
)

// ToggleQuery handles HTTP/2 gRPC request for retrieve toggle .
type ToggleQuery struct {
	togglev1.UnimplementedToggleQueryServiceServer

	getter service.GetToggle
}

// NewToggleQuery creates an instance of ToggleQuery.
func NewToggleQuery(getter service.GetToggle) *ToggleQuery {
	return &ToggleQuery{getter: getter}
}

// GetToggleByKey handles HTTP/2 gRPC request similar to GET in HTTP/1.1.
// It gets a single toggle by its key.
func (tq *ToggleQuery) GetToggleByKey(ctx context.Context, request *togglev1.GetToggleByKeyRequest) (*togglev1.GetToggleByKeyResponse, error) {
	if request == nil {
		return nil, entity.ErrEmptyToggle()
	}

	toggle, err := tq.getter.GetByKey(ctx, request.GetKey())
	if err != nil {
		return nil, err
	}
	return createGetToggleByKeyResponse(toggle), nil
}

// GetAllToggles handles HTTP/2 gRPC request similar to GET in HTTP/1.1.
// It gets all available toggles in system.
func (tq *ToggleQuery) GetAllToggles(ctx context.Context, request *togglev1.GetAllTogglesRequest) (*togglev1.GetAllTogglesResponse, error) {
	if request == nil {
		return nil, entity.ErrEmptyToggle()
	}

	toggles, err := tq.getter.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return createGetAllTogglesResponse(toggles), nil
}

func createGetToggleByKeyResponse(toggle *entity.Toggle) *togglev1.GetToggleByKeyResponse {
	return &togglev1.GetToggleByKeyResponse{
		Toggle: createProtoToggle(toggle),
	}
}

func createGetAllTogglesResponse(toggles []*entity.Toggle) *togglev1.GetAllTogglesResponse {
	resp := &togglev1.GetAllTogglesResponse{}
	for _, toggle := range toggles {
		resp.Toggles = append(resp.Toggles, createProtoToggle(toggle))
	}
	return resp
}

func createProtoToggle(toggle *entity.Toggle) *togglev1.Toggle {
	return &togglev1.Toggle{
		Key:         toggle.Key,
		IsEnabled:   toggle.IsEnabled,
		Description: toggle.Description,
		CreatedAt:   timestamppb.New(toggle.CreatedAt),
		UpdatedAt:   timestamppb.New(toggle.UpdatedAt),
	}
}
