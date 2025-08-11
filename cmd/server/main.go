package main

import (
	"log"
	"net/http"

	serverconfig "github.com/Skywardkite/service-metrics/internal/config/serverConfig"
	"github.com/Skywardkite/service-metrics/internal/handler"
	logger "github.com/Skywardkite/service-metrics/internal/logger"
	"github.com/Skywardkite/service-metrics/internal/service"
	"github.com/Skywardkite/service-metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	if err := logger.Initialize(); err != nil {
        log.Fatal("Error to initialize logger:", err)
    }
	defer logger.Sync()

	cfg, err := serverconfig.ParseFlags()
    if err != nil {
		logger.Sugar.Fatalw("Error to parse flags", "error", err)
    }

    store := storage.NewMemStorage()
    metricService := service.NewMetricService(store)
    h := handler.NewHandler(metricService)

    r := chi.NewRouter()
	// Применяем middleware логирования ко всем роутам
	r.Use(logger.WithLogging)

	// Регистрируем обработчики
    r.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateHandler)
	r.Get("/value/{metricType}/{metricName}", h.GetHandler)
	r.Get("/", h.GetAllMetricsHandler)

	//r.Post("/update", h.UpdateJSONHandler)
   	if err := http.ListenAndServe(cfg.FlagRunAddr, r); err != nil {
		logger.Sugar.Fatalw("Error to listen server", err.Error(), "event", "start server")
	}
}