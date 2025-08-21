package repository

import "context"

func (r *PostgresStorage) SetGauge(ctx context.Context, name string, value float64) error {
	query := `INSERT INTO gauges (name, value) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, name, value)

	return err
}