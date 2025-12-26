// Package handler содержит HTTP-обработчики сервиса метрик и агента.
package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Skywardkite/service-metrics/internal/audit"
	"github.com/Skywardkite/service-metrics/internal/config/server_config"
	"github.com/Skywardkite/service-metrics/internal/service"
)

// Handler объединяет HTTP-обработчики сервиса.
type Handler struct {
	service service.MetricServiceInterface
	cfg     *server_config.Config
	logger  *zap.SugaredLogger
	audit   audit.AuditPublisherInterface
}

func NewHandler(s service.MetricServiceInterface, cfg *server_config.Config, logger *zap.SugaredLogger, audit audit.AuditPublisherInterface) *Handler {
	return &Handler{
		service: s,
		cfg:     cfg,
		logger:  logger,
		audit:   audit,
	}
}

// UpdateHandler обновляет метрику в хранилище.
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

	ip := clientIP(req)

	// Отправляем аудит
	h.audit.Publish(audit.AuditEvent{
		TS:        time.Now().Unix(),
		Metrics:   []string{metricName},
		IPAddress: ip,
	})

	// Собираем ответ
	currentTime := time.Now().UTC().Format(time.RFC1123)
	res.Header().Set("Date", currentTime)
	responseBody := ""
	res.Header().Set("Content-Length", fmt.Sprintf("%d", len(responseBody)))
	res.Header().Set("Content-Type", "application/json")

	if h.cfg.Key != "" {
		hash := SignBody([]byte(responseBody), h.cfg.Key)
		res.Header().Set("HashSHA256", hash)
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(responseBody))
}
