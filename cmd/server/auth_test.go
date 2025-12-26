package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Skywardkite/service-metrics/internal/handler"
	"gotest.tools/assert"
)

func TestAuthMiddleware(t *testing.T) {
	key := "secret-key"

	// простой handler, который записывает "ok"
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	tests := []struct {
		name       string
		key        string
		body       []byte
		hashHeader string
		wantStatus int
	}{
		{
			name:       "empty key, passes through",
			key:        "",
			body:       []byte("body"),
			hashHeader: "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "empty hash header, passes through",
			key:        key,
			body:       []byte("body"),
			hashHeader: "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid hash",
			key:        key,
			body:       []byte("body"),
			hashHeader: "wronghash",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "valid hash",
			key:        key,
			body:       []byte("body"),
			hashHeader: handler.SignBody([]byte("body"), key),
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tt.body))
			if tt.hashHeader != "" {
				req.Header.Set("HashSHA256", tt.hashHeader)
			}

			w := httptest.NewRecorder()
			mw := AuthMiddleware(tt.key, nextHandler)
			mw.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}
