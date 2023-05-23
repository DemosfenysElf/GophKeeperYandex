package storage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDBPing(t *testing.T) {
	type want struct {
		codeGet  int
		response string
	}
	tests := []struct {
		name string

		errPing1 error
		errPing2 error
	}{
		{
			name: "TestDBGetPing1",
		},
		{
			name:     "TestDBGetPing2err",
			errPing1: errors.New("ping fail"),
		},

		{
			name:     "TestDBGetPing3err",
			errPing2: errors.New("ping fail"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			dialector := postgres.New(postgres.Config{
				PreferSimpleProtocol: false,
				DriverName:           "postgres",
				Conn:                 db,
			})
			mock.ExpectPing().WillReturnError(tt.errPing1)
			DB, err := gorm.Open(dialector)
			if (err != nil) && (tt.errPing1 == nil) {
				t.Fatalf("error Gorm: %s", err)
			}
			mockDB := &Database{}
			mockDB.SetConnection(DB)

			defer db.Close()
			mock.ExpectPing().WillReturnError(tt.errPing2)
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			err = mockDB.Ping(ctx)

			if (err != nil) && (tt.errPing2 == nil) {
				t.Errorf("Ошибка проверки пинга")
			}
		})
	}
}
