package service

import (
	"context"
	"fmt"
	"strconv"

	serverconfig "github.com/Skywardkite/service-metrics/internal/config/server_config"
	model "github.com/Skywardkite/service-metrics/internal/model"
	"github.com/Skywardkite/service-metrics/internal/repository"
	"github.com/Skywardkite/service-metrics/internal/storage"
)

type MetricService struct {
	cfg 	*serverconfig.Config
	store repository.Storage
}

func NewMetricService(cfg 	*serverconfig.Config, s repository.Storage) *MetricService {
	return &MetricService{cfg: cfg, store: s}
}

func (s *MetricService) UpdateMetric(ctx context.Context, metricType, metricName, metricValue string) error {
	switch metricType {
      case model.Gauge:
        value, err := strconv.ParseFloat(metricValue, 64)
        if err != nil {
          return fmt.Errorf("invalid gauge value: %s", metricValue)
        }
		    err = s.store.SetGauge(context.Background(), metricName, value)
        if err != nil {
          return fmt.Errorf("failed to set gauge: %s", metricName)
        }

        if s.cfg.StoreInternal == 0 {
          storage.SaveMetrics(s.cfg.FileStoragePath, map[string]float64{metricName: value}, nil)
        }
        return nil
      case model.Counter:
        value, err := strconv.ParseInt(metricValue, 10, 64)
        if err != nil {
          return fmt.Errorf("invalid counter value: %s", metricValue)
        }
		    err = s.store.SetCounter(ctx, metricName, value)
        if err != nil {
          return fmt.Errorf("failed to set counter: %s", metricName)
        }

        if s.cfg.StoreInternal == 0 && s.cfg.DatabaseDSN == "" {
          storage.SaveMetrics(s.cfg.FileStoragePath, nil, map[string]int64{metricName: value})
        }
        return nil
      default:
        return fmt.Errorf("unsupported metric type: %s", metricType)
    }
}

func (s *MetricService) GetMetric(ctx context.Context, metricType, metricName string) (string, error) {
	switch metricType {
      case model.Gauge:
        value, err := s.store.GetGauge(ctx, metricName)
        if err != nil {
          return "", fmt.Errorf("unknown metric: %s", metricName)
        }
        return strconv.FormatFloat(value, 'f', -1, 64), nil
      case model.Counter:
        value, err := s.store.GetCounter(ctx, metricName)
        if err != nil {
          return "", fmt.Errorf("unknown metric: %s", metricName)
        }
        return strconv.FormatInt(value, 10), nil
      default:
        return "", fmt.Errorf("unsupported metric type: %s", metricType)
    }
}

func (s *MetricService) GetAllMetrics(ctx context.Context) (map[string]float64, map[string]int64, error) {
  gauges, counters, err := s.store.GetMetrics(ctx)
  if err != nil {
    return nil, nil, fmt.Errorf("failed to get metrics: %s", err)
  }
	return gauges, counters, nil
} 