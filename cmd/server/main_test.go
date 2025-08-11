package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	serverconfig "github.com/Skywardkite/service-metrics/internal/config/serverConfig"
	"github.com/Skywardkite/service-metrics/internal/handler"
	"github.com/Skywardkite/service-metrics/internal/service"
	"github.com/Skywardkite/service-metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"gotest.tools/assert"
)

func TestMain(t *testing.T) {
	store := storage.NewMemStorage()
	cfg := serverconfig.Config{
		StoreInternal: 30,
	}
	metricService := service.NewMetricService(&cfg, store)
	h := handler.NewHandler(metricService)

	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateHandler)

	tests := []struct {
		name       string
		url        string
		method     string
		wantStatus int
	}{
		{
			name:       "valid update request",
			url:        "/update/gauge/test_metric/123.45",
			method:     http.MethodPost,
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid path",
			url:        "/invalid/path",
			method:     http.MethodPost,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "wrong method",
			url:        "/update/gauge/test_metric/123.45",
			method:     http.MethodGet,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Errorf("failed to create request: %v", err)
				return
			}
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantStatus, rr.Code, "unexpected status code")
		})
	}
}