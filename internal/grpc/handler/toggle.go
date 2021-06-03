package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/indrasaputra/toggle/entity"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
	"github.com/indrasaputra/toggle/service"
)

// Toggle handles HTTP/2 gRPC request for toggle .
type Toggle struct {
	togglev1.UnimplementedToggleServiceServer

	creator service.CreateToggle
	getter  service.GetToggle
	deleter service.DeleteToggle
}

// NewToggle creates an instance of Toggle.
func NewToggle(creator service.CreateToggle, getter service.GetToggle, deleter service.DeleteToggle) *Toggle {
	return &Toggle{
		creator: creator,
		getter:  getter,
		deleter: deleter,
	}
}

// CreateToggle handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (th *Toggle) CreateToggle(ctx context.Context, request *togglev1.CreateToggleRequest) (*togglev1.CreateToggleResponse, error) {
	if request == nil {
		return nil, entity.ErrEmptyToggle()
	}

	err := th.creator.Create(ctx, createToggleFromCreateToggleRequest(request))
	if err != nil {
		return nil, err
	}
	return &togglev1.CreateToggleResponse{}, nil
}

// GetToggleByKey handles HTTP/2 gRPC request similar to GET in HTTP/1.1.
// It gets a single toggle by its key.
func (th *Toggle) GetToggleByKey(ctx context.Context, request *togglev1.GetToggleByKeyRequest) (*togglev1.GetToggleByKeyResponse, error) {
	if request == nil {
		return nil, entity.ErrEmptyToggle()
	}

	toggle, err := th.getter.GetByKey(ctx, request.GetKey())
	if err != nil {
		return nil, err
	}
	return createGetToggleByKeyResponse(toggle), nil
}

// GetAllToggles handles HTTP/2 gRPC request similar to GET in HTTP/1.1.
// It gets all available toggles in system.
func (th *Toggle) GetAllToggles(ctx context.Context, request *togglev1.GetAllTogglesRequest) (*togglev1.GetAllTogglesResponse, error) {
	if request == nil {
		return nil, entity.ErrEmptyToggle()
	}

	toggles, err := th.getter.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return createGetAllTogglesResponse(toggles), nil
}

// DeleteToggle handles HTTP/2 gRPC request similar to DELETE in HTTP/1.1.
// It delete the toggle.
func (th *Toggle) DeleteToggle(ctx context.Context, request *togglev1.DeleteToggleRequest) (*togglev1.DeleteToggleResponse, error) {
	if request == nil {
		return nil, entity.ErrEmptyToggle()
	}

	err := th.deleter.DeleteByKey(ctx, request.GetKey())
	if err != nil {
		return nil, err
	}
	return &togglev1.DeleteToggleResponse{}, nil
}

func createToggleFromCreateToggleRequest(request *togglev1.CreateToggleRequest) *entity.Toggle {
	return &entity.Toggle{
		Key:         request.GetKey(),
		Description: request.GetDescription(),
	}
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
