package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Skywardkite/service-metrics/internal/service"
	"github.com/Skywardkite/service-metrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandler_UpdateHandler(t *testing.T) {
	type fields struct {
		service *service.MetricService
	}
	type args struct {
		method  string
		path    string
	}

	store := storage.NewMetricsStorage()
    metricService := service.NewMetricService(store)

	tests := []struct {
		name   string
		args   	args
		expectedStatus     int
		expectedHeaders    map[string]string
	}{
		{
			name: "successful gauge update",
			args: args{
				method: http.MethodPost,
				path:   "/update/gauge/temperature/23.5",
			},
			expectedStatus:    http.StatusOK,
			expectedHeaders: map[string]string{
				"Content-Type":   "text/plain; charset=utf-8",
				"Content-Length": "0",
			},
		},
		{
			name: "successful counter update",
			args: args{
				method: http.MethodPost,
				path:   "/update/counter/temperature/23",
			},
			expectedStatus:    http.StatusOK,
			expectedHeaders: map[string]string{
				"Content-Type":   "text/plain; charset=utf-8",
				"Content-Length": "0",
			},
		},
		{
			name: "method not allowed",
			args: args{
				method: http.MethodGet,
				path:   "/update/counter/temperature/23",
			},
			expectedStatus:    http.StatusMethodNotAllowed,
		},
		{
			name: "invalid path - too short",
			args: args{
				method: http.MethodPost,
				path:   "/update/counter",
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "invalid path - missing update prefix",
			args: args{
				method: http.MethodPost,
				path:   "/counter/temperature/23",
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "empty metric name",
			args: args{
				method: http.MethodPost,
				path:   "/update/counter//23",
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:              "service returns error",
			args: args{
				method: http.MethodPost,
				path:   "/update/gauge/temperature/invalid",
			},
			expectedStatus:   http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(metricService)

			req := httptest.NewRequest(tt.args.method, tt.args.path, nil)
			res := httptest.NewRecorder()

			h.UpdateHandler(res, req)

			assert.Equal(t, tt.expectedStatus, res.Code)
			for key, value := range tt.expectedHeaders {
				assert.Equal(t, value, res.Header().Get(key))
			}
			if tt.expectedStatus == http.StatusOK {
				assert.NotEmpty(t, res.Header().Get("Date"))
			}
		})
	}
}