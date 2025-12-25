package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Skywardkite/service-metrics/internal/audit"
	model "github.com/Skywardkite/service-metrics/internal/model"
)

// UpdateJSONHandler - обновляет одну метрику. Принимает значения в json.
func (h *Handler) UpdateJSONHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var metric model.Metrics
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(req.Body); err != nil {
		h.logger.Errorf("error reading request body", "err", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &metric); err != nil {
		h.logger.Errorf("error unmarshalling request body", "err", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if metric.ID == "" {
		h.logger.Info("metric ID is empty")
		http.Error(res, "Invalid metric name", http.StatusNotFound)
		return
	}

	var value string
	switch metric.MType {
	case model.Gauge:
		if metric.Value == nil {
			h.logger.Info("invalid gauge value")
			http.Error(res, "invalid gauge value", http.StatusBadRequest)
			return
		}
		value = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	case model.Counter:
		if metric.Delta == nil {
			h.logger.Info("invalid counter value")
			http.Error(res, "invalid counter value", http.StatusBadRequest)
			return
		}
		value = strconv.FormatInt(*metric.Delta, 10)
	}

	if err := h.service.UpdateMetric(ctx, metric.MType, metric.ID, value); err != nil {
		h.logger.Errorf("error updating metric", "err", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// Отправка события
	ip := clientIP(req)
	h.audit.Publish(audit.AuditEvent{
		TS:        time.Now().Unix(),
		Metrics:   []string{metric.ID},
		IPAddress: ip,
	})

	r, err := json.Marshal(metric)
	if err != nil {
		h.logger.Errorf("error marshalling response", "err", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if h.cfg.Key != "" {
		hash := SignBody(r, h.cfg.Key)
		res.Header().Set("HashSHA256", hash)
	}

	res.WriteHeader(http.StatusOK)
	res.Write(r)
}
