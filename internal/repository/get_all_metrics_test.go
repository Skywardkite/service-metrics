package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestPostgresStorage_GetMetrics(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(sqlmock.Sqlmock)
		wantGauges   map[string]float64
		wantCounters map[string]int64
		wantErr      error
	}{
		{
			name: "success",
			setup: func(mock sqlmock.Sqlmock) {
				gaugeRows := sqlmock.NewRows([]string{"name", "value"}).
					AddRow("Alloc", 123.45).
					AddRow("Heap", 678.9)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM gauges")).
					WillReturnRows(gaugeRows)

				counterRows := sqlmock.NewRows([]string{"name", "value"}).
					AddRow("PollCount", 10).
					AddRow("Errors", 2)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM counters")).
					WillReturnRows(counterRows)
			},
			wantGauges: map[string]float64{
				"Alloc": 123.45,
				"Heap":  678.9,
			},
			wantCounters: map[string]int64{
				"PollCount": 10,
				"Errors":    2,
			},
			wantErr: nil,
		},
		{
			name: "empty tables",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM gauges")).
					WillReturnRows(sqlmock.NewRows([]string{"name", "value"}))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM counters")).
					WillReturnRows(sqlmock.NewRows([]string{"name", "value"}))
			},
			wantGauges:   map[string]float64{},
			wantCounters: map[string]int64{},
			wantErr:      nil,
		},
		{
			name: "gauges db error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM gauges")).
					WillReturnError(errors.New("db error"))
			},
			wantGauges:   nil,
			wantCounters: nil,
			wantErr:      errors.New("db error"),
		},
		{
			name: "db counters error",
			setup: func(mock sqlmock.Sqlmock) {
				gaugeRows := sqlmock.NewRows([]string{"name", "value"}).
					AddRow("Alloc", 123.45).
					AddRow("Heap", 678.9)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM gauges")).
					WillReturnRows(gaugeRows)

				mock.ExpectQuery(regexp.QuoteMeta(queryAllGetCounter)).
					WillReturnError(errors.New("db error"))
			},
			wantCounters: nil,
			wantGauges:   nil,
			wantErr:      errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")

			tt.setup(mock)

			storage := &PostgresStorage{
				db: sqlxDB,
			}

			gotGauges, gotCounters, err := storage.GetMetrics(context.Background())
			if tt.wantErr != nil {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
			}

			assert.DeepEqual(t, tt.wantGauges, gotGauges)
			assert.DeepEqual(t, tt.wantCounters, gotCounters)
			assert.NilError(t, mock.ExpectationsWereMet())
		})
	}
}
