package service

import (
	"testing"

	serverconfig "github.com/Skywardkite/service-metrics/internal/config/serverConfig"
	"github.com/Skywardkite/service-metrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestMetricService_UpdateMetric(t *testing.T) {
	type fields struct {
		store *storage.MemStorage
		cfg    *serverconfig.Config
	}
	type args struct {
		metricType  string
		metricName  string
		metricValue string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		errorMessage string
	}{
		{
			name: "successful gauge update",
			fields: fields{
				store: &storage.MemStorage{
					Gauge:   map[string]float64{"some_metric": 1.5},
					Counter: make(map[string]int64),
				},
				cfg: &serverconfig.Config{
					StoreInternal: 0,
				},
			},
			args: args{
				metricType: "gauge",
				metricName: "some_metric_1",
				metricValue: "65.5",
			},
			wantErr: false,
		},
		{
			name: "successful counter update",
			fields: fields{
				store: &storage.MemStorage{
					Gauge:   map[string]float64{"some_metric": 1.5},
					Counter: make(map[string]int64),
				},
				cfg: &serverconfig.Config{
					StoreInternal: 0,
				},
			},
			args: args{
				metricType: "counter",
				metricName: "some_metric_1",
				metricValue: "20",
			},
			wantErr: false,
		},
		{
			name: "invalid gauge value",
			args: args{
				metricType: "gauge",
				metricName: "some_metric_1",
				metricValue: "invalid",
			},
			wantErr:  true,
			errorMessage: "invalid gauge value: invalid",
		},
		{
			name: "invalid counter value",
			args: args{
				metricType: "counter",
				metricName: "some_metric_1",
				metricValue: "invalid",
			},
			wantErr:  true,
			errorMessage: "invalid counter value: invalid",
		},
		{
			name: "unsupported metric type",
			args: args{
				metricType: "wrongCounter",
				metricName: "some_metric_1",
				metricValue: "20",
			},
			wantErr:  true,
			errorMessage: "unsupported metric type: wrongCounter",
		},
		{
			name: "empty metric type",
			args: args{
				metricType: "",
				metricName: "some_metric_1",
				metricValue: "20",
			},
			wantErr:  true,
			errorMessage: "unsupported metric type: ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricService{
				store: tt.fields.store,
				cfg: tt.fields.cfg,
			}
			err := s.UpdateMetric(tt.args.metricType, tt.args.metricName, tt.args.metricValue)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}