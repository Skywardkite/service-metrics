package repository

import (
	"context"
)

const querySetCounter = `
	INSERT INTO counters (name, value) 
	VALUES ($1, $2)
	ON CONFLICT (name) DO UPDATE SET value = counters.value + $2
`

func (r *PostgresStorage) SetCounter(ctx context.Context, name string, value int64) error {

	_, err := r.db.ExecContext(ctx, querySetCounter, name, value)

	return err
}