package handler

import (
	"context"

	"github.com/indrasaputra/toggle/entity"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
	"github.com/indrasaputra/toggle/service"
)

// Toggle handles HTTP/2 gRPC request for toggle .
type Toggle struct {
	togglev1.UnimplementedToggleServiceServer
	creator service.CreateToggle
}

// NewToggle creates an instance of Toggle.
func NewToggle(creator service.CreateToggle) *Toggle {
	return &Toggle{
		creator: creator,
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

func createToggleFromCreateToggleRequest(request *togglev1.CreateToggleRequest) *entity.Toggle {
	return &entity.Toggle{
		Key: request.GetKey(),
	}
}
