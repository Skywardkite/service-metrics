package handler

import (
	"context"

	model "github.com/Skywardkite/service-metrics/internal/model"
)

type MockStorage struct{}

func (d *MockStorage) SetCounter(ctx context.Context, name string, value int64) error {
	return nil
}

func (d *MockStorage) SetGauge(ctx context.Context, name string, value float64) error {
	return nil
}

func (d *MockStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	return 0, nil
}

func (d *MockStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	return 0, nil
}

func (d *MockStorage) GetMetrics(ctx context.Context) (map[string]float64, map[string]int64, error) {
	return nil, nil, nil
}

func (d *MockStorage) SetMetricsBatch(ctx context.Context, metrics []model.Metrics) error {
	return nil
}

func (d *MockStorage) Ping() error {
	return nil
}

func (d *MockStorage) SaveMetrics(filePath string, gauges map[string]float64, counters map[string]int64) error {
	return nil
}
