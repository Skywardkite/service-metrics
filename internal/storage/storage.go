package storage

type MetricsStorage struct {
	Gauge 		map[string]float64
	Counter 	map[string]int64
}

func NewMetricsStorage() *MetricsStorage {
	return &MetricsStorage{
		Gauge: 		make(map[string]float64),
		Counter: 	make(map[string]int64),
	}
}

func (s *MetricsStorage) SetGauge(name string, value float64) {
	s.Gauge[name] = value
}

func (s *MetricsStorage) SetCounter(name string, value int64) {
	s.Counter[name] += value
}