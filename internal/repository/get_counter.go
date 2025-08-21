package repository

import (
	"context"
	"database/sql"
	"errors"
)

var ErrCounterNotFound = errors.New("metric counter not found")

func (r *PostgresStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	query := `SELECT value FROM counters WHERE name = $1`

	var value int64
	err := r.db.GetContext(ctx, &value, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrCounterNotFound
		}

		return 0, err
	}

	return value, nil
}