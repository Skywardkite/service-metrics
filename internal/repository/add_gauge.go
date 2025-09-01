package repository

import (
	"context"
)

const querySetGauge = `
	INSERT INTO gauges (name, value) 
    VALUES ($1, $2)
    ON CONFLICT (name) DO UPDATE SET value = $2
`

func (r *PostgresStorage) SetGauge(ctx context.Context, name string, value float64) error {
	
	_, err := r.db.ExecContext(ctx, querySetGauge, name, value)

	return err
}