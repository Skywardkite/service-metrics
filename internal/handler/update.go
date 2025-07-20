package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Skywardkite/service-metrics/internal/service"
)

type Handler struct {
	service *service.MetricService
}

func NewHandler(s *service.MetricService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(parts) != 4 || parts[0] != "update" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	metricType := parts[1]
	metricName := parts[2]
	metricValue := parts[3]
		
	if metricName == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	err := h.service.UpdateMetric(metricType, metricName, metricValue)
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