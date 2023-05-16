package menu

import (
	"fmt"
	"os"
)

func (ad *allData) Operations() {
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
		fmt.Println("Введите команду:")
	}
}

func (ad *allData) writeData() {
	var command int
	for {
		fmt.Println("Какие данные вы хотите записать?\nВведите номер: \n 1. Card \n 2. Password \n 3. Текст \n 4. Вернуться назад")
		fmt.Fscan(os.Stdin, &command)
		switch command {
		case 1:
			fmt.Println("1. Card")
			ad.writeCard()
		case 2:
			fmt.Println("2. Password")
			ad.writePassword()
		case 3:
			fmt.Println("3. Текст")
			ad.writeText()

		case 4:
			return
		}
		fmt.Println("Введите команду:")
	}
}

func (ad *allData) readData() {
	var command int
	for {
		fmt.Println("Какие данные вы хотите посмотреть?\nВведите номер: \n 1. Card \n 2. Password \n 3. Text \n 4. Вернуться назад")
		fmt.Fscan(os.Stdin, &command)
		switch command {
		case 1:
			fmt.Println("1. Card")
			err := ad.readCard()
			if err != nil {
				return
			}
		case 2:
			fmt.Println("2. Password")
		case 3:
			fmt.Println("3. Text")
		case 4:
			return
		}
		fmt.Println("Введите команду:")
	}
}
