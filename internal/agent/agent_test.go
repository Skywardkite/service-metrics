package agent

import (
	"testing"
	"time"
)


func TestPollRuntimeMetrics(t *testing.T) {
	storage := NewAgentMetrics()
	PollRuntimeMetrics(storage)

	requiredGauges := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
		"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
		"MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
		"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
		"Sys", "TotalAlloc", "RandomValue",
	}

	for _, metric := range requiredGauges {
		if _, ok := storage.Gauge[metric]; !ok {
			t.Errorf("Metric %s was not recorded", metric)
		}
	}

	if storage.Counter["PollCount"] != 1 {
		t.Errorf("PollCount should be 1, got %d", storage.Counter["PollCount"])
	}

	// RandomValue в допустимом диапазоне
	if storage.Gauge["RandomValue"] < 0 || storage.Gauge["RandomValue"] > 1000 {
		t.Errorf("RandomValue should be between 0 and 1000, got %f", storage.Gauge["RandomValue"])
	}

	// Метрики времени обновления актуальны
	if storage.Gauge["LastGC"] > float64(time.Now().UnixNano()) {
		t.Error("LastGC value is in the future")
	}
}