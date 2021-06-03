package postgres

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

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
func (t *Toggle) Insert(ctx context.Context, toggle *entity.Toggle) error {
	if toggle == nil {
		return entity.ErrEmptyToggle()
	}

	query := "INSERT INTO " +
		"toggles (key, is_enabled, description, created_at, updated_at) " +
		"VALUES ($1, $2, $3, $4, $5)"

	_, err := t.pool.Exec(ctx, query,
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

// GetByKey gets a toggle from database.
// It returns entity.ErrNotFound if toggle can't be found.
func (t *Toggle) GetByKey(ctx context.Context, key string) (*entity.Toggle, error) {
	query := "SELECT key, is_enabled, description, created_at, updated_at FROM toggles WHERE key = $1 LIMIT 1"
	row := t.pool.QueryRow(ctx, query, key)

	res := entity.Toggle{}
	err := row.Scan(&res.Key, &res.IsEnabled, &res.Description, &res.CreatedAt, &res.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, entity.ErrNotFound()
	}
	if err != nil {
		return nil, entity.ErrInternal(err.Error())
	}
	return &res, nil
}

// GetAll gets all available toggles from storage.
// If there isn't any toggle in repository, it returns empty list of toggle and nil error.
func (t *Toggle) GetAll(ctx context.Context, limit uint) ([]*entity.Toggle, error) {
	query := "SELECT key, is_enabled, description, created_at, updated_at FROM toggles LIMIT $1"
	rows, err := t.pool.Query(ctx, query, limit)
	if err != nil {
		return []*entity.Toggle{}, entity.ErrInternal(err.Error())
	}
	defer rows.Close()

	res := []*entity.Toggle{}
	for rows.Next() {
		var tmp entity.Toggle
		if err := rows.Scan(&tmp.Key, &tmp.IsEnabled, &tmp.Description, &tmp.CreatedAt, &tmp.UpdatedAt); err != nil {
			log.Printf("[Toggle-GetAll] scan rows error: %s", err.Error())
			continue
		}
		res = append(res, &tmp)
	}
	if rows.Err() != nil {
		return []*entity.Toggle{}, entity.ErrInternal(rows.Err().Error())
	}
	return res, nil
}

// UpdateIsEnabled updates the toggle's is_enabled value in the storage.
// It should handle if the toggle doesn't exist.
func (t *Toggle) UpdateIsEnabled(ctx context.Context, key string, value bool) error {
	if err := t.checkIfToggleExists(ctx, key); err != nil {
		return err
	}

	query := "UPDATE toggles SET is_enabled = $1, updated_at = $2 WHERE key = $3"
	_, err := t.pool.Exec(ctx, query, value, time.Now(), key)
	if err != nil {
		return entity.ErrInternal(err.Error())
	}
	return nil
}

// Delete deletes a toggle from PostgreSQL.
// If the group doesn't exist, it doesn't returns error.
func (t *Toggle) Delete(ctx context.Context, key string) error {
	query := "DELETE FROM toggles WHERE key = $1"
	_, err := t.pool.Exec(ctx, query, key)
	if err != nil {
		return entity.ErrInternal(err.Error())
	}
	return nil
}

func (t *Toggle) checkIfToggleExists(ctx context.Context, key string) error {
	query := "SELECT EXISTS(SELECT 1 FROM toggles WHERE key = $1)"
	row := t.pool.QueryRow(ctx, query, key)

	var found bool
	if err := row.Scan(&found); err != nil {
		return entity.ErrInternal(err.Error())
	}
	if !found {
		return entity.ErrNotFound()
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
