package handler

import (
	"encoding/json"
	"net/http"

	model "github.com/Skywardkite/service-metrics/internal/model"
)

func (h *Handler) UpdateMetricsBatchJSONHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

    if h.service.Cfg.Key != "" {
		success, err := checkKey(r, h.service.Cfg.Key)
		if err != nil {
			h.logger.Errorf("Error read body", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if !success {
            h.logger.Info("no rights")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

    // Декодируем JSON
    var metrics []model.Metrics
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&metrics); err != nil {
		h.logger.Errorf("Error decoding JSON: %v", err)
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    defer r.Body.Close()

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
        hash := signBody(responseBody, h.service.Cfg.Key)
        w.Header().Set("HashSHA256", hash)
    }

    w.WriteHeader(http.StatusOK)
    w.Write(responseBody)
}
