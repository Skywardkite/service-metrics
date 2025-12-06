package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Skywardkite/service-metrics/internal/audit"
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

	// Отправка события
	ip := clientIP(r)

	var metricNames []string
	for _, m := range metrics {
		metricNames = append(metricNames, m.ID)
	}

	h.audit.Publish(audit.AuditEvent{
		TS:        time.Now().Unix(),
		Metrics:   metricNames,
		IPAddress: ip,
	})

	w.Header().Set("Content-Type", "application/json")

	responseBody := []byte("{}")
	if h.service.Cfg.Key != "" {
		hash := SignBody(responseBody, h.service.Cfg.Key)
		w.Header().Set("HashSHA256", hash)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
