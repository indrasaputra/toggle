package service

import (
	"context"
	"log"

	"github.com/indrasaputra/toggle/entity"
)

// EnableToggle defines the interface to enable a toggle.
type EnableToggle interface {
	// Enable enables the toggle in the system.
	// It sets the is_enabled to be true.
	// It must ensure that the toggle exists.
	// Otherwise, it returns error.
	Enable(ctx context.Context, key string) error
}

// EnableToggleRepository defines the interface to save an enabled toggle into the repository.
type EnableToggleRepository interface {
	// Enable updates the toggle's is_enabled value in the repository.
	// It returns NotFound error if the toggle doesn't exist.
	Enable(ctx context.Context, key string, value bool) error
}

// ToggleEnabler is responsible for enabling a toggle.
type ToggleEnabler struct {
	repo      EnableToggleRepository
	publisher TogglePublisher
}

// NewToggleEnabler creates an instance of ToggleEnabler.
func NewToggleEnabler(repo EnableToggleRepository, publisher TogglePublisher) *ToggleEnabler {
	return &ToggleEnabler{
		repo:      repo,
		publisher: publisher,
	}
}

// Enable enables a toggle.
// It just sets the toogle's is_enabled to be true.
// It doesn't validate toggle's key like ToggleCreator.Create does.
// But, it returns NotFound error if the toggle doesn't exist.
func (te *ToggleEnabler) Enable(ctx context.Context, key string) error {
	if err := te.repo.Enable(ctx, key, true); err != nil {
		return err
	}
	if err := te.publisher.Publish(ctx, entity.EventToggleEnabled(&entity.Toggle{Key: key})); err != nil {
		log.Printf("publish on toggle enabler error: %v", err)
	}
	return nil
}
