package main

import (
	"net/http"

	"github.com/Skywardkite/service-metrics.git/internal/handler"
	"github.com/Skywardkite/service-metrics.git/internal/service"
	"github.com/Skywardkite/service-metrics.git/internal/storage"
)

func main() {
    store := storage.NewMemStorage()
    metricService := service.NewMetricService(store)
    h := handler.NewHandler(metricService)

    http.HandleFunc(`/update/`, h.UpdateHandler)

    err := http.ListenAndServe(`:8080`, nil)
    if err != nil {
        panic(err)
    }
}

