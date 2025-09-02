package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Skywardkite/service-metrics/internal/config/server_config"
	model "github.com/Skywardkite/service-metrics/internal/model"
	"github.com/Skywardkite/service-metrics/internal/repository"
	"github.com/Skywardkite/service-metrics/internal/storage"
)

type MetricService struct {
	cfg 	*server_config.Config
	store repository.Storage
}

func NewMetricService(cfg 	*server_config.Config, s repository.Storage) *MetricService {
	return &MetricService{cfg: cfg, store: s}
}

func (s *MetricService) UpdateMetric(ctx context.Context, metricType, metricName, metricValue string) error {
	switch metricType {
		case model.Gauge:
		  value, err := strconv.ParseFloat(metricValue, 64)
		  if err != nil {
				return fmt.Errorf("invalid gauge value: %s", metricValue)
		  }

		  err = withRetry(func() error {
				return s.store.SetGauge(ctx, metricName, value)
		  })
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
		
			err = withRetry(func() error {
				return s.store.SetCounter(ctx, metricName, value)
		  })
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
    var value float64
    err := withRetry(func() error {
      var err error
      value, err = s.store.GetGauge(ctx, metricName)
      return err
    })
    if err != nil {
      return "", fmt.Errorf("unknown metric: %s", metricName)
    }
    return strconv.FormatFloat(value, 'f', -1, 64), nil

  case model.Counter:
    var value int64
    err := withRetry(func() error {
    	var err error
      value, err = s.store.GetCounter(ctx, metricName)
      return err
    })
    if err != nil {
			return "", fmt.Errorf("unknown metric: %s", metricName)
    }
    return strconv.FormatInt(value, 10), nil

  default:
    return "", fmt.Errorf("unsupported metric type: %s", metricType)
  }
}

func (s *MetricService) GetAllMetrics(ctx context.Context) (map[string]float64, map[string]int64, error) {
  var (
    gauges   map[string]float64
    counters map[string]int64
  )

  err := withRetry(func() error {
    var err error
    gauges, counters, err = s.store.GetMetrics(ctx)
    return err
	})
  if err != nil {
    return nil, nil, fmt.Errorf("failed to get metrics: %w", err)
  }

  return gauges, counters, nil
}

func (s *MetricService) SaveMetricsBatch(ctx context.Context, metrics []model.Metrics) error {
	err := withRetry(func() error {
		return s.store.SetMetricsBatch(ctx, metrics)
	})
	if err != nil {
		return fmt.Errorf("failed to set metrics batch: %w", err)
	}
	
	return nil
}