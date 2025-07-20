package main

import (
	"net/http"

	"github.com/Skywardkite/service-metrics/internal/handler"
	"github.com/Skywardkite/service-metrics/internal/service"
	"github.com/Skywardkite/service-metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
    store := storage.NewMemStorage()
    metricService := service.NewMetricService(store)
    h := handler.NewHandler(metricService)

    r := chi.NewRouter()

    r.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateHandler)
	r.Get("/value/{metricType}/{metricName}", h.GetHandler)
	r.Get("/", h.GetAllMetricsHandler)

   if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}