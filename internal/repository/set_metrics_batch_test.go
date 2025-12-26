package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	model "github.com/Skywardkite/service-metrics/internal/model"
	pointers "github.com/Skywardkite/service-metrics/internal/utils"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestPostgresStorage_SetMetricsBatch(t *testing.T) {
	tests := []struct {
		name    string
		metrics []model.Metrics
		setup   func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success with gauge and counter",
			metrics: []model.Metrics{
				{ID: "Alloc", MType: "gauge", Value: pointers.To(float64(123.45))},
				{ID: "PollCount", MType: "counter", Delta: pointers.To(int64(10))},
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`
                    INSERT INTO gauges (name, value) 
                    VALUES ($1, $2)
                    ON CONFLICT (name) DO UPDATE SET value = $2
                `)).
					WithArgs("Alloc", 123.45).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(regexp.QuoteMeta(`
                    INSERT INTO counters (name, value) 
                    VALUES ($1, $2)
                    ON CONFLICT (name) DO UPDATE SET value = counters.value + $2
                `)).
					WithArgs("PollCount", int64(10)).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "exec error rollback",
			metrics: []model.Metrics{
				{ID: "Alloc", MType: "gauge", Value: pointers.To(float64(123.45))},
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`
                    INSERT INTO gauges (name, value) 
                    VALUES ($1, $2)
                    ON CONFLICT (name) DO UPDATE SET value = $2
                `)).
					WithArgs("Alloc", 123.45).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			storage := &PostgresStorage{db: sqlxDB}

			tt.setup(mock)

			err = storage.SetMetricsBatch(context.Background(), tt.metrics)
			if tt.wantErr {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
			}

			assert.NilError(t, mock.ExpectationsWereMet())
		})
	}
}
