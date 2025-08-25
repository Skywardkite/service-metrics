package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

func (r *PostgresStorage) GetMetrics(ctx context.Context) (map[string]float64, map[string]int64, error) {
	classifier := NewPostgresErrorClassifier()
	delays := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

	var lastErr error

	for attempt := 0; attempt <= len(delays); attempt++ {
		gauges, counters, err := r.tryGetMetrics(ctx)
		if err == nil {
			return gauges, counters, nil
		}

		if classifier.Classify(err) == Retriable {
			lastErr = err
			if attempt < len(delays) {
				time.Sleep(delays[attempt])
				continue
			}
		}

		return nil, nil, err
	}

	return nil, nil, lastErr
}

func (r *PostgresStorage) tryGetMetrics(ctx context.Context) (map[string]float64, map[string]int64, error) {
	var gaugesEntity []Gauge
	var countersEntity []Counter

	err := r.db.SelectContext(ctx, &gaugesEntity, `SELECT * FROM gauges`)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, err
	}

	err = r.db.SelectContext(ctx, &countersEntity, `SELECT * FROM counters`)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, err
	}

	gauges := make(map[string]float64, len(gaugesEntity))
	counters := make(map[string]int64, len(countersEntity))

	for _, gauge := range gaugesEntity {
		gauges[gauge.Name] = gauge.Value
	}

	for _, counter := range countersEntity {
		counters[counter.Name] = counter.Value
	}

	return gauges, counters, nil
}