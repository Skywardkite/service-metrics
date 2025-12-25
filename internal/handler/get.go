package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetMetric - вернет значение по типу метрики и ее названию.
func (h *Handler) GetMetric(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")

	if metricName == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	value, err := h.service.GetMetric(ctx, metricType, metricName)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if h.cfg.Key != "" {
		hash := SignBody([]byte(value), h.cfg.Key)
		res.Header().Set("HashSHA256", hash)
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(value))
}
