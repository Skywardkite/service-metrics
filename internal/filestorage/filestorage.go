package filestorage

import (
	"time"

	serverconfig "github.com/Skywardkite/service-metrics/internal/config/serverConfig"
	"github.com/Skywardkite/service-metrics/internal/logger"
	"github.com/Skywardkite/service-metrics/internal/storage"
)

type StorageConfig struct {
	cfg *serverconfig.Config
	store *storage.MemStorage
}

func NewStorageConfig(cfg *serverconfig.Config, store *storage.MemStorage) *StorageConfig {
	return &StorageConfig{
		cfg: cfg,
		store: store,
	}
}

func (c *StorageConfig) Run(){
	if c.cfg.Restore {
		if gauges, counters, err := storage.LoadMetrics(c.cfg.FileStoragePath); err == nil {
			for name, value := range gauges {
				c.store.SetGauge(name, value)
			}
			for name, delta := range counters {
				c.store.SetCounter(name, delta)
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
				gauges, counters := c.store.GetMetrics()
				err := storage.SaveMetrics(c.cfg.FileStoragePath, gauges, counters)
				if err != nil {
					logger.Sugar.Errorw("Failed to save metrics", "error", err)
					return
				}
			}
		}()
	}
}