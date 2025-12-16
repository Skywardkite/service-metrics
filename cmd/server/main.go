package main

import (
	"context"
	"log"
	"net/http"

	"net/http/pprof"

	"github.com/Skywardkite/service-metrics/internal/audit"
	"github.com/Skywardkite/service-metrics/internal/config/server_config"
	"github.com/Skywardkite/service-metrics/internal/filestorage"
	"github.com/Skywardkite/service-metrics/internal/handler"
	logger "github.com/Skywardkite/service-metrics/internal/logger"
	"github.com/Skywardkite/service-metrics/internal/repository"
	"github.com/Skywardkite/service-metrics/internal/service"
	"github.com/Skywardkite/service-metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

func mountPprof(r chi.Router) {
	r.HandleFunc("/", pprof.Index)
	r.HandleFunc("/cmdline", pprof.Cmdline)
	r.HandleFunc("/profile", pprof.Profile)
	r.HandleFunc("/symbol", pprof.Symbol)
	r.HandleFunc("/trace", pprof.Trace)
	r.HandleFunc("/allocs", pprof.Handler("allocs").ServeHTTP)
	r.HandleFunc("/block", pprof.Handler("block").ServeHTTP)
	r.HandleFunc("/goroutine", pprof.Handler("goroutine").ServeHTTP)
	r.HandleFunc("/heap", pprof.Handler("heap").ServeHTTP)
	r.HandleFunc("/mutex", pprof.Handler("mutex").ServeHTTP)
	r.HandleFunc("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
}

func main() {
	if err := logger.Initialize(); err != nil {
		log.Fatal("Error to initialize logger:", err)
	}
	defer logger.Sync()

	cfg, err := server_config.ParseFlags()
	if err != nil {
		logger.Sugar.Fatalw("Error to parse flags", "error", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var store repository.Storage

	if cfg.DatabaseDSN != "" {
		pgStore, err := repository.New(cfg.DatabaseDSN)
		if err != nil {
			logger.Sugar.Fatalw("Failed to connect to database", "error", err)
		}
		store = pgStore

		// Закрываем соединение при завершении программы
		defer func() {
			if err := pgStore.Close(); err != nil {
				logger.Sugar.Errorw("Failed to close database", "error", err)
			}
		}()
	} else {
		store = storage.NewMemStorage()
		fileStorage := filestorage.NewStorageConfig(&cfg, store)
		fileStorage.Run(ctx)
	}

	publisher := audit.NewAuditPublisher()

	if cfg.AuditFile != "" {
		fileObs := &audit.FileObserver{FilePath: cfg.AuditFile}
		publisher.Subscribe(fileObs)
	}

	if cfg.AuditURL != "" {
		httpObs := &audit.HttpObserver{URL: cfg.AuditURL}
		publisher.Subscribe(httpObs)
	}

	metricService := service.NewMetricService(&cfg, store)
	h := handler.NewHandler(metricService, store, logger.Sugar, publisher)

	r := chi.NewRouter()

	// Применяем middleware ко всем роутам
	r.Use(logger.WithLogging)
	r.Use(func(next http.Handler) http.Handler {
		return authMiddleware(cfg.Key, next)
	})
	r.Use(gzipMiddleware)

	// pprof на /debug/pprof
	r.Route("/debug/pprof", func(pp chi.Router) {
		mountPprof(pp)
	})

	// Регистрируем обработчики
	r.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateHandler)
	r.Get("/value/{metricType}/{metricName}", h.GetHandler)
	r.Get("/", h.GetAllMetricsHandler)
	r.Get("/ping", h.PingHandler)

	r.Post("/update/", h.UpdateJSONHandler)
	r.Post("/updates/", h.UpdateMetricsBatchJSONHandler)
	r.Post("/value/", h.GetMetricJSONHandler)

	if err := http.ListenAndServe(cfg.FlagRunAddr, r); err != nil {
		logger.Sugar.Fatalw("Error to listen server", err.Error(), "event", "start server")
	}
}
