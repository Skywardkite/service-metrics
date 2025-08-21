package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Skywardkite/service-metrics/internal/repository"
	"github.com/Skywardkite/service-metrics/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *service.MetricService
	store repository.Storage
}

func NewHandler(s *service.MetricService, store repository.Storage) *Handler {
	return &Handler{service: s, store: store}
}

func (h *Handler) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	metricValue := chi.URLParam(req, "metricValue")
	if metricName == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	err := h.service.UpdateMetric(ctx, metricType, metricName, metricValue)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// Собираем ответ
	res.WriteHeader(http.StatusOK)
	currentTime := time.Now().UTC().Format(time.RFC1123)
	res.Header().Set("Date", currentTime)
	responseBody := ""
	res.Header().Set("Content-Length", fmt.Sprintf("%d", len(responseBody)))
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	res.Write([]byte(responseBody))
}