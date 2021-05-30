package postgres_test

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/toggle/entity"
	"github.com/indrasaputra/toggle/internal/repository/postgres"
)

var (
	testCtx                = context.Background()
	testToggleKey          = "toggle-1"
	testToggleDesc         = "description"
	testToggle             = &entity.Toggle{Key: testToggleKey, Description: testToggleDesc}
	errPostgresInternalMsg = "database down"
	errPostgresInternal    = errors.New(errPostgresInternalMsg)
)

type ToggleExecutor struct {
	toggle *postgres.Toggle
	pgx    pgxmock.PgxPoolIface
}

func TestNewToggle(t *testing.T) {
	t.Run("successfully create an instance of Toggle", func(t *testing.T) {
		exec := createToggleExecutor()
		assert.NotNil(t, exec.toggle)
	})
}

func TestToggle_Insert(t *testing.T) {
	t.Run("nil toggle is prohibited", func(t *testing.T) {
		exec := createToggleExecutor()

		err := exec.toggle.Insert(testCtx, nil)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyToggle(), err)
	})

	t.Run("postgres database returns internal error", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.pgx.
			ExpectExec(`INSERT INTO toggles \(key, is_enabled, description, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`).
			WillReturnError(errPostgresInternal)

		err := exec.toggle.Insert(testCtx, testToggle)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(errPostgresInternalMsg), err)
	})

	t.Run("insert duplicate toggle", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.pgx.
			ExpectExec(`INSERT INTO toggles \(key, is_enabled, description, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`).
			WillReturnError(&pgconn.PgError{Code: "23505"})

		err := exec.toggle.Insert(testCtx, testToggle)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrAlreadyExists(), err)
	})

	t.Run("success insert a new toggle", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.pgx.
			ExpectExec(`INSERT INTO toggles \(key, is_enabled, description, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := exec.toggle.Insert(testCtx, testToggle)

		assert.Nil(t, err)
	})
}

func TestToggle_GetByKey(t *testing.T) {
	t.Run("select by key query returns empty row", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.pgx.
			ExpectQuery(`SELECT key, is_enabled, description, created_at, updated_at FROM toggles WHERE key = \$1 LIMIT 1`).
			WillReturnError(pgx.ErrNoRows)

		res, err := exec.toggle.GetByKey(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
		assert.Nil(t, res)
	})

	t.Run("select by key query returns error", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.pgx.
			ExpectQuery(`SELECT key, is_enabled, description, created_at, updated_at FROM toggles WHERE key = \$1 LIMIT 1`).
			WillReturnError(errPostgresInternal)

		res, err := exec.toggle.GetByKey(testCtx, testToggleKey)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternal(errPostgresInternalMsg), err)
		assert.Nil(t, res)
	})

	t.Run("successfully retrieve row", func(t *testing.T) {
		exec := createToggleExecutor()
		exec.pgx.
			ExpectQuery(`SELECT key, is_enabled, description, created_at, updated_at FROM toggles WHERE key = \$1 LIMIT 1`).
			WillReturnRows(pgxmock.
				NewRows([]string{"key", "is_enabled", "description", "created_at", "updated_at"}).
				AddRow(testToggleKey, true, testToggleDesc, time.Now(), time.Now()),
			)

		res, err := exec.toggle.GetByKey(testCtx, testToggleKey)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
}

func createToggleExecutor() *ToggleExecutor {
	mock, err := pgxmock.NewPool(pgxmock.MonitorPingsOption(true))
	if err != nil {
		log.Panicf("error opening a stub database connection: %v\n", err)
	}

	toggle := postgres.NewToggle(mock)
	return &ToggleExecutor{
		toggle: toggle,
		pgx:    mock,
	}
}
