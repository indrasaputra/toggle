package service

import (
	"context"
	"log"
	"regexp"
	"strings"

	"github.com/indrasaputra/toggle/entity"
	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
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

// TogglePublisher defines the interface to publish toggle to message queue.
type TogglePublisher interface {
	// Publish publishes event toggle to message queue.
	Publish(ctx context.Context, event *togglev1.EventToggle) error
}

// ToggleCreator is responsible for creating a new toggle.
type ToggleCreator struct {
	repo      CreateToggleRepository
	publisher TogglePublisher
}

// NewToggleCreator creates an instance of ToggleCreator.
func NewToggleCreator(repo CreateToggleRepository, publisher TogglePublisher) *ToggleCreator {
	return &ToggleCreator{
		repo:      repo,
		publisher: publisher,
	}
}

// Create creates a new toggle.
// It will reject if the toggle's key already exists in the system.
// The toggle's key is converted to lower case.
func (tc *ToggleCreator) Create(ctx context.Context, toggle *entity.Toggle) error {
	if err := validateToggle(toggle); err != nil {
		return err
	}
	sanitizeToggle(toggle)

	if err := tc.repo.Insert(ctx, toggle); err != nil {
		return err
	}
	if err := tc.publisher.Publish(ctx, entity.EventToggleCreated(toggle)); err != nil {
		log.Printf("publish on toggle creator error: %v", err)
	}
	return nil
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
