package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestPostgresStorage_GetCounter(t *testing.T) {
	tests := []struct {
		name    string
		argName string
		setup   func(sqlmock.Sqlmock)
		want    int64
		wantErr error
	}{
		{
			name:    "success",
			argName: "Count",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(queryGetCounter)).
					WithArgs("Count").
					WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(54))
			},
			want:    54,
			wantErr: nil,
		},
		{
			name:    "counter not found",
			argName: "Unknown",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(queryGetCounter)).
					WithArgs("Unknown").
					WillReturnError(sql.ErrNoRows)
			},
			want:    0,
			wantErr: ErrCounterNotFound,
		},
		{
			name:    "db error",
			argName: "Count",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(queryGetCounter)).
					WithArgs("Count").
					WillReturnError(errors.New("db error"))
			},
			want:    0,
			wantErr: errors.New("db error"),
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

			got, err := storage.GetCounter(context.Background(), tt.argName)
			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NilError(t, err)
			}

			assert.Equal(t, tt.want, got)
			assert.NilError(t, mock.ExpectationsWereMet())
		})
	}
}
