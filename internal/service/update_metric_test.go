package service

import (
	"context"
	"testing"

	"github.com/Skywardkite/service-metrics/internal/config/server_config"
	model "github.com/Skywardkite/service-metrics/internal/model"
	"github.com/Skywardkite/service-metrics/internal/repository"
	"github.com/Skywardkite/service-metrics/internal/repository/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_UpdateMetric(t *testing.T) {
	type fields struct {
		Cfg   *server_config.Config
		store repository.Storage
	}

	type args struct {
		metricType  string
		metricName  string
		metricValue string
	}

	repErr := errors.New("repository error")

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr string
	}{
		{
			name: "success update gauge",
			fields: fields{
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().SetGauge(mock.Anything, "name", float64(67.8)).Return(nil)
					return repoMock
				}(),
				Cfg: &server_config.Config{
					StoreInternal: 10,
				},
			},
			args: args{
				metricType:  model.Gauge,
				metricName:  "name",
				metricValue: "67.8",
			},
			wantErr: "",
		},
		{
			name: "success update counter",
			fields: fields{
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().SetCounter(mock.Anything, "name", int64(234)).Return(nil)
					return repoMock
				}(),
				Cfg: &server_config.Config{
					StoreInternal: 10,
				},
			},
			args: args{
				metricType:  model.Counter,
				metricName:  "name",
				metricValue: "234",
			},
			wantErr: "",
		},
		{
			name: "error update gauge",
			fields: fields{
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().SetGauge(mock.Anything, "name", float64(67.8)).Return(repErr)
					return repoMock
				}(),
				Cfg: &server_config.Config{
					StoreInternal: 10,
				},
			},
			args: args{
				metricType:  model.Gauge,
				metricName:  "name",
				metricValue: "67.8",
			},
			wantErr: "failed to set gauge: name",
		},
		{
			name: "error update counter",
			fields: fields{
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().SetCounter(mock.Anything, "name", int64(234)).Return(repErr)
					return repoMock
				}(),
				Cfg: &server_config.Config{
					StoreInternal: 10,
				},
			},
			args: args{
				metricType:  model.Counter,
				metricName:  "name",
				metricValue: "234",
			},
			wantErr: "failed to set counter: name",
		},
		{
			name: "success update gauge with SaveMetrics",
			fields: fields{
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().SetGauge(mock.Anything, "name", float64(67.8)).Return(nil)
					repoMock.EXPECT().SaveMetrics("fileStoragePath", map[string]float64{"name": float64(67.8)}, map[string]int64(nil)).Return(nil)
					return repoMock
				}(),
				Cfg: &server_config.Config{
					StoreInternal:   0,
					FileStoragePath: "fileStoragePath",
				},
			},
			args: args{
				metricType:  model.Gauge,
				metricName:  "name",
				metricValue: "67.8",
			},
			wantErr: "",
		},
		{
			name: "success update counter with SaveMetrics",
			fields: fields{
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().SetCounter(mock.Anything, "name", int64(234)).Return(nil)
					repoMock.EXPECT().SaveMetrics("fileStoragePath", map[string]float64(nil), map[string]int64{"name": int64(234)}).Return(nil)
					return repoMock
				}(),
				Cfg: &server_config.Config{
					StoreInternal:   0,
					FileStoragePath: "fileStoragePath",
				},
			},
			args: args{
				metricType:  model.Counter,
				metricName:  "name",
				metricValue: "234",
			},
			wantErr: "",
		},
		{
			name: "error update gauge with SaveMetrics",
			fields: fields{
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().SetGauge(mock.Anything, "name", float64(67.8)).Return(nil)
					repoMock.EXPECT().SaveMetrics("fileStoragePath", map[string]float64{"name": float64(67.8)}, map[string]int64(nil)).Return(repErr)
					return repoMock
				}(),
				Cfg: &server_config.Config{
					StoreInternal:   0,
					FileStoragePath: "fileStoragePath",
				},
			},
			args: args{
				metricType:  model.Gauge,
				metricName:  "name",
				metricValue: "67.8",
			},
			wantErr: "repository error",
		},
		{
			name: "error update counter with SaveMetrics",
			fields: fields{
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().SetCounter(mock.Anything, "name", int64(234)).Return(nil)
					repoMock.EXPECT().SaveMetrics("fileStoragePath", map[string]float64(nil), map[string]int64{"name": int64(234)}).Return(repErr)
					return repoMock
				}(),
				Cfg: &server_config.Config{
					StoreInternal:   0,
					FileStoragePath: "fileStoragePath",
				},
			},
			args: args{
				metricType:  model.Counter,
				metricName:  "name",
				metricValue: "234",
			},
			wantErr: "repository error",
		},
		{
			name:   "error invalid gauge",
			fields: fields{},
			args: args{
				metricType:  model.Gauge,
				metricName:  "name",
				metricValue: "invalid",
			},
			wantErr: "invalid gauge value: invalid",
		},
		{
			name:   "error invalid counter",
			fields: fields{},
			args: args{
				metricType:  model.Counter,
				metricName:  "name",
				metricValue: "invalid",
			},
			wantErr: "invalid counter value: invalid",
		},
		{
			name: "error invalid metric type",
			args: args{
				metricType:  "invalid Type",
				metricName:  "name",
				metricValue: "234",
			},
			wantErr: "unsupported metric type: invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricService{
				Cfg:   tt.fields.Cfg,
				store: tt.fields.store,
			}

			err := s.UpdateMetric(context.Background(), tt.args.metricType, tt.args.metricName, tt.args.metricValue)
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}
