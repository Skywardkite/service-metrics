package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	model "github.com/Skywardkite/service-metrics/internal/model"
)

func (h *Handler) GetJSONHandler(res http.ResponseWriter, req *http.Request) {
	var metric model.Metrics
    if err := json.NewDecoder(req.Body).Decode(&metric); err != nil {
        http.Error(res, "Invalid JSON format", http.StatusBadRequest)
        return
    }

	if metric.ID == "" {
		http.Error(res, "Invalid metric name", http.StatusNotFound)
        return
	}

    responseMetric := model.Metrics{
        ID:    metric.ID,
        MType: metric.MType,
    }

	value, err := h.service.GetMetric(metric.MType, metric.ID)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

    // Получаем значение метрики из хранилища
    switch metric.MType {
    case "gauge":
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(res, "invalid gauge value", http.StatusBadRequest)
			return
		}
        responseMetric.Value = &floatValue
    case "counter":
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(res, "invalid counter value", http.StatusBadRequest)
			return
		}
        responseMetric.Delta = &intValue
    }

    // Возвращаем метрику с заполненными значениями
    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusOK)
    json.NewEncoder(res).Encode(responseMetric)
}