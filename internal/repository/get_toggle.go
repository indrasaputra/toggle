package repository

import (
	"context"

	"github.com/indrasaputra/toggle/entity"
)

// GetToggleDatabase defines the interface to get toggle from database.
type GetToggleDatabase interface {
	// GetByKey gets a toggle from database.
	// It must return codes.NotFound from package package google.golang.org/grpc/codes if data can't be found.
	GetByKey(ctx context.Context, key string) (*entity.Toggle, error)
	// GetAll gets all available toggles from database.
	// If there isn't any toggle in repository, it returns empty list of toggle and nil error.
	GetAll(ctx context.Context) ([]*entity.Toggle, error)
}

// GetToggleCache defines the interface to get a toggle in cache.
type GetToggleCache interface {
	// Get gets a toggle in cache.
	// It only returns error of there is error in the system.
	// If the data can't be found but the system is fine, it returns nil.
	Get(ctx context.Context, key string) (*entity.Toggle, error)
	// Set sets a toggle in cache.
	Set(ctx context.Context, toggle *entity.Toggle) error
}

// ToggleGetter is responsible to get the toggle from storage.
// It uses database and cache.
type ToggleGetter struct {
	database GetToggleDatabase
	cache    GetToggleCache
}

// NewToggleGetter creates an instance of ToggleGetter.
func NewToggleGetter(database GetToggleDatabase, cache GetToggleCache) *ToggleGetter {
	return &ToggleGetter{database: database, cache: cache}
}

// GetByKey gets the toggle from the storage.
// First, it accessess the cache. If success, the data will be returned instantly..
// Otherwise, it checks the data in database.
func (tg *ToggleGetter) GetByKey(ctx context.Context, key string) (*entity.Toggle, error) {
	toggle, err := tg.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if toggle != nil {
		return toggle, nil
	}

	toggle, err = tg.database.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	_ = tg.cache.Set(ctx, toggle)
	return toggle, nil
}

// GetAll gets all available toggles from storage.
// If there isn't any toggle in repository, it returns empty list of toggle and nil error.
func (tg *ToggleGetter) GetAll(ctx context.Context) ([]*entity.Toggle, error) {
	return tg.database.GetAll(ctx)
}
