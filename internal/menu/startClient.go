package menu

type allData struct {
	tokenJWT string
	login    string
	password string

	serverAddress string
}

func initData() *allData {
	return &allData{serverAddress: "http://localhost:8080"}
}

func StartClient() {
	client := initData()
	err := client.cheakUser()
	if (err != nil) || (client.tokenJWT == "") {
		return
	}
	client.operations()
}
