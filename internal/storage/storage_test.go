package storage

import (
	"testing"

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
			want:   65.5,
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
			want:   20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsStorage{
				Gauge:   tt.fields.Gauge,
				Counter: tt.fields.Counter,
			}
			s.SetGauge(tt.args.name, tt.args.value)

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
			want:   75,
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
			want:   20,
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
			want:   -20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsStorage{
				Gauge:   tt.fields.Gauge,
				Counter: tt.fields.Counter,
			}
			s.SetCounter(tt.args.name, tt.args.value)

			assert.Equal(t, s.Counter[tt.args.name], tt.want)
		})
	}
}
