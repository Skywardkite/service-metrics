package main

import (
	"log"
	"net/http"

	serverconfig "github.com/Skywardkite/service-metrics/internal/config/server_config"
	"github.com/Skywardkite/service-metrics/internal/database"
	"github.com/Skywardkite/service-metrics/internal/filestorage"
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
	fileStorage := filestorage.NewStorageConfig(&cfg, store)
	fileStorage.Run()

	var db *database.DB
	if cfg.DatabaseDSN != "" {
		db, err = database.New(cfg.DatabaseDSN)
		if err != nil {
			logger.Sugar.Fatalw("Failed to connect to database", "error", err)
		}
		defer db.Close()
	}
	
    metricService := service.NewMetricService(&cfg, store)
    h := handler.NewHandler(metricService, db)

    r := chi.NewRouter()
	// Применяем middleware ко всем роутам
	r.Use(logger.WithLogging)
	r.Use(gzipMiddleware)

	// Регистрируем обработчики
    r.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateHandler)
	r.Get("/value/{metricType}/{metricName}", h.GetHandler)
	r.Get("/", h.GetAllMetricsHandler)
	r.Get("/ping", h.PingHandler)

	r.Post("/update/", h.UpdateJSONHandler)
	r.Post("/value/", h.GetMetricJSONHandler)
   	if err := http.ListenAndServe(cfg.FlagRunAddr, r); err != nil {
		logger.Sugar.Fatalw("Error to listen server", err.Error(), "event", "start server")
	}
}