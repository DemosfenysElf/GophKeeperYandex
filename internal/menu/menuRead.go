package menu

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"PasManagerGophKeeper/internal/service"
)

func (ad allData) readCard() error {
	var bC bankCard
	var bCs []bankCard
	var number int

	cards, err := ad.getRead("/read/card")
	if err != nil {
		return err
	}

	for _, bytes := range cards {
		err = json.Unmarshal(bytes, &bC)
		if err != nil {
			fmt.Println(err)
			return err
		}
		bCs = append(bCs, bC)
	}

	if len(bCs) == 0 {
		fmt.Println("Сохраненные карты отсутствуют")

		return errDataNil
	}

	fmt.Println("Список всех сохраненных карт:")
	for i, card := range bCs {
		fmt.Println(i+1, " ", card.CardName)
	}
	for {
		fmt.Println("Введите номер карты для отображения:")
		fmt.Fscan(os.Stdin, &number)
		if !(number <= 0 || (number) > len(bCs)) {
			break
		}
		fmt.Println("Несоответствующий номер")

	}
	fmt.Println("Ваша карта              ", bCs[number-1].CardName)
	fmt.Println("Номер карты             ", bCs[number-1].Number)
	fmt.Println("Дата окончания действия ", bCs[number-1].Date)
	fmt.Println("Зарегистрирована на     ", bCs[number-1].Name)
	fmt.Println("Код Csv                 ", bCs[number-1].Csv)

	return nil
}

func (ad allData) readPassword() error {
	var sP savePassword
	var sPs []savePassword
	var number int

	passwords, err := ad.getRead("/read/password")
	if err != nil {
		return err
	}

	for _, bytes := range passwords {
		err = json.Unmarshal(bytes, &sP)
		if err != nil {
			fmt.Println(err)
			return err
		}
		sPs = append(sPs, sP)
	}

	if len(sPs) == 0 {
		fmt.Println("Сохраненные пароли отсутствуют")
		return errDataNil
	}

	fmt.Println("Список всех сохраненных паролей:")
	for i, pass := range sPs {
		fmt.Println(i+1, " ", pass.ServiceName)
	}
	for {
		fmt.Println("Введите номер сохраненного пароля для отображения:")
		fmt.Fscan(os.Stdin, &number)
		if !(number <= 0 || (number) > len(sPs)) {
			break
		}
		fmt.Println("Несоответствующий номер")

	}
	fmt.Println("Сервис, где используется   ", sPs[number-1].ServiceName)
	fmt.Println("Логин                      ", sPs[number-1].Login)
	fmt.Println("Пароль                     ", sPs[number-1].Password)
	return nil
}

func (ad allData) readText() error {
	var sT saveText
	var sTs []saveText
	var number int

	passwords, err := ad.getRead("/read/text")
	if err != nil {
		return err
	}

	for _, bytes := range passwords {
		err = json.Unmarshal(bytes, &sT)
		if err != nil {
			fmt.Println(err)
			return err
		}
		sTs = append(sTs, sT)
	}

	if len(sTs) == 0 {
		fmt.Println("Сохраненные заметки отсутствуют")

		return errDataNil
	}

	fmt.Println("Список всех сохраненных заметок:")
	for i, text := range sTs {
		fmt.Println(i+1, " ", text.TextName)
	}
	for {
		fmt.Println("Введите номер заметки для отображения:")
		fmt.Fscan(os.Stdin, &number)
		if !(number <= 0 || (number) > len(sTs)) {
			break
		}
		fmt.Println("Несоответствующий номер")

	}
	fmt.Println("Название заметки   ", sTs[number-1].TextName)
	fmt.Println("Текст заметки      ", sTs[number-1].Text)

	return nil
}

func (ad allData) getRead(operation string) ([][]byte, error) {
	get := ad.serverAddress + operation
	req, err := http.NewRequest("GET", get, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(service.Authorization, ad.tokenJWT)
	cli := http.Client{}
	do, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	switch do.StatusCode {
	case 202:
		body, err := io.ReadAll(do.Body)
		defer do.Body.Close()
		if err != nil {
			return nil, err
		}
		var mByte [][]byte
		err = json.Unmarshal(body, &mByte)
		if err != nil {
			return nil, err
		}
		for i := range mByte {
			deCryptData, err := service.DeCrypt([]byte(mByte[i]), ad.password)
			if err != nil {
				return nil, err
			}
			mByte[i] = deCryptData
		}
		return mByte, nil
	case 404:
		fmt.Println("Сервер не отвечает. Попробуйте позже")
	case 500:
		fmt.Println("Внутренняя ошибка сервера. Попробуйте позже")
		return nil, errFailPost
	}

	return nil, errAllBroken
}
