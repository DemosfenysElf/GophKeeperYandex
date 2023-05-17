package menu

import (
	"fmt"
	"os"
)

func (ad *allData) operations() {
	var command int
	for {
		fmt.Println("Выберите действие.\nВведите номер: \n 1. Запись данных \n 2. Чтение данных \n 3. Выход")
		fmt.Fscan(os.Stdin, &command)
		switch command {
		case 1:
			ad.writeData()
		case 2:
			ad.readData()
		case 3:
			return
		}
		fmt.Println("Необходимо ввести команду.")
	}
}

func (ad *allData) writeData() {
	var command int
	for {
		fmt.Println("Какие данные вы хотите записать?\nВведите номер: \n " +
			"1. Данные банковской карты \n 2. Пару Логин/Пароль \n 3. Текст \n  4. Файл \n 5. Вернуться назад")
		fmt.Fscan(os.Stdin, &command)
		switch command {
		case 1:
			ad.writeCard()
		case 2:
			ad.writePassword()
		case 3:
			ad.writeText()
		case 4:
			ad.writeFile()
		case 5:
			return
		}
		fmt.Println("Необходимо ввести команду.")
	}
}

func (ad *allData) readData() {
	var command int
	for {
		fmt.Println("Какие данные вы хотите получить?\nВведите номер: \n " +
			"1. Данные банковской карты \n 2. Пару Логин/Пароль \n 3. Текст \n  4. Файл \n 5. Вернуться назад")
		fmt.Fscan(os.Stdin, &command)
		switch command {
		case 1:
			err := ad.readCard()
			if err != nil {
				return
			}
		case 2:
			err := ad.readPassword()
			if err != nil {
				return
			}
		case 3:
			err := ad.readText()
			if err != nil {
				return
			}
		case 4:
			err := ad.readFile()
			if err != nil {
				return
			}
		case 5:
			return
		}
		fmt.Println("Необходимо ввести команду.")
	}
}
