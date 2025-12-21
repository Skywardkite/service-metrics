package service

import (
	"context"
	"testing"

	model "github.com/Skywardkite/service-metrics/internal/model"
	"github.com/Skywardkite/service-metrics/internal/repository"
	"github.com/Skywardkite/service-metrics/internal/repository/mocks"
	pointers "github.com/Skywardkite/service-metrics/internal/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_SaveMetricsBatch(t *testing.T) {
	type fields struct {
		store repository.Storage
	}

	metrics := []model.Metrics{
		{
			ID:    "23",
			MType: "type",
			Delta: pointers.To(int64(345)),
			Value: nil,
		},
		{
			ID:    "24",
			MType: "type",
			Delta: nil,
			Value: pointers.To(float64(67.8)),
		},
	}

	repErr := errors.New("repository error")

	tests := []struct {
		name    string
		fields  fields
		wantErr string
	}{
		{
			name: "success",
			fields: fields{
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().SetMetricsBatch(mock.Anything, metrics).Return(nil)
					return repoMock
				}(),
			},
			wantErr: "",
		},
		{
			name: "error",
			fields: fields{
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().SetMetricsBatch(mock.Anything, metrics).Return(repErr)
					return repoMock
				}(),
			},
			wantErr: "failed to set metrics batch",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricService{
				store: tt.fields.store,
			}

			err := s.SaveMetricsBatch(context.Background(), metrics)
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}
