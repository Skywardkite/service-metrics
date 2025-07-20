package main

import (
	"net/http"

	"github.com/Skywardkite/service-metrics/internal/handler"
	"github.com/Skywardkite/service-metrics/internal/service"
	"github.com/Skywardkite/service-metrics/internal/storage"
)

func main() {
    store := storage.NewMetricsStorage()
    metricService := service.NewMetricService(store)
    h := handler.NewHandler(metricService)

    http.HandleFunc(`/update/`, h.UpdateHandler)

    err := http.ListenAndServe(`:8080`, nil)
    if err != nil {
        panic(err)
    }
}

