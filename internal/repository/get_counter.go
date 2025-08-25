package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const queryGetCounter = `SELECT value FROM counters WHERE name = $1`

var ErrCounterNotFound = errors.New("metric counter not found")

func (r *PostgresStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	classifier := NewPostgresErrorClassifier()
    delays := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

    var lastErr error
    for attempt := 0; attempt <= len(delays); attempt++ {
		var value int64
		err := r.db.GetContext(ctx, &value, queryGetCounter, name)
		if err == nil {
			return value, nil
		}

		classification := classifier.Classify(err)

        if classification == NonRetriable {
			if errors.Is(err, sql.ErrNoRows) {
				return 0, ErrCounterNotFound
			}
            return 0, err
        }
	}

	return 0, fmt.Errorf("операция прервана после %d попыток: %w", len(delays), lastErr)
}