package menu

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"PasManagerGophKeeper/internal/router"
	"PasManagerGophKeeper/internal/service"
)

// возможные ошибки
var errFailPost = fmt.Errorf("ошибка при попытке отправки запроса на сервер")
var errDuplicateLogin = fmt.Errorf("ввёденный логин уже занят, выберите другой")
var errAllBroken = fmt.Errorf("всё поломалось, непредвиденная ошибка")
var errDataNil = fmt.Errorf("нет сохраненных данных")
var errExit = fmt.Errorf("выход")

// cheakUser меню выбора входа в клиент
func (ad *allData) cheakUser() error {
	var command int
	for {
		fmt.Println("Выберите действие.\nВведите номер: \n 1. Регистрация нового пользователя \n 2. Логин \n 3. Выход")
		fmt.Fscan(os.Stdin, &command)
		switch command {
		case 1:
			ad.registration()
			return nil
		case 2:
			ad.loginUser()
			return nil
		case 3:
			return errExit
		default:
			fmt.Println("Введено неправильное число.")
		}
		fmt.Println("Введите команду:")
	}

}

// registration регистрация нового пользователя
func (ad *allData) registration() {
	for {
		logpas := ad.testLogPass()
		err := ad.postRegistration(logpas)
		if err == nil {
			break
		}
	}
}

// loginUser вход под существующим пользователем
func (ad *allData) loginUser() {
	for {
		logpas := ad.testLogPass()
		err := ad.postLogin(logpas)
		if err == nil {
			break
		}
	}
}

// postRegistration запрос на сервер для регистрации пользователя и входа в учетную запись
func (ad *allData) postRegistration(logpas []byte) error {
	postUrl := ad.serverAddress + "/api/user/register"
	resp, err := http.Post(postUrl, "application/json", bytes.NewBuffer(logpas))
	if err != nil {
		return errFailPost
	}
	switch resp.StatusCode {
	case 200:
		// всё ок, получаем токен и идём работать дальше
		aut := resp.Header.Get(service.Authorization)
		ad.tokenJWT = aut
		return nil
	case 400:
		fmt.Println("Неверный формат запроса") //не должно сработать, причины отлавливается в testLogPass()
		return err
	case 404:
		fmt.Println("Сервер не отвечает. Попробуйте позже")
	case 409:
		fmt.Println("Логин уже занят")
		return errDuplicateLogin
	case 500:
		fmt.Println("Внутренняя ошибка сервера. Попробуйте позже")
		return errFailPost
	}
	return errAllBroken
}

// postRegistration запрос на сервер для входа в учетную запись
func (ad *allData) postLogin(logpas []byte) error {
	postUrl := ad.serverAddress + "/api/user/login"
	resp, err := http.Post(postUrl, "application/json", bytes.NewBuffer(logpas))
	if err != nil {
		return errFailPost
	}
	switch resp.StatusCode {
	case 200:
		// всё ок, получаем токен и идём работать дальше
		ad.tokenJWT = resp.Header.Get(service.Authorization)
		return nil
	case 400:
		fmt.Println("Неверный формат запроса") //не должно сработать, причины отлавливается в testLogPass()
		return err
	case 401:
		fmt.Println("Неверная пара логин/пароль") //не должно сработать, причины отлавливается в testLogPass()
		return err
	case 404:
		fmt.Println("Сервер не отвечает. Попробуйте позже")
	case 500:
		fmt.Println("Внутренняя ошибка сервера. Попробуйте позже")
		return err
	}
	return errAllBroken
}

// testLogPass меню ввода данных для входа или регистрации, с проверкой на валидность символов
func (ad *allData) testLogPass() []byte {
	var key = "u0283tyuhgjfn"

	newUser := router.User{}
	for (len(newUser.Login) == 0) && (!isTrueSym(newUser.Login)) {
		fmt.Println("Введите логин\nЛогин должен состоять из латинских букв и цифр")
		fmt.Fscan(os.Stdin, &newUser.Login)
	}

	for (!isTrueLen(newUser.Password)) && (!isTrueSym(newUser.Password)) {
		fmt.Println("Введите пароль\nПароль должен состоять из латинских букв и цифр\n и содержать 8-16 символов")
		fmt.Fscan(os.Stdin, &newUser.Password)
	}
	ad.login = newUser.Login
	ad.password = newUser.Password

	//
	hexPassword := []byte(newUser.Password)
	crypt, err := service.EnCrypt(hexPassword, key)
	if err != nil {
		return nil
	}
	newUser.Password = hex.EncodeToString(crypt)

	marshalUser, err := json.Marshal(newUser)
	if err != nil {
		return nil
	}
	return marshalUser
}

// isTrueSym проверерка на валидность символов
func isTrueSym(str string) bool {
	for _, r := range str {
		if ((r > '\u002F') && (r < '\u003A')) || ((r > '\u0040') && (r < '\u005B')) {
			return true
		}
	}
	return false
}

// isTrueLen проверерка на валидность длины пароля
func isTrueLen(str string) bool {
	if len(str) >= 8 && len(str) <= 16 {
		return true
	}
	return false
}
