package service

import (
	"context"
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
	repo DisableToggleRepository
}

// NewToggleDisabler creates an instance of ToggleDisabler.
func NewToggleDisabler(repo DisableToggleRepository) *ToggleDisabler {
	return &ToggleDisabler{repo: repo}
}

// Disable disables a toggle.
// It just sets the toogle's is_enabled to be false.
// It doesn't validate toggle's key like ToggleCreator.Create does.
// But, it returns NotFound error if the toggle doesn't exist.
func (tc *ToggleDisabler) Disable(ctx context.Context, key string) error {
	return tc.repo.Disable(ctx, key, false)
}
