package repository

import (
	"context"

	"github.com/indrasaputra/toggle/entity"
)

// DeleteToggleDatabase defines the interface to delete toggle from database.
type DeleteToggleDatabase interface {
	// GetByKey gets a toggle from database.
	// It must return codes.NotFound from package package google.golang.org/grpc/codes if data can't be found.
	GetByKey(ctx context.Context, key string) (*entity.Toggle, error)
	// Delete deletes a toggle from database.
	// It doesn't return any error if toggle is not found.
	Delete(ctx context.Context, key string) error
}

// DeleteToggleCache defines the interface to delete a toggle in cache.
type DeleteToggleCache interface {
	// Delete deletes a toggle from database.
	// It doesn't return any error if toggle is not found.
	Delete(ctx context.Context, key string) error
}

// ToggleDeleter is responsible to delete the toggle from storage.
// It uses database and cache.
type ToggleDeleter struct {
	database DeleteToggleDatabase
	cache    DeleteToggleCache
}

// NewToggleDeleter creates an instance of ToggleDeleter.
func NewToggleDeleter(database DeleteToggleDatabase, cache DeleteToggleCache) *ToggleDeleter {
	return &ToggleDeleter{database: database, cache: cache}
}

// GetByKey gets the toggle from the storage.
// It accessess the database directly without checking the cache.
func (td *ToggleDeleter) GetByKey(ctx context.Context, key string) (*entity.Toggle, error) {
	toggle, err := td.database.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return toggle, nil
}

// DeleteByKey deletes the toggle from the storage.
// It doesn't return any error if toggle is not found.
func (td *ToggleDeleter) DeleteByKey(ctx context.Context, key string) error {
	if err := td.cache.Delete(ctx, key); err != nil {
		return err
	}
	return td.database.Delete(ctx, key)
}
