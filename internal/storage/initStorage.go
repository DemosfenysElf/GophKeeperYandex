package storage

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        int
	Login     string `gorm:"uniqueIndex"`
	Password  string
	Cards     []Card
	Passwords []Password
	Texts     []Text
	Bins      []Bin
}

type Card struct {
	UserID int
	ID     int
	Data   string
}

type Password struct {
	UserID int
	ID     int
	Data   string
}

type Text struct {
	UserID int
	ID     int
	Data   string
}

type Bin struct {
	UserID int
	ID     int
	Data   string
}

func (User) TableName() string {
	return "users"
}

func (Password) TableName() string {
	return "password"
}

func (Text) TableName() string {
	return "text"
}

func (Bin) TableName() string {
	return "binare"
}

func (Card) TableName() string {
	return "card"
}

type Database struct {
	connection *gorm.DB
}

func InitDB() (*Database, error) {
	return &Database{}, nil
}

func (db *Database) Connect(ctx context.Context, connStr string) (err error) {
	pdb, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return err
	}

	db.connection = pdb
	pdb.AutoMigrate(&User{})
	pdb.AutoMigrate(&Password{})
	pdb.AutoMigrate(&Text{})
	pdb.AutoMigrate(&Bin{})
	pdb.AutoMigrate(&Card{})
	return nil
}

func (db *Database) Close() error {
	db1, err := db.connection.DB()
	if err != nil {
		return err
	}
	return db1.Close()
}

func (db *Database) Ping(ctx context.Context) error {
	db1, _ := db.connection.DB()
	if err := db1.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
