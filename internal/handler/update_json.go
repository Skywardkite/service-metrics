package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	model "github.com/Skywardkite/service-metrics/internal/model"
)

func (h *Handler) UpdateJSONHandler(res http.ResponseWriter, req *http.Request) {
	var metric model.Metrics
    var buf bytes.Buffer
	ctx := req.Context()
	
    if _, err := buf.ReadFrom(req.Body); err != nil {
        http.Error(res, err.Error(), http.StatusBadRequest)
        return
    }
    if err := json.Unmarshal(buf.Bytes(), &metric); err != nil {
        http.Error(res, err.Error(), http.StatusBadRequest)
        return
    }

	if metric.ID == "" {
		http.Error(res, "Invalid metric name", http.StatusNotFound)
		return
	}

	var value string
	switch metric.MType {
	case model.Gauge:
		if metric.Value == nil {
			http.Error(res, "invalid gauge value", http.StatusBadRequest)
			return
		}
		value = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	case model.Counter:
		if metric.Delta == nil {
			http.Error(res, "invalid counter value", http.StatusBadRequest)
			return
		}
		value = strconv.FormatInt(*metric.Delta, 10)
	}

	if err := h.service.UpdateMetric(ctx, metric.MType, metric.ID, value); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
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