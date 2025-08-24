package storage

import (
	"context"
	"errors"
	"maps"

	model "github.com/Skywardkite/service-metrics/internal/model"
)

var (
	ErrCounterNotFound = errors.New("metric counter not found")
	ErrGaugeNotFound = errors.New("metric gauge not found")
	ErrConnectToDB = errors.New("metric counter not found")
	ErrUnsupportedMetricType = errors.New("unsupported metric type")
)
type MemStorage struct {
	Gauge 		map[string]float64
	Counter 	map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge: 		make(map[string]float64),
		Counter: 	make(map[string]int64),
	}
}

func (s *MemStorage) SetGauge(ctx context.Context, name string, value float64) error {
	s.Gauge[name] = value
	return nil
}

func (s *MemStorage) SetCounter(ctx context.Context, name string, value int64) error {
	s.Counter[name] += value
	return nil
}

func (s *MemStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	value, ok := s.Gauge[name]
	if !ok {
		return 0, ErrGaugeNotFound
	}
	return value, nil
}

func (s *MemStorage) GetCounter(ctx context.Context, name string) (int64, error){
	value, ok := s.Counter[name]
	if !ok {
		return 0, ErrCounterNotFound
	}
	return value, nil
}

func (s *MemStorage) GetMetrics(ctx context.Context) (map[string]float64, map[string]int64, error) {
	gauges := make(map[string]float64)
	maps.Copy(gauges, s.Gauge)

	counters := make(map[string]int64)
	maps.Copy(counters, s.Counter)
		
	return gauges, counters, nil
}

func (s *MemStorage) Ping() error {
	// Соединение с бд отсутствует
	return ErrConnectToDB
}

func (s *MemStorage) SetMetricsBatch(ctx context.Context, metrics []model.Metrics) error {
	for _, metric := range metrics {
		switch metric.MType {
		case model.Gauge:
			s.Gauge[metric.ID] = *metric.Value
		case model.Counter:
			s.Counter[metric.ID] += *metric.Delta
		default: 
			return ErrUnsupportedMetricType
		}
	}
	return nil
}