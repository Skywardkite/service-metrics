// Package service - вся бизнес логика сервиса, который получает, хранит и отдает метрики
package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Skywardkite/service-metrics/internal/config/server_config"
	model "github.com/Skywardkite/service-metrics/internal/model"
	"github.com/Skywardkite/service-metrics/internal/repository"
)

type MetricService struct {
	Cfg   *server_config.Config
	store repository.Storage
}

func NewMetricService(cfg *server_config.Config, s repository.Storage) *MetricService {
	return &MetricService{
		Cfg:   cfg,
		store: s,
	}
}

type MetricServiceInterface interface {
	UpdateMetric(ctx context.Context, metricType, metricName, metricValue string) error
	GetMetric(ctx context.Context, metricType, metricName string) (string, error)
	GetAllMetrics(ctx context.Context) (map[string]float64, map[string]int64, error)
	SaveMetricsBatch(ctx context.Context, metrics []model.Metrics) error
}

// UpdateMetric обновляет метрики.
// Обновление проходит с ретраем.
// Если при запуске сервиса не выставляли StoreInternal, после каждого обновления происходит сохранение метрик в storage.
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

		if s.Cfg.StoreInternal == 0 {
			err = s.store.SaveMetrics(s.Cfg.FileStoragePath, map[string]float64{metricName: value}, nil)
			if err != nil {
				return err
			}
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

		if s.Cfg.StoreInternal == 0 && s.Cfg.DatabaseDSN == "" {
			err = s.store.SaveMetrics(s.Cfg.FileStoragePath, nil, map[string]int64{metricName: value})
			if err != nil {
				return err
			}
		}

		return nil

	default:
		return fmt.Errorf("unsupported metric type: %s", metricType)
	}
}

// GetMetric отдает значение метрики по ее типу и названию.
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

// GetAllMetrics отдает все метрики из store что знает.
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

// SaveMetricsBatch сохраняет сразу несколько метрик вида model.Metrics.
func (s *MetricService) SaveMetricsBatch(ctx context.Context, metrics []model.Metrics) error {
	err := withRetry(func() error {
		return s.store.SetMetricsBatch(ctx, metrics)
	})
	if err != nil {
		return fmt.Errorf("failed to set metrics batch: %w", err)
	}

	return nil
}
