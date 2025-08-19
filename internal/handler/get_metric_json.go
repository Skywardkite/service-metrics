package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	model "github.com/Skywardkite/service-metrics/internal/model"
)

func (h *Handler) GetMetricJSONHandler(res http.ResponseWriter, req *http.Request) {
	var metric model.Metrics
    var buf bytes.Buffer
    if _, err := buf.ReadFrom(req.Body); err != nil {
        http.Error(res, err.Error(), http.StatusBadRequest)
        return
    }
    if err := json.Unmarshal(buf.Bytes(), &metric); err != nil {
        http.Error(res, err.Error(), http.StatusBadRequest)
        return
    }

	value, err := h.service.GetMetric(metric.MType, metric.ID)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

    // Получаем значение метрики из хранилища
    switch metric.MType {
    case model.Gauge:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(res, "invalid gauge value", http.StatusBadRequest)
			return
		}
        metric.Value = &floatValue
    case model.Counter:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(res, "invalid counter value", http.StatusBadRequest)
			return
		}
        metric.Delta = &intValue
    }

	r, err := json.Marshal(metric)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    }

	res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusOK)
    res.Write(r)
}