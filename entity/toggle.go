package entity

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
)

// Toggle defines the logical data of a toggle.
type Toggle struct {
	// Key defines the toggle's identifier.
	// It must be unique from the rest.
	Key string
	// IsEnabled defines the toggle's value.
	IsEnabled bool
	// Description defines toggle's description.
	Description string
	// CreatedAt defines the time when the toggle was created.
	CreatedAt time.Time
	// Updated defines the time when the toggle was last updated.
	UpdatedAt time.Time
}

// EventToggleCreated creates an event for created toggle.
func EventToggleCreated(toggle *Toggle) *togglev1.ToggleEvent {
	return &togglev1.ToggleEvent{
		Name:      togglev1.ToggleEventName_TOGGLE_EVENT_NAME_CREATED,
		Toggle:    createAPIToggle(toggle),
		CreatedAt: timestamppb.Now(),
	}
}

// EventToggleEnabled creates an event for enabled toggle.
func EventToggleEnabled(toggle *Toggle) *togglev1.ToggleEvent {
	return &togglev1.ToggleEvent{
		Name:      togglev1.ToggleEventName_TOGGLE_EVENT_NAME_ENABLED,
		Toggle:    createAPIToggle(toggle),
		CreatedAt: timestamppb.Now(),
	}
}

// EventToggleDisabled creates an event for disabled toggle.
func EventToggleDisabled(toggle *Toggle) *togglev1.ToggleEvent {
	return &togglev1.ToggleEvent{
		Name:      togglev1.ToggleEventName_TOGGLE_EVENT_NAME_DISABLED,
		Toggle:    createAPIToggle(toggle),
		CreatedAt: timestamppb.Now(),
	}
}

// EventToggleDeleted creates an event for deleted toggle.
func EventToggleDeleted(toggle *Toggle) *togglev1.ToggleEvent {
	return &togglev1.ToggleEvent{
		Name:      togglev1.ToggleEventName_TOGGLE_EVENT_NAME_DELETED,
		Toggle:    createAPIToggle(toggle),
		CreatedAt: timestamppb.Now(),
	}
}

func createAPIToggle(toggle *Toggle) *togglev1.Toggle {
	return &togglev1.Toggle{
		Key:         toggle.Key,
		IsEnabled:   toggle.IsEnabled,
		Description: toggle.Description,
	}
}
