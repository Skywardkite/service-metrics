package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	model "github.com/Skywardkite/service-metrics/internal/model"
)

func (h *Handler) UpdateJSONHandler(res http.ResponseWriter, req *http.Request) {
	var metric model.Metrics
    var buf bytes.Buffer
	
    _, err := buf.ReadFrom(req.Body)
    if err != nil {
        http.Error(res, err.Error(), http.StatusBadRequest)
        return
    }
    if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
        http.Error(res, err.Error(), http.StatusBadRequest)
        return
    }

	if metric.ID == "" {
		http.Error(res, "Invalid metric name", http.StatusNotFound)
		return
	}

	var value string
	switch metric.MType {
	case "gauge":
		if metric.Value == nil {
			http.Error(res, "invalid gauge value", http.StatusBadRequest)
			return
		}
		value = fmt.Sprintf("%v", *metric.Value)
	case "counter":
		if metric.Delta == nil {
			http.Error(res, "invalid counter value", http.StatusBadRequest)
			return
		}
		value = strconv.FormatInt(*metric.Delta, 10)
	}

	err = h.service.UpdateMetric(metric.MType, metric.ID, value)
	if err != nil {
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