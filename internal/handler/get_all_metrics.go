package handler

import (
	"html/template"
	"net/http"
)

// MetricsPageData собирает все метрики, чтобы передать их в шаблон metrics.html.
type MetricsPageData struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

// GetAllMetricsHandler отдаем все метрики.
// Подставляем все значения в шаблон metrics.html.
func (h *Handler) GetAllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	gauges, counters, err := h.service.GetAllMetrics(ctx)
	if err != nil {
		h.logger.Errorw("Failed to get metrics", "error", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := MetricsPageData{
		Gauges:   gauges,
		Counters: counters,
	}

	tmpl, err := template.ParseFiles("internal/templates/metrics.html")
	if err != nil {
		h.logger.Errorw("Failed to parse file", "error", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "text/html")

	if err := tmpl.Execute(res, data); err != nil {
		h.logger.Errorw("Failed to execute file", "error", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody := []byte("OK")
	if h.cfg.Key != "" {
		hash := SignBody(responseBody, h.cfg.Key)
		res.Header().Set("HashSHA256", hash)
	}

	res.WriteHeader(http.StatusOK)
	res.Write(responseBody)
}
