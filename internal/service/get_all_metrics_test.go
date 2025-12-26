package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/Skywardkite/service-metrics/internal/repository"
	"github.com/Skywardkite/service-metrics/internal/repository/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetAllMetrics(t *testing.T) {
	type fields struct {
		store repository.Storage
	}

	repErr := errors.New("repository error")

	tests := []struct {
		name        string
		fields      fields
		wantConters map[string]int64
		wantGauges  map[string]float64
		wantErr     string
	}{
		{
			name: "success",
			fields: fields{
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().GetMetrics(mock.Anything).Return(map[string]float64{"name gauge": 34.5}, map[string]int64{"name counter": 345}, nil)
					return repoMock
				}(),
			},
			wantConters: map[string]int64{"name counter": 345},
			wantGauges:  map[string]float64{"name gauge": 34.5},
			wantErr:     "",
		},
		{
			name: "error get gauge",
			fields: fields{
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().GetMetrics(mock.Anything).Return(map[string]float64{}, map[string]int64{}, repErr)
					return repoMock
				}(),
			},
			wantErr: "failed to get metrics:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricService{
				store: tt.fields.store,
			}

			gotGauges, gotCounter, err := s.GetAllMetrics(context.Background())
			assert.True(t, reflect.DeepEqual(gotCounter, tt.wantConters))
			assert.True(t, reflect.DeepEqual(gotGauges, tt.wantGauges))
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}
