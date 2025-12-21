package storage

import (
	"context"
	"testing"

	model "github.com/Skywardkite/service-metrics/internal/model"
	pointers "github.com/Skywardkite/service-metrics/internal/utils"
	"gotest.tools/v3/assert"
)

func TestMemStorage_SetGauge(t *testing.T) {
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
			s := &MemStorage{
				Gauge:   tt.fields.Gauge,
				Counter: tt.fields.Counter,
			}
			s.SetGauge(context.Background(), tt.args.name, tt.args.value)

			assert.Equal(t, s.Gauge[tt.args.name], tt.want)
		})
	}
}

func TestMemStorage_SetCounter(t *testing.T) {
	type fields struct {
		Gauge   map[string]float64
		Counter map[string]int64
	}
	type args struct {
		name  string
		value int64
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
				name:  "some_metric",
				value: 65,
			},
			want: 75,
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
		{
			name: "success_negative_value",
			fields: fields{
				Gauge:   make(map[string]float64),
				Counter: make(map[string]int64),
			},
			args: args{
				name:  "some_metric",
				value: -20,
			},
			want: -20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Gauge:   tt.fields.Gauge,
				Counter: tt.fields.Counter,
			}
			s.SetCounter(context.Background(), tt.args.name, tt.args.value)

			assert.Equal(t, s.Counter[tt.args.name], tt.want)
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type fields struct {
		Gauge map[string]float64
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				Gauge: map[string]float64{"some_metric": 1.5},
			},
			args: args{
				name: "some_metric",
			},
			want: 1.5,
		},
		{
			name: "error",
			fields: fields{
				Gauge: map[string]float64{"some_metric": 1.5},
			},
			args: args{
				name: "different_metric",
			},
			want:    0,
			wantErr: ErrGaugeNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Gauge: tt.fields.Gauge,
			}

			got, err := s.GetGauge(context.Background(), tt.args.name)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	type fields struct {
		Counter map[string]int64
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				Counter: map[string]int64{"some_metric": 15},
			},
			args: args{
				name: "some_metric",
			},
			want: 15,
		},
		{
			name: "error",
			fields: fields{
				Counter: map[string]int64{"some_metric": 15},
			},
			args: args{
				name: "different_metric",
			},
			want:    0,
			wantErr: ErrCounterNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Counter: tt.fields.Counter,
			}

			got, err := s.GetCounter(context.Background(), tt.args.name)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestMemStorage_GetMetrics(t *testing.T) {
	tests := []struct {
		name         string
		storage      *MemStorage
		wantCounters map[string]int64
		wantGauges   map[string]float64
		wantErr      error
	}{
		{
			name: "success",
			storage: &MemStorage{
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
			wantErr: nil,
		},
		{
			name: "success empty",
			storage: &MemStorage{
				Gauge:   map[string]float64{},
				Counter: map[string]int64{},
			},
			wantGauges:   map[string]float64{},
			wantCounters: map[string]int64{},
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gauges, counters, err := tt.storage.GetMetrics(context.Background())
			assert.Equal(t, tt.wantErr, err)
			assert.DeepEqual(t, tt.wantGauges, gauges)
			assert.DeepEqual(t, tt.wantCounters, counters)
		})
	}
}

func TestMemStorage_GetMetrics_ReturnsCopies(t *testing.T) {
	st := &MemStorage{
		Gauge: map[string]float64{
			"Alloc": 100,
		},
		Counter: map[string]int64{
			"PollCount": 1,
		},
	}

	gauges, counters, err := st.GetMetrics(context.Background())
	assert.NilError(t, err)

	gauges["Alloc"] = 999
	counters["PollCount"] = 999

	// Внутреннее состояние не должно измениться при мутировании результатов
	assert.Equal(t, float64(100), st.Gauge["Alloc"])
	assert.Equal(t, int64(1), st.Counter["PollCount"])
}

func TestMemStorage_SetMetricsBatch(t *testing.T) {
	tests := []struct {
		name        string
		initial     *MemStorage
		metrics     []model.Metrics
		wantGauges  map[string]float64
		wantCounter map[string]int64
		wantErr     error
	}{
		{
			name: "success gauge and counter",
			initial: &MemStorage{
				Gauge:   map[string]float64{},
				Counter: map[string]int64{},
			},
			metrics: []model.Metrics{
				{
					ID:    "Alloc",
					MType: model.Gauge,
					Value: pointers.To(float64(1.5)),
				},
				{
					ID:    "PollCount",
					MType: model.Counter,
					Delta: pointers.To(int64(10)),
				},
			},
			wantGauges: map[string]float64{
				"Alloc": 1.5,
			},
			wantCounter: map[string]int64{
				"PollCount": 10,
			},
			wantErr: nil,
		},
		{
			name: "counter increments existing value",
			initial: &MemStorage{
				Gauge: map[string]float64{},
				Counter: map[string]int64{
					"PollCount": 5,
				},
			},
			metrics: []model.Metrics{
				{
					ID:    "PollCount",
					MType: model.Counter,
					Delta: pointers.To(int64(3)),
				},
			},
			wantGauges: map[string]float64{},
			wantCounter: map[string]int64{
				"PollCount": 8,
			},
			wantErr: nil,
		},
		{
			name: "unsupported metric type",
			initial: &MemStorage{
				Gauge:   map[string]float64{},
				Counter: map[string]int64{},
			},
			metrics: []model.Metrics{
				{
					ID:    "Unknown",
					MType: "unknown",
				},
			},
			wantGauges:  map[string]float64{},
			wantCounter: map[string]int64{},
			wantErr:     ErrUnsupportedMetricType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.initial.SetMetricsBatch(context.Background(), tt.metrics)

			assert.Equal(t, tt.wantErr, err)
			assert.DeepEqual(t, tt.wantGauges, tt.initial.Gauge)
			assert.DeepEqual(t, tt.wantCounter, tt.initial.Counter)
		})
	}
}
