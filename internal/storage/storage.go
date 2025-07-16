package storage

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