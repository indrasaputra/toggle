package service

import (
	"context"
	"regexp"
	"strings"

	"github.com/indrasaputra/toggle/entity"
)

var (
	regexCompiler *regexp.Regexp
)

func init() {
	regexCompiler = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)
}

// CreateToggle defines the interface to create a toggle.
type CreateToggle interface {
	// Create creates a new toggle and store it in the storage.
	// It must check the uniqueness of the toggle's key.
	Create(ctx context.Context, toggle *entity.Toggle) error
}

// CreateToggleRepository defines the interface to save a toggle into the repository.
type CreateToggleRepository interface {
	// Insert inserts the toggle into the repository.
	// It also validates if the toggle's key is unique.
	Insert(ctx context.Context, toggle *entity.Toggle) error
}

// ToggleCreator is responsible for creating a new toggle.
type ToggleCreator struct {
	repo CreateToggleRepository
}

// NewToggleCreator creates an instance of ToggleCreator.
func NewToggleCreator(repo CreateToggleRepository) *ToggleCreator {
	return &ToggleCreator{repo: repo}
}

// Create creates a new toggle.
// It will reject if the toggle's key already exists in the system.
// The toggle's key is converted to lower case.
func (tc *ToggleCreator) Create(ctx context.Context, toggle *entity.Toggle) error {
	if err := validateToggle(toggle); err != nil {
		return err
	}
	sanitizeToggle(toggle)

	return tc.repo.Insert(ctx, toggle)
}

func sanitizeToggle(toggle *entity.Toggle) {
	toggle.Key = strings.TrimSpace(toggle.Key)
	toggle.Key = strings.ToLower(toggle.Key)
	toggle.Description = strings.TrimSpace(toggle.Description)
	toggle.IsEnabled = false
}

func validateToggle(toggle *entity.Toggle) error {
	if toggle == nil {
		return entity.ErrEmptyToggle()
	}

	if !regexCompiler.MatchString(toggle.Key) {
		return entity.ErrInvalidKey()
	}
	return nil
}
