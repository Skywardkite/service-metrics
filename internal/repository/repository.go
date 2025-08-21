package repository

import (
	"context"
	"fmt"
	"log"

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