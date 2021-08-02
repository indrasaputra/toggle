package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/entity"
)

func TestEventToggleCreated(t *testing.T) {
	t.Run("successfully create event toggle created", func(t *testing.T) {
		event := entity.EventToggleCreated(&entity.Toggle{})
		assert.NotNil(t, event)
	})
}

func TestEventToggleEnabled(t *testing.T) {
	t.Run("successfully create event toggle enabled", func(t *testing.T) {
		event := entity.EventToggleEnabled(&entity.Toggle{})
		assert.NotNil(t, event)
	})
}

func TestEventToggleDisabled(t *testing.T) {
	t.Run("successfully create event toggle disabled", func(t *testing.T) {
		event := entity.EventToggleDisabled(&entity.Toggle{})
		assert.NotNil(t, event)
	})
}

func TestEventToggleDeleted(t *testing.T) {
	t.Run("successfully create event toggle deleted", func(t *testing.T) {
		event := entity.EventToggleDeleted(&entity.Toggle{})
		assert.NotNil(t, event)
	})
}
