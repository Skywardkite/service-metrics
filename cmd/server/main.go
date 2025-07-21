package main

import (
	"net/http"

	"github.com/Skywardkite/service-metrics/internal/handler"
	"github.com/Skywardkite/service-metrics/internal/service"
	"github.com/Skywardkite/service-metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	if err := parseFlagsServer(); err != nil {
        panic(err)
    }

    store := storage.NewMemStorage()
    metricService := service.NewMetricService(store)
    h := handler.NewHandler(metricService)

    r := chi.NewRouter()
    r.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateHandler)
	r.Get("/value/{metricType}/{metricName}", h.GetHandler)
	r.Get("/", h.GetAllMetricsHandler)

   	if err := http.ListenAndServe(flagRunAddr, r); err != nil {
		panic(err)
	}
}