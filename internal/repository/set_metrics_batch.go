package repository

import (
	"context"

	model "github.com/Skywardkite/service-metrics/internal/model"
)

func (r *PostgresStorage) SetMetricsBatch(ctx context.Context, metrics []model.Metrics) error {
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