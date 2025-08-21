package storage

import (
	serverconfig "github.com/Skywardkite/service-metrics/internal/config/server_config"
	"github.com/Skywardkite/service-metrics/internal/repository"
)

func NewStorage(cfg *serverconfig.Config) (repository.Storage, error) {
	if cfg.DatabaseDSN != "" {
		// Используем PostgreSQL
		return repository.New(cfg.DatabaseDSN)
	}
	
	// Используем память с файловой синхронизацией
	return NewMemStorage(), nil
}