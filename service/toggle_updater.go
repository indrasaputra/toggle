package service

import (
	"context"
)

// UpdateToggle defines the interface to update a toggle.
type UpdateToggle interface {
	// Enable enables the toggle in the system.
	// It sets the is_enabled to be true.
	// It must ensure that the toggle exists.
	// Otherwise, it returns error.
	Enable(ctx context.Context, key string) error
	// Disable disables the toggle in the system.
	// It sets the is_enabled to be false.
	// It must ensure that the toggle exists.
	// Otherwise, it returns error.
	Disable(ctx context.Context, key string) error
}

// UpdateToggleRepository defines the interface to save a toggle into the repository.
type UpdateToggleRepository interface {
	// UpdateIsEnabled updates the toggle's is_enabled value in the repository.
	// It returns NotFound error if the toggle doesn't exist.
	UpdateIsEnabled(ctx context.Context, key string, value bool) error
}

// ToggleUpdater is responsible for updating a toggle.
type ToggleUpdater struct {
	repo UpdateToggleRepository
}

// NewToggleUpdater creates an instance of ToggleUpdater.
func NewToggleUpdater(repo UpdateToggleRepository) *ToggleUpdater {
	return &ToggleUpdater{repo: repo}
}

// Enable enables a toggle.
// It just sets the toogle's is_enabled to be true.
// It doesn't validate toggle's key like ToggleCreator.Create does.
// But, it returns NotFound error if the toggle doesn't exist.
func (tc *ToggleUpdater) Enable(ctx context.Context, key string) error {
	return tc.repo.UpdateIsEnabled(ctx, key, true)
}

// Disable disables a toggle.
// It just sets the toogle's is_enabled to be false.
// It doesn't validate toggle's key like ToggleCreator.Create does.
// But, it returns NotFound error if the toggle doesn't exist.
func (tc *ToggleUpdater) Disable(ctx context.Context, key string) error {
	return tc.repo.UpdateIsEnabled(ctx, key, false)
}
