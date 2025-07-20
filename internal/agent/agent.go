package agent

import (
	"maps"
	"math/rand/v2"
	"runtime"
)

type AgentMetrics struct {
	Gauge 		map[string]float64
	Counter 	map[string]int64
}

func NewAgentMetrics() *AgentMetrics {
	return &AgentMetrics{
		Gauge: 		make(map[string]float64),
		Counter: 	make(map[string]int64),
	}
}

func (s *AgentMetrics) SetAgentGauge(name string, value float64) {
	s.Gauge[name] = value
}

func (s *AgentMetrics) SetAgentCounter(name string) {
	s.Counter[name]++
}

func (s *AgentMetrics) GetAgentMetrics() (map[string]float64, map[string]int64) {
	g := s.Gauge
	c := s.Counter
	maps.Copy(g, s.Gauge)
	maps.Copy(c, s.Counter)
	return g, c
}

func PollRuntimeMetrics(storage *AgentMetrics) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	storage.SetAgentGauge("Alloc", float64(m.Alloc))
	storage.SetAgentGauge("BuckHashSys", float64(m.BuckHashSys))
	storage.SetAgentGauge("Frees", float64(m.Frees))
	storage.SetAgentGauge("GCCPUFraction", m.GCCPUFraction)
	storage.SetAgentGauge("GCSys", float64(m.GCSys))
	storage.SetAgentGauge("HeapAlloc", float64(m.HeapAlloc))
	storage.SetAgentGauge("HeapIdle", float64(m.HeapIdle))
	storage.SetAgentGauge("HeapInuse", float64(m.HeapInuse))
	storage.SetAgentGauge("HeapObjects", float64(m.HeapObjects))
	storage.SetAgentGauge("HeapReleased", float64(m.HeapReleased))
	storage.SetAgentGauge("HeapSys", float64(m.HeapSys))
	storage.SetAgentGauge("LastGC", float64(m.LastGC))
	storage.SetAgentGauge("Lookups", float64(m.Lookups))
	storage.SetAgentGauge("MCacheInuse", float64(m.MCacheInuse))
	storage.SetAgentGauge("MCacheSys", float64(m.MCacheSys))
	storage.SetAgentGauge("MSpanInuse", float64(m.MSpanInuse))
	storage.SetAgentGauge("MSpanSys", float64(m.MSpanSys))
	storage.SetAgentGauge("Mallocs", float64(m.Mallocs))
	storage.SetAgentGauge("NextGC", float64(m.NextGC))
	storage.SetAgentGauge("NumForcedGC", float64(m.NumForcedGC))
	storage.SetAgentGauge("NumGC", float64(m.NumGC))
	storage.SetAgentGauge("OtherSys", float64(m.OtherSys))
	storage.SetAgentGauge("PauseTotalNs", float64(m.PauseTotalNs))
	storage.SetAgentGauge("StackInuse", float64(m.StackInuse))
	storage.SetAgentGauge("StackSys", float64(m.StackSys))
	storage.SetAgentGauge("Sys", float64(m.Sys))
	storage.SetAgentGauge("TotalAlloc", float64(m.TotalAlloc))

	storage.SetAgentGauge("RandomValue", rand.Float64()*1000)
	storage.SetAgentCounter("PollCount")
}