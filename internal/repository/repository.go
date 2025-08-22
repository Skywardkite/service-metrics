package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	SetCounter(ctx context.Context, name string, value int64) error
	SetGauge(ctx context.Context, name string, value float64) error
	GetCounter(ctx context.Context, name string) (int64, error)
	GetGauge(ctx context.Context, name string) (float64, error)
	GetMetrics(ctx context.Context) (map[string]float64, map[string]int64, error)
	Ping() error
}

type PostgresStorage struct {
	db *sqlx.DB
}

func New(dsn string) (*PostgresStorage, error) {
	// Используем Connect вместо Open - он проверяет соединение
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Printf("Connection failed: %v", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = applyMigrations(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("PostgreSQL storage initialized successfully")
	return &PostgresStorage{db: db}, nil
}

func (r *PostgresStorage) Close() error {
	if r.db == nil {
		return nil
	}
	return r.db.Close()
}

func (r *PostgresStorage) Ping() error {
	return r.db.Ping()
}

func applyMigrations(dsn string) error {
    m, err := migrate.New(
        "file://migrations",
        dsn,
    )
    if err != nil {
        return fmt.Errorf("failed to create migrate instance: %w", err)
    }

    // Применяем все новые миграции
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("failed to apply migrations: %w", err)
    }

    log.Println("Migrations applied successfully")
    return nil
}