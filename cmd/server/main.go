package main

import (
	"log"
	"net/http"

	serverconfig "github.com/Skywardkite/service-metrics/internal/config/serverConfig"
	"github.com/Skywardkite/service-metrics/internal/handler"
	"github.com/Skywardkite/service-metrics/internal/service"
	"github.com/Skywardkite/service-metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := serverconfig.ParseFlags()
    if err != nil {
        log.Fatal("Error to parse flags:", err)
    }

    store := storage.NewMemStorage()
    metricService := service.NewMetricService(store)
    h := handler.NewHandler(metricService)

    r := chi.NewRouter()
    r.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateHandler)
	r.Get("/value/{metricType}/{metricName}", h.GetHandler)
	r.Get("/", h.GetAllMetricsHandler)
   	if err := http.ListenAndServe(cfg.FlagRunAddr, r); err != nil {
		log.Fatal("Error to listen server:", err)
	}
}