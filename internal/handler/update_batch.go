package handler

import (
	"encoding/json"
	"net/http"

	model "github.com/Skywardkite/service-metrics/internal/model"
)

func (h *Handler) UpdateMetricsBatchJSONHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

    // Декодируем JSON
    var metrics []model.Metrics
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&metrics); err != nil {
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    defer r.Body.Close()

    // Проверяем, что батч не пустой
    if len(metrics) == 0 {
        http.Error(w, "Empty batch", http.StatusBadRequest)
        return
    }

	// Сохраняем метрики в базе в одной транзакции
    err := h.service.SaveMetricsBatch(ctx, metrics)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
}
