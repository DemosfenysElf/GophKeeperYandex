package testsService

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"PasManagerGophKeeper/internal/storage"
)

func DBGormMockOnTests() (*storage.Database, sqlmock.Sqlmock, error) {
	db, mock, _ := sqlmock.New()

	dialector := postgres.New(postgres.Config{
		PreferSimpleProtocol: false,
		DriverName:           "postgres",
		Conn:                 db,
	})
	DB, _ := gorm.Open(dialector)

	mockDB := &storage.Database{}
	mockDB.SetConnection(DB)
	return mockDB, mock, nil
}
