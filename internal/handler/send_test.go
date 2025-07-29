package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sendPlainPost(t *testing.T) {
	tests := []struct {
		name string
		serverHandler  http.HandlerFunc
		client *http.Client
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
			url:         "/test",
			wantErr: false,
		},
		{
			name: "server error",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			client:        &http.Client{},
			url:           "/test",
			wantErr:  		true,
			errorMessage: fmt.Sprintf("non-OK response status: %d", http.StatusInternalServerError),
		},
		{
			name: "invalid URL",
			serverHandler: nil,
			client:        &http.Client{},
			url:           "://invalid-url",
			wantErr:   		true,
			errorMessage: "failed to create request",
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

			err := sendPlainPost(tt.client, tt.url)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}