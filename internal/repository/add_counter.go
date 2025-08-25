package repository

import (
	"context"
	"fmt"
	"time"
)

const querySetCounter = `
	INSERT INTO counters (name, value) 
	VALUES ($1, $2)
	ON CONFLICT (name) DO UPDATE SET value = counters.value + $2
`

func (r *PostgresStorage) SetCounter(ctx context.Context, name string, value int64) error {
	classifier := NewPostgresErrorClassifier()
    delays := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

    var lastErr error
    for attempt := 0; attempt <= len(delays); attempt++ {
		_, err := r.db.ExecContext(ctx, querySetCounter, name, value)
		if err == nil {
			return nil
		}

		// Определяем классификацию ошибки
        classification := classifier.Classify(err)

        if classification == NonRetriable {
            // Нет смысла повторять, возвращаем ошибку
            return err
        }
	}

	return fmt.Errorf("операция прервана после %d попыток: %w", len(delays), lastErr)
}