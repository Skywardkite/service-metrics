package handler

import (
	"crypto/hmac"
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/Skywardkite/service-metrics/internal/model"
)

func (h *Handler) UpdateMetricsBatchJSONHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

    // Декодируем JSON
    var metrics []model.Metrics
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&metrics); err != nil {
		h.logger.Errorf("Error decoding JSON: %v", err)
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    defer r.Body.Close()

    jsonData, _ := json.Marshal(metrics)
    if h.service.Cfg.Key != "" {
        expected := signBody(jsonData, h.service.Cfg.Key)
        received := r.Header.Get("HashSHA256")

        if received == "" || !hmac.Equal([]byte(received), []byte(expected)) {
            h.logger.Info("no rights")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
        }
	}

    // Проверяем, что батч не пустой
    if len(metrics) == 0 {
        h.logger.Errorf("Empty batch to update metrics")
        http.Error(w, "Empty batch", http.StatusBadRequest)
        return
    }

	// Сохраняем метрики в базе в одной транзакции
    err := h.service.SaveMetricsBatch(ctx, metrics)
    if err != nil {
        h.logger.Errorf("Error saving metrics: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    responseBody := []byte(`{"status":"ok"}`)

	w.Header().Set("Content-Type", "application/json")

    if h.service.Cfg.Key != "" {
        fmt.Println("error here")
        hash := signBody(responseBody, h.service.Cfg.Key)
        w.Header().Set("HashSHA256", hash)
    }

    w.WriteHeader(http.StatusOK)
    w.Write(responseBody)
}
