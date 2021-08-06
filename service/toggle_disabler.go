package service

import (
	"context"
	"log"

	"github.com/indrasaputra/toggle/entity"
)

// DisableToggle defines the interface to disable a toggle.
type DisableToggle interface {
	// Disable disables the toggle in the system.
	// It sets the is_enabled to be false.
	// It must ensure that the toggle exists.
	// Otherwise, it returns error.
	Disable(ctx context.Context, key string) error
}

// DisableToggleRepository defines the interface to save a disabled toggle into the repository.
type DisableToggleRepository interface {
	// Disable updates the toggle's is_enabled value in the repository.
	// It returns NotFound error if the toggle doesn't exist.
	Disable(ctx context.Context, key string, value bool) error
}

// ToggleDisabler is responsible for disabling a toggle.
type ToggleDisabler struct {
	repo      DisableToggleRepository
	publisher TogglePublisher
}

// NewToggleDisabler creates an instance of ToggleDisabler.
func NewToggleDisabler(repo DisableToggleRepository, publisher TogglePublisher) *ToggleDisabler {
	return &ToggleDisabler{
		repo:      repo,
		publisher: publisher,
	}
}

// Disable disables a toggle.
// It just sets the toogle's is_enabled to be false.
// It doesn't validate toggle's key like ToggleCreator.Create does.
// But, it returns NotFound error if the toggle doesn't exist.
func (td *ToggleDisabler) Disable(ctx context.Context, key string) error {
	if err := td.repo.Disable(ctx, key, false); err != nil {
		return err
	}
	if err := td.publisher.Publish(ctx, entity.EventToggleDisabled(&entity.Toggle{Key: key})); err != nil {
		log.Printf("publish on toggle deleter error: %v", err)
	}
	return nil
}
