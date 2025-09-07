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
        // Обработка gzip
        if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
            gz, err := gzip.NewReader(r.Body)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            defer gz.Close()
            r.Body = gz
        }

        // Обработка авторизации
        if key != "" {
            // читаем тело для проверки подписи
            bodyBytes, err := io.ReadAll(r.Body)
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

            // восстанавливаем тело для обработки gzip
            r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
        }

        // Создаем обертку для ResponseWriter для сжатия ответа
        writer := &gzipResponseWriter{
            ResponseWriter: w,
            request:        r,
        }
        defer writer.Close()

        next.ServeHTTP(writer, r)
    })
}