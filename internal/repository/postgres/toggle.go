package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgconn"

	"github.com/indrasaputra/toggle/entity"
)

const (
	// errCodeUniqueViolation is derived from https://www.postgresql.org/docs/11/errcodes-appendix.html
	errCodeUniqueViolation = "23505"
)

// Toggle is responsible to connect toggle entity with toggles table in PostgreSQL.
type Toggle struct {
	pool PgxPoolIface
}

// NewToggle creates an instance of Toggle.
func NewToggle(pool PgxPoolIface) *Toggle {
	return &Toggle{pool: pool}
}

// Insert inserts the toggle into the toggles table.
func (g *Toggle) Insert(ctx context.Context, toggle *entity.Toggle) error {
	if toggle == nil {
		return entity.ErrEmptyToggle()
	}

	query := "INSERT INTO " +
		"toggles (key, is_enabled, description, created_at, updated_at) " +
		"VALUES ($1, $2, $3, $4, $5)"

	_, err := g.pool.Exec(ctx, query,
		toggle.Key,
		toggle.IsEnabled,
		toggle.Description,
		time.Now(),
		time.Now(),
	)

	if err != nil && isUniqueViolationErr(err) {
		return entity.ErrAlreadyExists()
	}
	if err != nil {
		return entity.ErrInternal(err.Error())
	}
	return nil
}

func isUniqueViolationErr(err error) bool {
	pgerr, ok := err.(*pgconn.PgError)
	if !ok {
		return false
	}
	return pgerr.Code == errCodeUniqueViolation
}
