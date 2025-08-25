package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	model "github.com/Skywardkite/service-metrics/internal/model"
	"github.com/stretchr/testify/assert"
)

func Test_sendPlainPost(t *testing.T) {
    testValue := 45.979785
    testDelta := int64(6767884)

	tests := []struct {
		name string
		serverHandler  http.HandlerFunc
		client *http.Client
        metric model.Metrics
		url    string
		wantErr   bool
		errorMessage string
	}{
		{
			name: "successful request",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			client:      &http.Client{},
            metric:     model.Metrics{
                ID: "test",
                MType: model.Gauge,
                Value: &testValue,
            },
			url:         "/test",
			wantErr: false,
		},
		{
			name: "server error",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			client:        &http.Client{},
            metric:     model.Metrics{
                ID:     "test",
                MType:  model.Counter,
                Delta:  &testDelta,
            },
			url:           "/test",
			wantErr:  		true,
			errorMessage: "all retries failed, last error: server error: 500",
		},
		{
			name: "invalid URL",
			serverHandler: nil,
			client:        &http.Client{},
			url:           "://invalid-url",
			wantErr:   		true,
			errorMessage: "failed to create request: parse \"://invalid-url\": missing protocol scheme",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var server *httptest.Server
			if tt.serverHandler != nil {
				server = httptest.NewServer(tt.serverHandler)
				defer server.Close()
				tt.url = server.URL + tt.url
			}

			err := sendPlainPost(tt.client, tt.url, tt.metric)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}