package storage

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User таблица пользователей
type User struct {
	ID       int
	Login    string `gorm:"uniqueIndex"`
	Password string
	Cards    []DataTable
}

// DataTable таблица сохраненных данных
type DataTable struct {
	UserID   int
	ID       int
	TypeData string
	Data     string
}

func (User) TableName() string {
	return "users"
}

func (DataTable) TableName() string {
	return "data"
}

type Database struct {
	connection *gorm.DB
}

// InitDB инициализация БД
func InitDB() (*Database, error) {
	return &Database{}, nil
}

// Connect подключение к БД
func (db *Database) Connect(ctx context.Context, connStr string) (err error) {
	pdb, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return err
	}

	db.connection = pdb
	err = pdb.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	err = pdb.AutoMigrate(&DataTable{})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) Close() error {
	db1, err := db.connection.DB()
	if err != nil {
		return err
	}
	return db1.Close()
}

// Ping пинг БД
func (db *Database) Ping(ctx context.Context) error {
	db1, _ := db.connection.DB()
	if err := db1.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

// SetConnection для тестирования с помощью mock
func (db *Database) SetConnection(conn *gorm.DB) {
	db.connection = conn
}
