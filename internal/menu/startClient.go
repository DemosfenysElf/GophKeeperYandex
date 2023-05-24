package menu

import (
	"fmt"

	"github.com/caarlos0/env"
)

// allData конфигурация клиента
type allData struct {
	tokenJWT string
	login    string
	password string

	serverAddress  string
	LocalDownloads string `env:"USERPROFILE"` // или куда скачивать файлы
}

// initData инициализация клиента
func initData() *allData {
	return &allData{serverAddress: "http://localhost:8080"}
}

// StartClient запуск клиента
func StartClient() {
	client := initData()

	err := env.Parse(client)
	if err != nil {
		fmt.Println(err)
	}
	client.LocalDownloads += `\Downloads\`

	err = client.cheakUser()
	if (err != nil) || (client.tokenJWT == "") {
		return
	}
	client.operations()
}
