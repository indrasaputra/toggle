package service

import (
	"context"
	"log"

	"github.com/indrasaputra/toggle/entity"
)

// DeleteToggle defines the interface to delete a toggle.
type DeleteToggle interface {
	// DeleteByKey deletes a single toggle by its key.
	DeleteByKey(ctx context.Context, key string) error
}

// DeleteToggleRepository defines the interface to delete toggle from the repository.
type DeleteToggleRepository interface {
	// GetByKey gets a single toggle from the repository.
	// If the toggle can't be found, it returns NotFound error.
	GetByKey(ctx context.Context, key string) (*entity.Toggle, error)
	// DeleteByKey deletes a single toggle from the repository.
	// If the toggle can't be found, it doesn't return error.
	DeleteByKey(ctx context.Context, key string) error
}

// ToggleDeleter is responsible for deleting a toggle.
type ToggleDeleter struct {
	repo      DeleteToggleRepository
	publisher TogglePublisher
}

// NewToggleDeleter creates an instance of ToggleDeleter.
func NewToggleDeleter(repo DeleteToggleRepository, publisher TogglePublisher) *ToggleDeleter {
	return &ToggleDeleter{
		repo:      repo,
		publisher: publisher,
	}
}

// DeleteByKey deletes a toggle by its key.
// It only deletes disabled toggle.
func (td *ToggleDeleter) DeleteByKey(ctx context.Context, key string) error {
	toggle, err := td.repo.GetByKey(ctx, key)
	if err != nil {
		return err
	}
	if toggle.IsEnabled {
		return entity.ErrProhibitedToDelete()
	}
	if err := td.repo.DeleteByKey(ctx, key); err != nil {
		return err
	}
	if err := td.publisher.Publish(ctx, entity.EventToggleDeleted(toggle)); err != nil {
		log.Printf("publish on toggle deleter error: %v", err)
	}
	return nil
}
