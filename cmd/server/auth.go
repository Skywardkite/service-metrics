package main

import (
	"bytes"
	"crypto/hmac"
	"io"
	"net/http"

	"github.com/Skywardkite/service-metrics/internal/handler"
)

func authMiddleware(key string, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if key == "" {
            next.ServeHTTP(w, r)
            return
        }

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

        // Восстанавливаем оригинальное тело из буфера
        r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

        next.ServeHTTP(w, r)
    })
}