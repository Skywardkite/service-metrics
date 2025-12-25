package handler

import (
	"context"

	model "github.com/Skywardkite/service-metrics/internal/model"
)

type MockService struct{}

func (m *MockService) UpdateMetric(ctx context.Context, metricType, metricName, metricValue string) error {
	return nil
}

func (s *MockService) GetMetric(ctx context.Context, metricType, metricName string) (string, error) {
	return "", nil
}

func (s *MockService) GetAllMetrics(ctx context.Context) (map[string]float64, map[string]int64, error) {
	return nil, nil, nil
}

func (s *MockService) SaveMetricsBatch(ctx context.Context, metrics []model.Metrics) error {
	return nil
}
