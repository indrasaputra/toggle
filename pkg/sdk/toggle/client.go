package toggle

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/indrasaputra/toggle/entity"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
)

// Client acts as a client to connect to Toggle.
type Client struct {
	command togglev1.ToggleCommandServiceClient
	query   togglev1.ToggleQueryServiceClient
}

// NewClient creates an instance of Client.
func NewClient(ctx context.Context, host string, options ...grpc.DialOption) (*Client, error) {
	conn, err := grpc.DialContext(ctx, host, options...)
	if err != nil {
		return nil, status.New(codes.Unavailable, "").Err()
	}

	return &Client{
		command: togglev1.NewToggleCommandServiceClient(conn),
		query:   togglev1.NewToggleQueryServiceClient(conn),
	}, nil
}

// Create creates a new toggle.
func (c *Client) Create(ctx context.Context, toggle *entity.Toggle) error {
	req := &togglev1.CreateToggleRequest{Toggle: &togglev1.Toggle{
		Key:         toggle.Key,
		Description: toggle.Description,
	}}
	_, err := c.command.CreateToggle(ctx, req)
	return err
}

// Get gets a single toggle by its key.
func (c *Client) Get(ctx context.Context, key string) (*entity.Toggle, error) {
	req := &togglev1.GetToggleByKeyRequest{Key: key}
	resp, err := c.query.GetToggleByKey(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.GetToggle() == nil {
		return nil, entity.ErrNotFound()
	}

	toggle := &entity.Toggle{
		Key:         resp.GetToggle().GetKey(),
		IsEnabled:   resp.GetToggle().GetIsEnabled(),
		Description: resp.GetToggle().GetDescription(),
		CreatedAt:   resp.GetToggle().GetCreatedAt().AsTime(),
		UpdatedAt:   resp.GetToggle().GetUpdatedAt().AsTime(),
	}
	return toggle, nil
}

// Enable enables a toggle.
// It sets toggle's `is_enabled` attribute to be true.
func (c *Client) Enable(ctx context.Context, key string) error {
	req := &togglev1.EnableToggleRequest{Key: key}
	_, err := c.command.EnableToggle(ctx, req)
	return err
}

// Disable disables a toggle.
// It sets toggle's `is_enabled` attribute to be false.
func (c *Client) Disable(ctx context.Context, key string) error {
	req := &togglev1.DisableToggleRequest{Key: key}
	_, err := c.command.DisableToggle(ctx, req)
	return err
}

// Delete deletes a toggle.
// It only deletes a nonactive toggle (is_enabled == false).
func (c *Client) Delete(ctx context.Context, key string) error {
	req := &togglev1.DeleteToggleRequest{Key: key}
	_, err := c.command.DeleteToggle(ctx, req)
	return err
}
