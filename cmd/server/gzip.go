package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

var gzipPool = sync.Pool{
	New: func() any {
		w, _ := gzip.NewWriterLevel(io.Discard, gzip.BestSpeed)
		return w
	},
}

type gzipResponseWriter struct {
	http.ResponseWriter
	request     *http.Request
	gzipWriter  *gzip.Writer
	wroteHeader bool
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gz.Close()

			bodyBytes, err := io.ReadAll(gz)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Восстанавливаем тело, чтобы дальше хэндлеры могли читать
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		writer := &gzipResponseWriter{
			ResponseWriter: w,
			request:        r,
		}
		defer writer.Close()

		next.ServeHTTP(writer, r)
	})
}

func (w *gzipResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true

	// Проверяем, нужно ли сжимать
	contentType := w.Header().Get("Content-Type")
	acceptsGzip := strings.Contains(w.request.Header.Get("Accept-Encoding"), "gzip")

	if acceptsGzip && (contentType == "text/html" || contentType == "application/json") {
		w.Header().Set("Content-Encoding", "gzip")

		gz := gzipPool.Get().(*gzip.Writer)
		gz.Reset(w.ResponseWriter)
		w.gzipWriter = gz
	}

	w.ResponseWriter.WriteHeader(code)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	if w.gzipWriter != nil {
		return w.gzipWriter.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *gzipResponseWriter) Close() {
	if w.gzipWriter != nil {
		w.gzipWriter.Close()
		gzipPool.Put(w.gzipWriter)
	}
}
