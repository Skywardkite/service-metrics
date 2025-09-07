package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetHandler(res http.ResponseWriter, req *http.Request) {
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

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
	if h.service.Cfg.Key != "" {
        hash := SignBody([]byte(value), h.service.Cfg.Key)
        res.Header().Set("HashSHA256", hash)
    }
	res.Write([]byte(value))
}