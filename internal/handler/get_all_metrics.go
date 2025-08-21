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
    gauges, counters, err := h.service.GetAllMetrics(ctx)
    if err != nil {
        res.WriteHeader(http.StatusInternalServerError)
        return
    }

    data := MetricsPageData{
		Gauges:    gauges,
		Counters:  counters,
	}

    tmpl, err := template.ParseFiles("internal/templates/metrics.html")
    if err != nil {
        res.WriteHeader(http.StatusInternalServerError)
        return
    }

    res.Header().Set("Content-Type", "text/html")
    res.WriteHeader(http.StatusOK)

    if err := tmpl.Execute(res, data); err != nil {
        res.WriteHeader(http.StatusInternalServerError)
        return
	}
}