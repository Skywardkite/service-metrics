package repository

import "context"

func (r *PostgresStorage) SetCounter(ctx context.Context, name string, value int64) error {
	query := `INSERT INTO counters (name, value) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, name, value)

	return err
}