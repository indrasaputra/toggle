package toggle

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/indrasaputra/toggle/entity"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
)

// DialConfig defines configuration to work with Client.
type DialConfig struct {
	// Host defines server host.
	Host string
	// Options defines list of dial option used to make a connection to server.
	Options []grpc.DialOption
}

// CircuitBreaker defines interface for circuit breaker.
type CircuitBreaker interface {
	// Execute executes the given function parameter.
	// It must implement circuit breaker concept during execution.
	Execute(func() (interface{}, error)) (interface{}, error)
}

// Client acts as a client to connect to Toggle.
type Client struct {
	command togglev1.ToggleCommandServiceClient
	query   togglev1.ToggleQueryServiceClient

	breaker CircuitBreaker
}

// NewClient creates an instance of Client.
func NewClient(cfg *DialConfig, breaker CircuitBreaker) (*Client, error) {
	conn, err := grpc.DialContext(context.Background(), cfg.Host, cfg.Options...)
	if err != nil {
		return nil, status.New(codes.Unavailable, "").Err()
	}

	if breaker == nil {
		breaker = noBreaker{}
	}

	return &Client{
		command: togglev1.NewToggleCommandServiceClient(conn),
		query:   togglev1.NewToggleQueryServiceClient(conn),
		breaker: breaker,
	}, nil
}

// Create creates a new toggle.
func (c *Client) Create(ctx context.Context, toggle *entity.Toggle) error {
	req := &togglev1.CreateToggleRequest{Toggle: &togglev1.Toggle{
		Key:         toggle.Key,
		Description: toggle.Description,
	}}

	_, err := c.breaker.Execute(func() (interface{}, error) {
		_, err := c.command.CreateToggle(ctx, req)
		if isServerError(err) {
			return nil, err
		}
		return nil, nil
	})
	return err
}

// Get gets a single toggle by its key.
func (c *Client) Get(ctx context.Context, key string) (*entity.Toggle, error) {
	req := &togglev1.GetToggleByKeyRequest{Key: key}

	tmp, err := c.breaker.Execute(func() (interface{}, error) {
		x, err := c.query.GetToggleByKey(ctx, req)
		if isServerError(err) {
			return nil, err
		}
		return x, nil
	})
	if err != nil {
		return nil, err
	}

	resp := tmp.(*togglev1.GetToggleByKeyResponse)
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

	_, err := c.breaker.Execute(func() (interface{}, error) {
		_, err := c.command.EnableToggle(ctx, req)
		if isServerError(err) {
			return nil, err
		}
		return nil, nil
	})
	return err
}

// Disable disables a toggle.
// It sets toggle's `is_enabled` attribute to be false.
func (c *Client) Disable(ctx context.Context, key string) error {
	req := &togglev1.DisableToggleRequest{Key: key}
	_, err := c.breaker.Execute(func() (interface{}, error) {
		_, err := c.command.DisableToggle(ctx, req)
		if isServerError(err) {
			return nil, err
		}
		return nil, nil
	})
	return err
}

// Delete deletes a toggle.
// It only deletes a nonactive toggle (is_enabled == false).
func (c *Client) Delete(ctx context.Context, key string) error {
	req := &togglev1.DeleteToggleRequest{Key: key}
	_, err := c.breaker.Execute(func() (interface{}, error) {
		_, err := c.command.DeleteToggle(ctx, req)
		if isServerError(err) {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func isServerError(err error) bool {
	switch status.Code(err) {
	case codes.ResourceExhausted, codes.Internal, codes.Unavailable:
		return true
	default:
		return false
	}
}

type noBreaker struct {
}

func (nb noBreaker) Execute(req func() (interface{}, error)) (interface{}, error) {
	return req()
}
