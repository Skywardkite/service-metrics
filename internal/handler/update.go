package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Skywardkite/service-metrics.git/internal/service"
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
		
		// При попытке передать запрос без имени метрики возвращать http.StatusNotFound
		if metricName == "" {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		err := h.service.UpdateMetric(metricType, metricName, metricValue)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// Собираем ответа
		// Устанавливаем статус ответа
		res.WriteHeader(http.StatusOK)

		// Устанавливаем заголовок "Date"
		currentTime := time.Now().UTC().Format(time.RFC1123)
		res.Header().Set("Date", currentTime)

		// Устанавливаем заголовок "Content-Length"
		responseBody := ""
		res.Header().Set("Content-Length", fmt.Sprintf("%d", len(responseBody)))

		// Устанавливаем заголовок "Content-Type"
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")

		// Отправляем тело ответа
		res.Write([]byte(responseBody))

}