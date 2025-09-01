package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Skywardkite/service-metrics/internal/repository"
	"github.com/Skywardkite/service-metrics/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	service *service.MetricService
	store repository.Storage
	logger *zap.SugaredLogger
}

func NewHandler(s *service.MetricService, store repository.Storage, logger *zap.SugaredLogger) *Handler {
	return &Handler{service: s, store: store, logger: logger}
}

func (h *Handler) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	metricValue := chi.URLParam(req, "metricValue")
	if metricName == "" {
		h.logger.Warn("Error updating metric, metricName is empty")
		res.WriteHeader(http.StatusNotFound)
		return
	}

	err := h.service.UpdateMetric(ctx, metricType, metricName, metricValue)
	if err != nil {
		h.logger.Errorf("Error updating metric", err)
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