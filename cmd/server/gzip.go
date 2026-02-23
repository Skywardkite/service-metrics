package main

import (
	"bytes"
	"compress/gzip"
	"crypto/rsa"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/Skywardkite/service-metrics/internal/crypto"
)

// GzipPool пул для повторного использования gzip.Writer.
var GzipPool = sync.Pool{
	New: func() any {
		w, _ := gzip.NewWriterLevel(io.Discard, gzip.BestSpeed)
		return w
	},
}

// gzipResponseWriter оборачивает http.ResponseWriter и добавляет поддержку gzip сжатия для ответов.
type gzipResponseWriter struct {
	http.ResponseWriter
	request     *http.Request
	gzipWriter  *gzip.Writer
	wroteHeader bool
}

// GzipMiddleware функция-обертка для добавления middleware сжатия ответов.
func GzipMiddleware(privKey *rsa.PrivateKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var bodyBytes []byte
			var err error

			bodyBytes, err = io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()

			// Расшифровка, если есть ключ
			if r.Method == http.MethodPost && privKey != nil {
				bodyBytes, err = crypto.Decrypt(privKey, bodyBytes)
				if err != nil {
					http.Error(w, "decrypt error: "+err.Error(), http.StatusBadRequest)
					return
				}
			}

			// Разжатие gzip
			if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				gz, err := gzip.NewReader(bytes.NewReader(bodyBytes))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				bodyBytes, err = io.ReadAll(gz)
				gz.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			// Восстанавливаем тело, чтобы дальше хэндлеры могли читать
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

			next.ServeHTTP(w, r)
		})
	}
}

// WriteHeader реализует интерфейс http.ResponseWriter.
// Записывает статус ответа и устанавливает заголовки для сжатия, если это необходимо.
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

		gz := GzipPool.Get().(*gzip.Writer)
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
		GzipPool.Put(w.gzipWriter)
	}
}
