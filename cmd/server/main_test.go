package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Skywardkite/service-metrics/internal/handler"
	"github.com/Skywardkite/service-metrics/internal/service"
	"github.com/Skywardkite/service-metrics/internal/storage"
)

func TestMain(t *testing.T) {
	store := storage.NewMetricsStorage()
	metricService := service.NewMetricService(store)
	h := handler.NewHandler(metricService)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/update/") {
			h.UpdateHandler(w, r)
			return
		}
		http.NotFound(w, r)
	}))
	defer ts.Close()

	tests := []struct {
		name       string
		url        string
		method     string
		wantStatus int
	}{
		{
			name:       "valid update request",
			url:        ts.URL + "/update/gauge/test_metric/123.45",
			method:     http.MethodPost,
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid path",
			url:        ts.URL + "/invalid/path",
			method:     http.MethodPost,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "wrong method",
			url:        ts.URL + "/update/gauge/test_metric/123.45",
			method:     http.MethodGet,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}