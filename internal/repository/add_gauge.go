package repository

import (
	"context"
	"fmt"
	"time"
)

const querySetGauge = `
	INSERT INTO gauges (name, value) 
    VALUES ($1, $2)
    ON CONFLICT (name) DO UPDATE SET value = $2
`

func (r *PostgresStorage) SetGauge(ctx context.Context, name string, value float64) error {
	classifier := NewPostgresErrorClassifier()
    delays := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

    var lastErr error
    for attempt := 0; attempt <= len(delays); attempt++ {
		_, err := r.db.ExecContext(ctx, querySetGauge, name, value)
		if err == nil {
			return nil
		}

        classification := classifier.Classify(err)

        if classification == NonRetriable {
            return err
        }
	}

	return fmt.Errorf("операция прервана после %d попыток: %w", len(delays), lastErr)
}