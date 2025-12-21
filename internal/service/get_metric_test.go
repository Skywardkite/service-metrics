package service

import (
	"context"
	"reflect"
	"testing"

	model "github.com/Skywardkite/service-metrics/internal/model"
	"github.com/Skywardkite/service-metrics/internal/repository"
	"github.com/Skywardkite/service-metrics/internal/repository/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetMetric(t *testing.T) {
	type fields struct {
		store repository.Storage
	}

	type args struct {
		metricType string
		metricName string
	}

	repErr := errors.New("repository error")

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr string
	}{
		{
			name: "success get gauge",
			fields: fields{
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().GetGauge(mock.Anything, "name").Return(float64(34.5), nil)
					return repoMock
				}(),
			},
			args: args{
				metricType: model.Gauge,
				metricName: "name",
			},
			want:    "34.5",
			wantErr: "",
		},
		{
			name: "success get counter",
			fields: fields{
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().GetCounter(mock.Anything, "name").Return(int64(234), nil)
					return repoMock
				}(),
			},
			args: args{
				metricType: model.Counter,
				metricName: "name",
			},
			want:    "234",
			wantErr: "",
		},
		{
			name: "error get gauge",
			fields: fields{
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().GetGauge(mock.Anything, "name").Return(0, repErr)
					return repoMock
				}(),
			},
			args: args{
				metricType: model.Gauge,
				metricName: "name",
			},
			want:    "",
			wantErr: "unknown metric: name",
		},
		{
			name: "error get counter",
			fields: fields{
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().GetCounter(mock.Anything, "name").Return(0, repErr)
					return repoMock
				}(),
			},
			args: args{
				metricType: model.Counter,
				metricName: "name",
			},
			want:    "",
			wantErr: "unknown metric: name",
		},
		{
			name: "error invalid metric type",
			args: args{
				metricType: "invalid Type",
				metricName: "name",
			},
			wantErr: "unsupported metric type: invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricService{
				store: tt.fields.store,
			}

			got, err := s.GetMetric(context.Background(), tt.args.metricType, tt.args.metricName)
			assert.True(t, reflect.DeepEqual(got, tt.want))
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}
