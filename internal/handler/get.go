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
	res.Write([]byte(value))
}