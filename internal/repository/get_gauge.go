package repository

import (
	"context"
	"database/sql"
	"errors"
)

const queryGetGauger = `SELECT value FROM gauges WHERE name = $1`

var ErrGaugeNotFound = errors.New("metric gauge not found")

func (r *PostgresStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	var value float64
	err := r.db.GetContext(ctx, &value, queryGetGauger, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrGaugeNotFound
		}

		return 0, err
	}

	return value, nil
}
