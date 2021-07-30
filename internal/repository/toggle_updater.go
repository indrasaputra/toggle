package repository

import (
	"context"
)

// UpdateToggleDatabase defines the interface to update a toggle in database.
type UpdateToggleDatabase interface {
	// UpdateIsEnabled updates the toggle's is_enabled value in the repository.
	// It should handle if the toggle doesn't exist.
	UpdateIsEnabled(ctx context.Context, key string, value bool) error
}

// UpdateToggleCache defines the interface to set (there is no update in cache) a toggle in cache.
type UpdateToggleCache interface {
	// SetIsEnabled sets is_enabled field in cache.
	SetIsEnabled(ctx context.Context, key string, value bool) error
}

// ToggleUpdater is responsible to update the toggle in storage.
// It uses database and cache.
type ToggleUpdater struct {
	database UpdateToggleDatabase
	cache    UpdateToggleCache
}

// NewToggleUpdater creates an instance of ToggleUpdater.
func NewToggleUpdater(database UpdateToggleDatabase, cache UpdateToggleCache) *ToggleUpdater {
	return &ToggleUpdater{database: database, cache: cache}
}

// Enable updates the toggle's is_enabled value to be true in the storage.
// First, it updates the data in database. If success, the data will be set to cache.
// It ignores the error from cache since it can always be generated when retrieving the data.
// But, it doesn't ignore the error from the database.
func (ti *ToggleUpdater) Enable(ctx context.Context, key string, value bool) error {
	return ti.updateIsEnabled(ctx, key, value)
}

// Disable updates the toggle's is_enabled value to be false in the storage.
// First, it updates the data in database. If success, the data will be set to cache.
// It ignores the error from cache since it can always be generated when retrieving the data.
// But, it doesn't ignore the error from the database.
func (ti *ToggleUpdater) Disable(ctx context.Context, key string, value bool) error {
	return ti.updateIsEnabled(ctx, key, value)
}

func (ti *ToggleUpdater) updateIsEnabled(ctx context.Context, key string, value bool) error {
	if err := ti.database.UpdateIsEnabled(ctx, key, value); err != nil {
		return err
	}
	_ = ti.cache.SetIsEnabled(ctx, key, value)
	return nil
}
