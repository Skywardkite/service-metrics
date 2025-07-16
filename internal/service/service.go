package service

import (
	"errors"
	"strconv"

	"github.com/Skywardkite/service-metrics.git/internal/storage"
)

type MetricService struct {
	store *storage.MemStorage
}

func NewMetricService(s *storage.MemStorage) *MetricService {
	return &MetricService{store: s}
}

func (s *MetricService) UpdateMetric(metricType, metricName, metricValue string) error {
	switch metricType {
      case "gauge":
        value, err := strconv.ParseFloat(metricValue, 64)
        if err != nil {
          return errors.New("invalid gauge value")
        }
		s.store.SetGauge(metricName, value)
        return nil
      case "counter":
        value, err := strconv.ParseInt(metricValue, 10, 64)
        if err != nil {
          return errors.New("invalid counter value")
        }
		s.store.SetCounter(metricName, value)
        return nil
      default:
        return errors.New("unsupported metric type")
    }
}