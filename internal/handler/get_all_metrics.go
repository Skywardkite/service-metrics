package handler

import (
	"html/template"
	"net/http"
)

type MetricsPageData struct {
	Gauges    map[string]float64
	Counters  map[string]int64
}

func (h *Handler) GetAllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

    if h.service.Cfg.Key != "" {
		success, err := checkKey(req, h.service.Cfg.Key)
		if err != nil {
			h.logger.Errorf("Error read body", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if !success {
            h.logger.Info("no rights")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}

    gauges, counters, err := h.service.GetAllMetrics(ctx)
    if err != nil {
        h.logger.Errorw("Failed to get metrics", "error", err)
        res.WriteHeader(http.StatusInternalServerError)
        return
    }

    data := MetricsPageData{
		Gauges:    gauges,
		Counters:  counters,
	}

    tmpl, err := template.ParseFiles("internal/templates/metrics.html")
    if err != nil {
        h.logger.Errorw("Failed to parse file", "error", err)
        res.WriteHeader(http.StatusInternalServerError)
        return
    }

    res.Header().Set("Content-Type", "text/html")

    if h.service.Cfg.Key != "" {
        responseBody := []byte("OK")
        hash := signBody(responseBody, h.service.Cfg.Key)
        res.Header().Set("HashSHA256", hash)
        res.Write(responseBody)
    }

    res.WriteHeader(http.StatusOK)

    if err := tmpl.Execute(res, data); err != nil {
        h.logger.Errorw("Failed to execute file", "error", err)
        res.WriteHeader(http.StatusInternalServerError)
        return
	}
}