package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/entity"
)

func TestErrInternal(t *testing.T) {
	t.Run("success get internal error", func(t *testing.T) {
		err := entity.ErrInternal("")

		assert.Contains(t, err.Error(), "rpc error: code = Internal")
	})
}

func TestErrEmptyToggle(t *testing.T) {
	t.Run("success get empty toggle error", func(t *testing.T) {
		err := entity.ErrEmptyToggle()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrAlreadyExists(t *testing.T) {
	t.Run("success get already exists error", func(t *testing.T) {
		err := entity.ErrAlreadyExists()

		assert.Contains(t, err.Error(), "rpc error: code = AlreadyExists")
	})
}

func TestErrInvalidKey(t *testing.T) {
	t.Run("success get invalid key error", func(t *testing.T) {
		err := entity.ErrInvalidKey()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrInvalidValue(t *testing.T) {
	t.Run("success get invalid value error", func(t *testing.T) {
		err := entity.ErrInvalidValue()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrNotFound(t *testing.T) {
	t.Run("success get not found error", func(t *testing.T) {
		err := entity.ErrNotFound()

		assert.Contains(t, err.Error(), "rpc error: code = NotFound")
	})
}

func TestErrProhibitedToDelete(t *testing.T) {
	t.Run("success get prohibited to delete error", func(t *testing.T) {
		err := entity.ErrProhibitedToDelete()

		assert.Contains(t, err.Error(), "rpc error: code = FailedPrecondition")
	})
}
