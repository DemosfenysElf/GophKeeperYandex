package menu

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"

	"PasManagerGophKeeper/internal/router"
	"PasManagerGophKeeper/internal/service"
)

var key = "111111n"

type AllDataTest struct {
	tokenJWT string
	login    string
	password string

	serverAddress  string
	LocalDownloads string `env:"USERPROFILE"` // или куда скачивать файлы
}

func testDataCFG(data *allData, ad AllDataTest) error {
	data.login = ad.login
	data.serverAddress = ad.serverAddress
	data.LocalDownloads = ad.LocalDownloads
	data.password = ad.password

	JWT, err := service.EncodeJWT(data.login)
	if err != nil {
		return err
	}
	data.tokenJWT = service.Bearer + " " + JWT

	return nil
}

func testDataMockReturn(data interface{}, pass string) (string, error) {
	marshal, _ := json.Marshal(data)
	crypt, _ := service.EnCrypt(marshal, pass)
	bodyToString := hex.EncodeToString(crypt)
	return bodyToString, nil
}

func testLogPass(login string, pass string) ([]byte, string) {
	hexPassword := []byte(pass)
	crypt, err := service.EnCrypt(hexPassword, key)
	if err != nil {
		return nil, ""
	}
	pass = hex.EncodeToString(crypt)
	newUser := router.User{Login: login, Password: pass}
	marshalUser, err := json.Marshal(newUser)
	if err != nil {
		return nil, ""
	}

	h := md5.New()
	h.Write([]byte(newUser.Password))
	passForMock := hex.EncodeToString(h.Sum(nil))

	return marshalUser, passForMock
}
