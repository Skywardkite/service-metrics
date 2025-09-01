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
		h.logger.Errorf("error reading request body", "err", err)
        http.Error(res, err.Error(), http.StatusBadRequest)
        return
    }
    if err := json.Unmarshal(buf.Bytes(), &metric); err != nil {
		h.logger.Errorf("error unmarshalling request body", "err", err)
        http.Error(res, err.Error(), http.StatusBadRequest)
        return
    }

	if metric.ID == "" {
		h.logger.Info("metric ID is empty")
		http.Error(res, "Invalid metric name", http.StatusNotFound)
		return
	}

	var value string
	switch metric.MType {
	case model.Gauge:
		if metric.Value == nil {
			h.logger.Info("invalid gauge value")
			http.Error(res, "invalid gauge value", http.StatusBadRequest)
			return
		}
		value = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	case model.Counter:
		if metric.Delta == nil {
			h.logger.Info("invalid counter value")
			http.Error(res, "invalid counter value", http.StatusBadRequest)
			return
		}
		value = strconv.FormatInt(*metric.Delta, 10)
	}

	if err := h.service.UpdateMetric(ctx, metric.MType, metric.ID, value); err != nil {
		h.logger.Errorf("error updating metric", "err", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	r, err := json.Marshal(metric)
    if err != nil {
		h.logger.Errorf("error marshalling response", "err", err)
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    }

	res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusOK)
    res.Write(r)
}