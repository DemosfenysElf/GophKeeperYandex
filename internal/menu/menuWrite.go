package menu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"PasManagerGophKeeper/internal/service"
)

type bankCard struct {
	CardName string
	Number   int
	Name     string
	Date     string
	Csv      int
}

type savePassword struct {
	ServiceName string
	Login       string
	Password    string
}

type saveText struct {
	TextName string
	Text     string
}

type saveFile struct {
	FileName string
	FileData []byte
}

func (ad *allData) writeCard() {
	var bC bankCard
	var command int
	for {
		fmt.Println("Введите данные карты.")
		fmt.Print("Введите название карты: ")
		fmt.Fscan(os.Stdin, &bC.CardName)
		fmt.Print("Введите номер карты: ")
		fmt.Fscan(os.Stdin, &bC.Number)
		fmt.Print("Введите держателя карты: ")
		fmt.Fscan(os.Stdin, &bC.Name)
		fmt.Print("Введите срок действия карты: ")
		fmt.Fscan(os.Stdin, &bC.Date)
		fmt.Print("Введите код CSV карты: ")
		fmt.Fscan(os.Stdin, &bC.Csv)

		fmt.Println("\nСохранить карту:", bC.CardName, "?")
		fmt.Println("Введите номер: \n 1. Сохранить \n 2. Ввести данные заново \n 3. Вернуться назад")
		fmt.Fscan(os.Stdin, &command)
		switch command {
		case 3:
			return
		case 1:
			marshal, err := json.Marshal(bC)
			if err != nil {
				fmt.Println("Ошибка при формировании данных для отправки на сервер")
				return
			}
			err = ad.postWrite(marshal, "/write/card")
			if err != nil {
				fmt.Println("Ошибка отправки данных на сервер")
				return
			}
			return

		default:
			fmt.Println("2. Цикл заново")
		}
	}
}

func (ad *allData) writePassword() {
	var sP savePassword
	var command int
	for {
		fmt.Println("Введите пару логин/пароль которые хотите сохранить.")
		fmt.Print("Введите название сервиса: ")
		fmt.Fscan(os.Stdin, &sP.ServiceName)
		fmt.Print("Введите логин: ")
		fmt.Fscan(os.Stdin, &sP.Login)
		fmt.Print("Введите пароль: ")
		fmt.Fscan(os.Stdin, &sP.Password)
		fmt.Print("Введите срок действия карты: ")

		fmt.Println("\nСохранить данные для ", sP.ServiceName, "?")
		fmt.Println("Введите номер: \n 1. Сохранить \n 2. Ввести данные заново \n 3. Вернуться назад")
		fmt.Fscan(os.Stdin, &command)
		switch command {
		case 3:
			return
		case 1:
			marshal, err := json.Marshal(sP)
			if err != nil {
				fmt.Println("Ошибка при формировании данных для отправки на сервер")
				return
			}
			err = ad.postWrite(marshal, "/write/password")
			if err != nil {
				fmt.Println("Ошибка отправки данных на сервер")
				return
			}
			return

		default:
			fmt.Println("2. Цикл заново")
		}
	}
}

func (ad *allData) writeText() {
	var sT saveText
	var command int
	for {
		fmt.Println("Введите пару логин/пароль которые хотите сохранить.")
		fmt.Print("Введите название тестовой заметки: ")
		fmt.Fscan(os.Stdin, &sT.TextName)
		fmt.Print("Введите текст заметки: ")
		fmt.Fscan(os.Stdin, &sT.Text)

		fmt.Println("\nСохранить данные для ", sT.TextName, "?")
		fmt.Println("Введите номер: \n 1. Сохранить \n 2. Ввести данные заново \n 3. Вернуться назад")
		fmt.Fscan(os.Stdin, &command)
		switch command {
		case 3:
			return
		case 1:
			marshal, err := json.Marshal(sT)
			if err != nil {
				fmt.Println("Ошибка при формировании данных для отправки на сервер")
				return
			}
			err = ad.postWrite(marshal, "/write/text")
			if err != nil {
				fmt.Println("Ошибка отправки данных на сервер")
				return
			}
			return

		default:
			fmt.Println("2. Цикл заново")
		}
	}
}

func (ad *allData) writeFile() {
	var sF saveFile
	var command int
	for {
		fmt.Println("Введите пару логин/пароль которые хотите сохранить.")
		fmt.Print("Введите название путь до файла: ")
		fmt.Fscan(os.Stdin, &sF.FileName)
		//fmt.Print("Введите текст заметки: ")
		//fmt.Fscan(os.Stdin, &sF.FileData)

		//

		file, err := os.Open(sF.FileName)
		if err != nil {
			log.Fatalln(err)
		}

		_, err = file.Read(sF.FileData)
		if err != nil {
			return
		}
		sF.FileName = file.Name()

		//

		fmt.Println("\nСохранить данные для ", sF.FileName, "?")
		fmt.Println("Введите номер: \n 1. Сохранить \n 2. Ввести данные заново \n 3. Вернуться назад")
		fmt.Fscan(os.Stdin, &command)
		switch command {
		case 3:
			return
		case 1:
			marshal, err := json.Marshal(sF)
			if err != nil {
				fmt.Println("Ошибка при формировании данных для отправки на сервер")
				return
			}
			err = ad.postWrite(marshal, "/write/bin")
			if err != nil {
				fmt.Println("Ошибка отправки данных на сервер")
				return
			}
			return

		default:
			fmt.Println("2. Цикл заново")
		}
	}
}

func (ad *allData) postWrite(data []byte, operation string) error {
	cryptData, err := service.EnCrypt(data, ad.password)
	if err != nil {
		return err
	}

	post := ad.serverAddress + operation
	req, err := http.NewRequest("POST", post, bytes.NewBuffer(cryptData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(service.Authorization, ad.tokenJWT)

	cli := http.Client{}

	do, err := cli.Do(req)
	if err != nil {
		return err
	}

	if err != nil {
		return errAllBroken
	}

	switch do.StatusCode {
	case 202:
		fmt.Println("Данные сохранены")
		return nil
	case 404:
		fmt.Println("Сервер не отвечает. Попробуйте позже")
	case 500:
		fmt.Println("Внутренняя ошибка сервера. Попробуйте позже")
		return errFailPost
	}

	return errAllBroken
}
