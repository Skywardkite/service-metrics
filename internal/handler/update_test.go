package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	serverconfig "github.com/Skywardkite/service-metrics/internal/config/server_config"
	model "github.com/Skywardkite/service-metrics/internal/model"
	"github.com/Skywardkite/service-metrics/internal/service"
	"github.com/Skywardkite/service-metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestHandler_UpdateHandler(t *testing.T) {
	type fields struct {
		service *service.MetricService
	}
	type args struct {
		method  string
		metricType      string
		metricName      string
		metricValue     string
	}

	store := storage.NewMemStorage()
	cfg := serverconfig.Config{
		StoreInternal: 30,
	}
    metricService := service.NewMetricService(&cfg, store)

	tests := []struct {
		name   string
		args   	args
		expectedStatus     int
		expectedHeaders    map[string]string
	}{
		{
			name: "successful gauge update",
			args: args{
				method: 		http.MethodPost,
				metricType:		model.Gauge,
				metricName:		"temperature",
				metricValue:	"23.5",
			},
			expectedStatus:    http.StatusOK,
			expectedHeaders: 	map[string]string{
				"Content-Type":   "text/plain; charset=utf-8",
				"Content-Length": "0",
			},
		},
		{
			name: "successful counter update",
			args: args{
				method: 		http.MethodPost,
				metricType:		model.Counter,
				metricName:		"temperature",
				metricValue:	"23",
			},
			expectedStatus:    http.StatusOK,
			expectedHeaders: map[string]string{
				"Content-Type":   "text/plain; charset=utf-8",
				"Content-Length": "0",
			},
		},
		{
			name: "invalid path - too short",
			args: args{
				method: http.MethodPost,
				metricType:		model.Counter,
				metricName:		"",
				metricValue:	"",
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "empty metric name",
			args: args{
				method: http.MethodPost,
				metricType:		model.Counter,
				metricName:		"",
				metricValue:	"23",
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "service returns error",
			args: args{
				method: http.MethodPost,
				metricType:		model.Gauge,
				metricName:		"temperature",
				metricValue:	"invalid",
			},
			expectedStatus:   http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(metricService, nil)

			req := httptest.NewRequest(tt.args.method, "/", nil)
			res := httptest.NewRecorder()

			// Устанавливаем параметры в chi context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricType", tt.args.metricType)
			rctx.URLParams.Add("metricName", tt.args.metricName)
			rctx.URLParams.Add("metricValue", tt.args.metricValue)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))


			h.UpdateHandler(res, req)

			assert.Equal(t, tt.expectedStatus, res.Code)
			
			for key, value := range tt.expectedHeaders {
				assert.Equal(t, value, res.Header().Get(key))
			}
			if tt.expectedStatus == http.StatusOK {
				assert.NotEmpty(t, res.Header().Get("Date"))
				assert.Empty(t, res.Body.String())
			}
		})
	}
}