package repository

import (
	"context"

	"github.com/indrasaputra/toggle/entity"
)

// InsertToggleDatabase defines the interface to insert a toggle into the database.
type InsertToggleDatabase interface {
	// Insert inserts a new toggle to the database.
	// It should handle if the toggle already exists.
	Insert(ctx context.Context, toggle *entity.Toggle) error
}

// SetToggleCache defines the interface to set a toggle in cache.
type SetToggleCache interface {
	// Set sets a toggle in cache.
	Set(ctx context.Context, toggle *entity.Toggle) error
}

// ToggleInserter is responsible to insert the toggle into storage.
// It uses database and cache.
type ToggleInserter struct {
	database InsertToggleDatabase
	cache    SetToggleCache
}

// NewToggleInserter creates an instance of ToggleInserter.
func NewToggleInserter(database InsertToggleDatabase, cache SetToggleCache) *ToggleInserter {
	return &ToggleInserter{database: database, cache: cache}
}

// Insert inserts the toggle into the storage.
// First, it inserts to database. If success, the data will be set to cache.
// It ignores the error from cache since it can always be generated when retrieving the data.
// But, it doesn't ignore the error from the database.
func (ti *ToggleInserter) Insert(ctx context.Context, toggle *entity.Toggle) error {
	if toggle == nil {
		return entity.ErrEmptyToggle()
	}

	if err := ti.database.Insert(ctx, toggle); err != nil {
		return err
	}
	_ = ti.cache.Set(ctx, toggle)
	return nil
}
