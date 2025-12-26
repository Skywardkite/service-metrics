package agent

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestAgentMetrics_SetGauge(t *testing.T) {
	type fields struct {
		Gauge   map[string]float64
		Counter map[string]int64
	}
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "success_set_metric",
			fields: fields{
				Gauge:   map[string]float64{"some_metric": 1.5},
				Counter: make(map[string]int64),
			},
			args: args{
				name:  "some_metric",
				value: 65.5,
			},
			want: 65.5,
		},
		{
			name: "success_new_metric",
			fields: fields{
				Gauge:   make(map[string]float64),
				Counter: make(map[string]int64),
			},
			args: args{
				name:  "some_metric",
				value: 20,
			},
			want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AgentMetrics{
				Gauge:   tt.fields.Gauge,
				Counter: tt.fields.Counter,
			}
			s.SetAgentGauge(tt.args.name, tt.args.value)

			assert.Equal(t, s.Gauge[tt.args.name], tt.want)
		})
	}
}
func TestAgentMetrics_SetCounter(t *testing.T) {
	type fields struct {
		Gauge   map[string]float64
		Counter map[string]int64
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			name: "success_set_metric",
			fields: fields{
				Gauge:   make(map[string]float64),
				Counter: map[string]int64{"some_metric": 10},
			},
			args: args{
				name: "some_metric",
			},
			want: 11,
		},
		{
			name: "success_new_metric",
			fields: fields{
				Gauge:   make(map[string]float64),
				Counter: make(map[string]int64),
			},
			args: args{
				name: "some_metric",
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AgentMetrics{
				Gauge:   tt.fields.Gauge,
				Counter: tt.fields.Counter,
			}
			s.SetAgentCounter(tt.args.name)

			assert.Equal(t, s.Counter[tt.args.name], tt.want)
		})
	}
}
func TestMemStorage_GetMetrics(t *testing.T) {
	tests := []struct {
		name         string
		storage      *AgentMetrics
		wantCounters map[string]int64
		wantGauges   map[string]float64
	}{
		{
			name: "success",
			storage: &AgentMetrics{
				Gauge: map[string]float64{
					"some_gauge": 1.5,
				},
				Counter: map[string]int64{
					"some_counter": 15,
				},
			},
			wantGauges: map[string]float64{
				"some_gauge": 1.5,
			},
			wantCounters: map[string]int64{
				"some_counter": 15,
			},
		},
		{
			name: "success empty",
			storage: &AgentMetrics{
				Gauge:   map[string]float64{},
				Counter: map[string]int64{},
			},
			wantGauges:   map[string]float64{},
			wantCounters: map[string]int64{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gauges, counters := tt.storage.GetAgentMetrics()
			assert.DeepEqual(t, tt.wantGauges, gauges)
			assert.DeepEqual(t, tt.wantCounters, counters)
		})
	}
}
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
func TestPollSystemMetrics(t *testing.T) {
	storage := &AgentMetrics{
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
	}

	PollSystemMetrics(storage)

	// Проверяем метрики памяти
	total, ok := storage.Gauge["TotalMemory"]
	assert.Assert(t, ok)
	assert.Assert(t, total > 0)

	free, ok := storage.Gauge["FreeMemory"]
	assert.Assert(t, ok)
	assert.Assert(t, free >= 0)

	// Проверяем CPU метрику
	foundCPU := false
	for name := range storage.Gauge {
		if len(name) > 14 && name[:14] == "CPUutilization" {
			foundCPU = true
			assert.Assert(t, storage.Gauge[name] >= 0.0)
			assert.Assert(t, storage.Gauge[name] <= 100.0)
		}
	}
	assert.Assert(t, foundCPU, "no CPU metrics were set")
}
