package handler

import (
	"context"

	"github.com/indrasaputra/toggle/entity"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
	"github.com/indrasaputra/toggle/service"
)

// ToggleCommand handles HTTP/2 gRPC request for state-changing toggle .
type ToggleCommand struct {
	togglev1.UnimplementedToggleCommandServiceServer

	creator service.CreateToggle
	updater service.UpdateToggle
	deleter service.DeleteToggle
}

// NewToggleCommand creates an instance of ToggleCommand.
func NewToggleCommand(creator service.CreateToggle, updater service.UpdateToggle, deleter service.DeleteToggle) *ToggleCommand {
	return &ToggleCommand{
		creator: creator,
		updater: updater,
		deleter: deleter,
	}
}

// CreateToggle handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (tc *ToggleCommand) CreateToggle(ctx context.Context, request *togglev1.CreateToggleRequest) (*togglev1.CreateToggleResponse, error) {
	if request == nil || request.GetToggle() == nil {
		return nil, entity.ErrEmptyToggle()
	}

	err := tc.creator.Create(ctx, createToggleFromCreateToggleRequest(request))
	if err != nil {
		return nil, err
	}
	return &togglev1.CreateToggleResponse{}, nil
}

// EnableToggle handles HTTP/2 gRPC request similar to PUT in HTTP/1.1.
// It sets the toggle's is_enabled field to be true.
func (tc *ToggleCommand) EnableToggle(ctx context.Context, request *togglev1.EnableToggleRequest) (*togglev1.EnableToggleResponse, error) {
	if request == nil {
		return nil, entity.ErrEmptyToggle()
	}

	err := tc.updater.Enable(ctx, request.GetKey())
	if err != nil {
		return nil, err
	}
	return &togglev1.EnableToggleResponse{}, nil
}

// DisableToggle handles HTTP/2 gRPC request similar to PUT in HTTP/1.1.
// It sets the toggle's is_enabled field to be false.
func (tc *ToggleCommand) DisableToggle(ctx context.Context, request *togglev1.DisableToggleRequest) (*togglev1.DisableToggleResponse, error) {
	if request == nil {
		return nil, entity.ErrEmptyToggle()
	}

	err := tc.updater.Disable(ctx, request.GetKey())
	if err != nil {
		return nil, err
	}
	return &togglev1.DisableToggleResponse{}, nil
}

// DeleteToggle handles HTTP/2 gRPC request similar to DELETE in HTTP/1.1.
// It delete the toggle.
func (tc *ToggleCommand) DeleteToggle(ctx context.Context, request *togglev1.DeleteToggleRequest) (*togglev1.DeleteToggleResponse, error) {
	if request == nil {
		return nil, entity.ErrEmptyToggle()
	}

	err := tc.deleter.DeleteByKey(ctx, request.GetKey())
	if err != nil {
		return nil, err
	}
	return &togglev1.DeleteToggleResponse{}, nil
}

func createToggleFromCreateToggleRequest(request *togglev1.CreateToggleRequest) *entity.Toggle {
	return &entity.Toggle{
		Key:         request.GetToggle().GetKey(),
		Description: request.GetToggle().GetDescription(),
	}
}
