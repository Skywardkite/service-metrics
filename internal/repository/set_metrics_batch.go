package repository

import (
	"context"
	"time"

	model "github.com/Skywardkite/service-metrics/internal/model"
)

func (r *PostgresStorage) SetMetricsBatch(ctx context.Context, metrics []model.Metrics) error {
    classifier := NewPostgresErrorClassifier()
    delays := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

    var lastErr error
    for attempt := 0; attempt <= len(delays); attempt++ {
        err := r.trySetMetricsBatch(ctx, metrics)
        if err == nil {
            return nil
        }

        if classifier.Classify(err) == Retriable {
            lastErr = err
            if attempt < len(delays) {
                time.Sleep(delays[attempt])
                continue
            }
        }

        // не ретраим
        return err
    }

    return lastErr
}

func (r *PostgresStorage) trySetMetricsBatch(ctx context.Context, metrics []model.Metrics) error {
    tx, err := r.db.Beginx()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    for _, m := range metrics {
        switch m.MType {
        case "gauge":
            if m.Value != nil {
                _, err := tx.ExecContext(ctx, `
                    INSERT INTO gauges (name, value) 
                    VALUES ($1, $2)
                    ON CONFLICT (name) DO UPDATE SET value = $2
                `, m.ID, *m.Value)
                if err != nil {
                    return err
                }
            }
        case "counter":
            if m.Delta != nil {
                // Для счетчиков используем атомарное сложение
                _, err := tx.ExecContext(ctx, `
                    INSERT INTO counters (name, value) 
                    VALUES ($1, $2)
                    ON CONFLICT (name) DO UPDATE SET value = counters.value + $2
                `, m.ID, *m.Delta)
                if err != nil {
                    return err
                }
            }
        }
    }

    return tx.Commit()
}