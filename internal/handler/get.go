package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if h.service.Cfg.Key != "" {
		success, err := checkKey(req, h.service.Cfg.Key)
		if err != nil {
			h.logger.Errorf("Error read body", err)
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if !success {
			h.logger.Info("no rights")
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

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

	if h.service.Cfg.Key != "" {
        hash := signBody([]byte(value), h.service.Cfg.Key)
        res.Header().Set("HashSHA256", hash)
    }
	
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(value))
}