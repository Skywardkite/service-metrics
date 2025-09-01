package filestorage

import (
	"context"
	"time"

	"github.com/Skywardkite/service-metrics/internal/config/server_config"
	"github.com/Skywardkite/service-metrics/internal/logger"
	"github.com/Skywardkite/service-metrics/internal/repository"
	"github.com/Skywardkite/service-metrics/internal/storage"
)

type StorageConfig struct {
	cfg *server_config.Config
	store repository.Storage
}

func NewStorageConfig(cfg *server_config.Config, store repository.Storage) *StorageConfig {
	return &StorageConfig{
		cfg: cfg,
		store: store,
	}
}

func (c *StorageConfig) Run(ctx context.Context){
	if c.cfg.Restore {
		if gauges, counters, err := storage.LoadMetrics(c.cfg.FileStoragePath); err == nil {
			for name, value := range gauges {
				c.store.SetGauge(ctx, name, value)
			}
			for name, delta := range counters {
				c.store.SetCounter(ctx, name, delta)
			}
			logger.Sugar.Info("Metrics restored from file")
		} else {
			logger.Sugar.Errorw("Failed to restore metrics", "error", err)
		}
	}

	if c.cfg.StoreInternal > 0 {
		go func() {
			ticker := time.NewTicker(c.cfg.StoreInternal)
			defer ticker.Stop()

			for range ticker.C {
				gauges, counters, err := c.store.GetMetrics(ctx)
				if err != nil {
					logger.Sugar.Errorf("Failed to get metrics filestorage", "error", err)
					return
				}

				err = storage.SaveMetrics(c.cfg.FileStoragePath, gauges, counters)
				if err != nil {
					logger.Sugar.Errorw("Failed to save metrics", "error", err)
					return
				}
			}
		}()
	}
}