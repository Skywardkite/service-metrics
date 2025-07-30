package storage

import (
	"maps"
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

func (s *MemStorage) SetGauge(name string, value float64) {
	s.Gauge[name] = value
}

func (s *MemStorage) SetCounter(name string, value int64) {
	s.Counter[name] += value
}

func (s *MemStorage) GetGauge(name string) (float64, bool) {
	value, ok := s.Gauge[name]
	return value, ok
}

func (s *MemStorage) GetCounter(name string) (int64, bool){
	value, ok := s.Counter[name]
	return value, ok
}

func (s *MemStorage) GetMetrics() (map[string]float64, map[string]int64) {
	gauges := make(map[string]float64)
	maps.Copy(gauges, s.Gauge)

	counters := make(map[string]int64)
	maps.Copy(counters, s.Counter)
		
	return gauges, counters
}
