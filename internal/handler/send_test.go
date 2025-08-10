package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Skywardkite/service-metrics/internal/agent"
	"github.com/stretchr/testify/assert"
)

func Test_sendPlainPost(t *testing.T) {
    storage := &agent.AgentMetrics{
        Gauge:   map[string]float64{"TestMetric": 1.23},
        Counter: map[string]int64{"TestCounter": 1},
    }

    tests := []struct {
        name          string
        serverHandler http.HandlerFunc
        client        *http.Client
        port          string
        wantErr       bool
        errorContains string
    }{
        {
            name: "successful request",
            serverHandler: func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
            },
            client:    &http.Client{Timeout: 1 * time.Second},
            wantErr:   false,
        },
        {
            name: "server error",
            serverHandler: func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusInternalServerError)
            },
            client:        &http.Client{Timeout: 1 * time.Second},
            wantErr:      true,
            errorContains: "500",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            server := httptest.NewServer(tt.serverHandler)
            defer server.Close()

            // Используем URL тестового сервера вместо фиксированного порта
            err := SendMetrics(tt.client, storage, server.URL)

            if tt.wantErr {
                assert.Error(t, err)
                if tt.errorContains != "" {
                    assert.Contains(t, err.Error(), tt.errorContains)
                }
            } else {
                assert.NoError(t, err)
            }
        })
    }
}