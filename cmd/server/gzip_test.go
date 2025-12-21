package server

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestGzipMiddleware(t *testing.T) {
	// хэндлер, который просто читает тело и возвращает длину
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NilError(t, err)
		w.Write([]byte(fmt.Sprintf("%d", len(body))))
	})

	tests := []struct {
		name          string
		body          []byte
		gzipEncoding  bool
		wantStatus    int
		wantRespBody  string
		corruptGzip   bool
	}{
		{
			name:         "no gzip",
			body:         []byte("hello world"),
			gzipEncoding: false,
			wantStatus:   http.StatusOK,
			wantRespBody: "11",
		},
		{
			name:         "gzip encoded",
			body:         []byte("hello gzip"),
			gzipEncoding: true,
			wantStatus:   http.StatusOK,
			wantRespBody: "10",
		},
		{
			name:         "corrupt gzip",
			body:         []byte("bad data"),
			gzipEncoding: true,
			corruptGzip:  true,
			wantStatus:   http.StatusInternalServerError,
			wantRespBody: "gzip: invalid header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyReader io.Reader

			if tt.gzipEncoding {
				var buf bytes.Buffer
				if tt.corruptGzip {
					// пишем некорректный gzip
					buf.Write([]byte("not a gzip"))
				} else {
					gw := gzip.NewWriter(&buf)
					_, err := gw.Write(tt.body)
					assert.NilError(t, err)
					gw.Close()
				}
				bodyReader = &buf
			} else {
				bodyReader = bytes.NewReader(tt.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/", bodyReader)
			if tt.gzipEncoding {
				req.Header.Set("Content-Encoding", "gzip")
			}

			w := httptest.NewRecorder()
			GzipMiddleware(nextHandler).ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()
			respBody, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
			if tt.gzipEncoding && tt.corruptGzip {
				assert.Assert(t, strings.Contains(string(respBody), "invalid header"))
			} else {
				assert.Equal(t, tt.wantRespBody, string(respBody))
			}
		})
	}
}
