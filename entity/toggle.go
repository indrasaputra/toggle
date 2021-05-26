package entity

import "time"

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
