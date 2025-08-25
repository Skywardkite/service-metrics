package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const queryGetGauge = `SELECT value FROM gauges WHERE name = $1`

var ErrGaugeNotFound = errors.New("metric gauge not found")

func (r *PostgresStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	classifier := NewPostgresErrorClassifier()
    delays := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

    var lastErr error
    for attempt := 0; attempt <= len(delays); attempt++ {
		var value float64
		err := r.db.GetContext(ctx, &value, queryGetGauge, name)
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