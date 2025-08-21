package repository

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (r *PostgresStorage) GetMetrics(ctx context.Context) (map[string]float64, map[string]int64, error) {
	var gaugesEntity []Gauge
	var countersEntity []Counter

	err := r.db.SelectContext(ctx, &gaugesEntity, `SELECT * FROM gauges`)
	if err != nil && !errors.Is(err, sql.ErrNoRows){
		return nil, nil, err
	}
	err = r.db.SelectContext(ctx, &countersEntity, `SELECT * FROM counters`)
	if err != nil && !errors.Is(err, sql.ErrNoRows){
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