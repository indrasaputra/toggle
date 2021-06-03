package service

import (
	"context"

	"github.com/indrasaputra/toggle/entity"
)

// GetToggle defines the interface to get toggle.
type GetToggle interface {
	// GetByKey gets a single toggle by its key.
	GetByKey(ctx context.Context, key string) (*entity.Toggle, error)
	// GetAll gets all toggles available in system.
	GetAll(ctx context.Context) ([]*entity.Toggle, error)
}

// GetToggleRepository defines the interface to get toggle from the repository.
type GetToggleRepository interface {
	// GetByKey gets a single toggle from the repository.
	// If the toggle can't be found, it returns NotFound error.
	GetByKey(ctx context.Context, key string) (*entity.Toggle, error)
	// GetAll gets all toggles available in repository.
	// If there isn't any toggle in repository, it returns empty list of toggle and nil error.
	GetAll(ctx context.Context) ([]*entity.Toggle, error)
}

// ToggleGetter is responsible for getting toggle.
type ToggleGetter struct {
	repo GetToggleRepository
}

// NewToggleGetter creates an instance of ToggleGetter.
func NewToggleGetter(repo GetToggleRepository) *ToggleGetter {
	return &ToggleGetter{repo: repo}
}

// GetByKey gets a single toggle.
// If the toggle can't be found, it returns NotFound error.
func (tg *ToggleGetter) GetByKey(ctx context.Context, key string) (*entity.Toggle, error) {
	return tg.repo.GetByKey(ctx, key)
}

// GetAll gets all toggles available in repository.
// If there isn't any toggle in repository, it returns empty list of toggle and nil error.
func (tg *ToggleGetter) GetAll(ctx context.Context) ([]*entity.Toggle, error) {
	return tg.repo.GetAll(ctx)
}
