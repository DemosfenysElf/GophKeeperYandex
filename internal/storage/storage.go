package storage

import (
	"context"
	"crypto/md5"
	"encoding/hex"

	_ "github.com/jackc/pgx/v5/stdlib"

	"PasManagerGophKeeper/internal/service"
)

type DBI interface {
	Connect(ctx context.Context, connStr string) (err error)
	Ping(ctx context.Context) error
	Close() error

	RegisterUser(ctx context.Context, login string, pass string) (tokenJWT string, err error)
	LoginUser(ctx context.Context, login string, pass string) (tokenJWT string, err error)
	GetUserID(ctx context.Context, login string) (UserID int, err error)

	WriteData(ctx context.Context, data string, userID int, typeData string) (err error)
	ReadAllDataType(ctx context.Context, userID int, typeData string) (dataList []string, err error)
}

// RegisterUser регистрация пользователя
func (db *Database) RegisterUser(ctx context.Context, login string, pass string) (tokenJWT string, err error) {
	h := md5.New()
	h.Write([]byte(pass))
	passHex := hex.EncodeToString(h.Sum(nil))
	user := User{
		Login:    login,
		Password: passHex,
	}

	if err = db.connection.WithContext(ctx).Create(&user).Error; err != nil {
		return "", err
	}

	tokenJWT, err = service.EncodeJWT(login)
	if err != nil {
		return "", err
	}
	return tokenJWT, nil
}

// LoginUser авторизация пользователя
func (db *Database) LoginUser(ctx context.Context, login string, pass string) (tokenJWT string, err error) {
	h := md5.New()
	h.Write([]byte(pass))
	pass = hex.EncodeToString(h.Sum(nil))
	user := User{}

	if err = db.connection.WithContext(ctx).Find(&user, "login = ?", login).Error; err != nil {
		return "", err
	}

	if user.Password != pass {
		return "", nil
	}
	tokenJWT, err = service.EncodeJWT(login)
	if err != nil {
		return "", err
	}
	return tokenJWT, nil
}

// GetUserID получение userID, для дальнейшего сохранения данных пользователя в таблицы
func (db *Database) GetUserID(ctx context.Context, login string) (UserID int, err error) {
	user := User{}
	if err = db.connection.WithContext(ctx).Find(&user, "login = ?", login).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

// WriteData сохранение данных в таблицу
func (db *Database) WriteData(ctx context.Context, data string, userID int, typeData string) (err error) {
	oneData := DataTable{
		UserID:   userID,
		typeData: typeData,
		Data:     data,
	}
	if err = db.connection.WithContext(ctx).Create(&oneData).Error; err != nil {
		return err
	}
	return
}

// ReadAllDataType получение массива сохраненных данных из таблицы
func (db *Database) ReadAllDataType(ctx context.Context, userID int, typeData string) (dataList []string, err error) {
	if err = db.connection.WithContext(ctx).Table("card").Select("data").
		Where("user_id = ?", userID).Where("type_data = ?", typeData).Scan(&dataList).Error; err != nil {
		return nil, err
	}
	return
}
