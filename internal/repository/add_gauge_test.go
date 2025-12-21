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

func TestPostgresStorage_SetGauge(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(querySetGauge)).
					WithArgs("Alloc", float64(5.4)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "db error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(querySetGauge)).
					WithArgs("Alloc", float64(5.4)).
					WillReturnError(errors.New("db error"))
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

			tt.setup(mock)

			storage := &PostgresStorage{
				db: sqlxDB,
			}

			err = storage.SetGauge(context.Background(), "Alloc", 5.4)

			if tt.wantErr {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
			}

			assert.NilError(t, mock.ExpectationsWereMet())
		})
	}
}
