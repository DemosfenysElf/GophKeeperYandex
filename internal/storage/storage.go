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
	//
	RegisterUser(ctx context.Context, login string, pass string) (tokenJWT string, err error)
	LoginUser(ctx context.Context, login string, pass string) (tokenJWT string, err error)
	GetUserID(ctx context.Context, login string) (UserID int, err error)
	//
	WriteCard(ctx context.Context, data string, userID int) (err error)
	ReadAllCard(ctx context.Context, userID int) (cardList []string, err error)
	WritePassword(ctx context.Context, data string, userID int) (err error)
	ReadAllPassword(ctx context.Context, userID int) (passwordList []string, err error)
	WriteText(ctx context.Context, data string, userID int) (err error)
	ReadAllText(ctx context.Context, userID int) (textList []string, err error)
	WriteBin(ctx context.Context, data string, userID int) (err error)
	ReadAllBin(ctx context.Context, userID int) (binList []string, err error)
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

// WriteCard сохранение закодированных данных банковской карты в таблицу
func (db *Database) WriteCard(ctx context.Context, data string, userID int) (err error) {
	card := Card{
		UserID: userID,
		Data:   data,
	}
	if err = db.connection.WithContext(ctx).Create(&card).Error; err != nil {
		return err
	}
	return
}

// ReadAllCard получение массива сохраненных данных банковской карты из таблицы
func (db *Database) ReadAllCard(ctx context.Context, userID int) (cardList []string, err error) {
	if err = db.connection.WithContext(ctx).Table("card").Select("data").Where("user_id = ?", userID).Scan(&cardList).Error; err != nil {
		return nil, err
	}
	return
}

// WritePassword сохранение закодированной пары логин-пароль в таблицу
func (db *Database) WritePassword(ctx context.Context, data string, userID int) (err error) {
	password := Password{
		UserID: userID,
		Data:   data,
	}
	if err = db.connection.WithContext(ctx).Create(&password).Error; err != nil {
		return err
	}
	return
}

// ReadAllPassword получение массива сохраненных пар логин-пароль из таблицы
func (db *Database) ReadAllPassword(ctx context.Context, userID int) (passwordList []string, err error) {
	if err = db.connection.WithContext(ctx).Find(&passwordList, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return
}

// WriteText сохранение закодированной текстовой заметки в таблицу
func (db *Database) WriteText(ctx context.Context, data string, userID int) (err error) {
	text := Text{
		UserID: userID,
		Data:   data,
	}
	if err = db.connection.WithContext(ctx).Create(&text).Error; err != nil {
		return err
	}
	return
}

// ReadAllText получение массива сохраненных текстовых заметок из таблицы
func (db *Database) ReadAllText(ctx context.Context, userID int) (textList []string, err error) {
	if err = db.connection.WithContext(ctx).Find(&textList, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return
}

// WriteBin сохранение закодированного файла в таблицу
func (db *Database) WriteBin(ctx context.Context, data string, userID int) (err error) {
	bin := Bin{
		UserID: userID,
		Data:   data,
	}
	if err = db.connection.WithContext(ctx).Create(&bin).Error; err != nil {
		return err
	}
	return
}

// ReadAllBin получение массива сохраненных файлов из таблицы
func (db *Database) ReadAllBin(ctx context.Context, userID int) (binList []string, err error) {
	if err = db.connection.WithContext(ctx).Find(&binList, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return
}
