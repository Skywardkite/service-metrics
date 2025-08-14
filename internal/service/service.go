package service

import (
	"fmt"
	"strconv"

	serverconfig "github.com/Skywardkite/service-metrics/internal/config/serverConfig"
	"github.com/Skywardkite/service-metrics/internal/storage"
)

type MetricService struct {
	cfg 	*serverconfig.Config
	store *storage.MemStorage
}

func NewMetricService(cfg 	*serverconfig.Config, s *storage.MemStorage) *MetricService {
	return &MetricService{cfg: cfg, store: s}
}

func (s *MetricService) UpdateMetric(metricType, metricName, metricValue string) error {
	switch metricType {
      case "gauge":
        value, err := strconv.ParseFloat(metricValue, 64)
        if err != nil {
          return fmt.Errorf("invalid gauge value: %s", metricValue)
        }
		    s.store.SetGauge(metricName, value)

        if s.cfg.StoreInternal == 0 {
          storage.SaveMetrics(s.cfg.FileStoragePath, map[string]float64{metricName: value}, nil)
        }
        return nil
      case "counter":
        value, err := strconv.ParseInt(metricValue, 10, 64)
        if err != nil {
          return fmt.Errorf("invalid counter value: %s", metricValue)
        }
		    s.store.SetCounter(metricName, value)

        if s.cfg.StoreInternal == 0 {
          storage.SaveMetrics(s.cfg.FileStoragePath, nil, map[string]int64{metricName: value})
        }
        return nil
      default:
        return fmt.Errorf("unsupported metric type: %s", metricType)
    }
}

func (s *MetricService) GetMetric(metricType, metricName string) (string, error) {
	switch metricType {
      case "gauge":
        value, ok := s.store.GetGauge(metricName)
        if !ok {
          return "", fmt.Errorf("unknown metric: %s", metricName)
        }
        return strconv.FormatFloat(value, 'f', -1, 64), nil
      case "counter":
        value, ok := s.store.GetCounter(metricName)
        if !ok {
          return "", fmt.Errorf("unknown metric: %s", metricName)
        }
        return strconv.FormatInt(value, 10), nil
      default:
        return "", fmt.Errorf("unsupported metric type: %s", metricType)
    }
}

func (s *MetricService) GetAllMetrics() (map[string]float64, map[string]int64) {
	gauges, counters := s.store.GetMetrics()
	return gauges, counters
} 