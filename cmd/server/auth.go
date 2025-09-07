package main

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"io"
	"net/http"
	"strings"

	"github.com/Skywardkite/service-metrics/internal/handler"
)

func authAndGzipMiddleware(key string, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if key != "" {
            // Буферизуем оригинальное тело
            var bodyBuf bytes.Buffer
            tee := io.TeeReader(r.Body, &bodyBuf)
            
            // Читаем для проверки подписи
            bodyBytes, err := io.ReadAll(tee)
            if err != nil {
                http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
                return
            }

            expected := handler.SignBody(bodyBytes, key)
            received := r.Header.Get("HashSHA256")
            if received == "" || !hmac.Equal([]byte(received), []byte(expected)) {
                http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
                return
            }

            // Восстанавливаем оригинальное тело из буфера
            r.Body = io.NopCloser(&bodyBuf)
        }

        if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
            gz, err := gzip.NewReader(r.Body)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            defer gz.Close()
            r.Body = gz
        }

        writer := &gzipResponseWriter{
            ResponseWriter: w,
            request:        r,
        }
        defer writer.Close()

        next.ServeHTTP(writer, r)
    })
}