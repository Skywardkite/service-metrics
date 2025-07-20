package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) GetAllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	gauges, counters := h.service.GetAllMetrics()

    res.Header().Set("Content-Type", "text/html")
    res.WriteHeader(http.StatusOK)

    fmt.Fprintln(res, "<!DOCTYPE html>")
    fmt.Fprintln(res, "<html><head><title>Metrics</title></head><body>")
    fmt.Fprintln(res, "<h1>Current Metrics</h1>")
    fmt.Fprintln(res, "<h2>Gauges</h2><ul>")
    for name, value := range gauges {
        fmt.Fprintf(res, "<li>%s: %v</li>\n", name, value)
    }
    fmt.Fprintln(res, "</ul>")
    fmt.Fprintln(res, "<h2>Counters</h2><ul>")
    for name, value := range counters {
        fmt.Fprintf(res, "<li>%s: %v</li>\n", name, value)
    }
    fmt.Fprintln(res, "</ul>")
    fmt.Fprintln(res, "</body></html>")
}