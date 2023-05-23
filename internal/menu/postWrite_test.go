package menu

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"

	"PasManagerGophKeeper/internal/router"
	"PasManagerGophKeeper/internal/service"
	"PasManagerGophKeeper/internal/testsService"
)

func TestWrite(t *testing.T) {
	tests := []struct {
		name     string
		bCard    bankCard
		text     saveText
		typeOp   string
		typeData string
		ad       AllDataTest
		idUser   int
		port     string
	}{
		{
			name: "Test1",
			bCard: bankCard{
				CardName: "card1",
				Number:   1234567890,
				Name:     "Oleg",
				Date:     "02/25",
				Csv:      132,
			},
			typeOp:   "/write/card",
			typeData: "card",
			ad: AllDataTest{
				login:    "login1",
				password: "Password",

				serverAddress: "http://localhost:8081",
			},
			idUser: 1,
			port:   ":8081",
		},
		{
			name: "Test2",
			text: saveText{
				TextName: "card2",
				Text:     "1222222222222",
			},
			typeOp:   "/write/text",
			typeData: "text",
			ad: AllDataTest{
				login:    "login1",
				password: "Password",

				serverAddress: "http://localhost:8081",
			},
			idUser: 1,
			port:   ":8081",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataC := initData()
			testDataCFG(dataC, tt.ad)

			// Данные
			marshalData, err := json.Marshal(tt.bCard)
			if err != nil {
				t.Errorf("Marshal err %s", err)
			}
			cryptData, err := service.EnCrypt(marshalData, dataC.password)
			if err != nil {
				t.Errorf("EnCrypt data err %s", err)
			}
			dataDBtoString := hex.EncodeToString(cryptData)

			// server
			mockDB, mock, err := testsService.DBGormMockOnTests()
			if err != nil {
				return
			}
			// mockDB
			// getUserID
			row := sqlmock.NewRows([]string{"ID", "Login", "Password"}).
				AddRow(tt.idUser, tt.ad.login, "tt.user.password")
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE `)).
				WithArgs(tt.ad.login).WillReturnRows(row).WillReturnError(nil)
			// writeData
			mock.ExpectBegin()
			mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "data" `)).
				WithArgs(tt.idUser, tt.typeData, dataDBtoString).
				WillReturnError(nil).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			mock.ExpectCommit()

			// server
			rout := router.InitServer()
			e := echo.New()
			rout.DB = mockDB
			rout.InitRouter(e)
			rout.Cfg.ServerAddress = tt.ad.serverAddress
			go func() {
				fmt.Println("aaaa")
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
				defer cancel()
				defer func() {
					time.Sleep(time.Millisecond * 500)
					e.Shutdown(ctx)
				}()

				// само тестирование
				err = dataC.postWrite(marshalData, tt.typeOp)
				if err != nil {
					t.Errorf("Я уже не понимаю что происходит 2 %s", err)
				}

			}()
			e.Start(tt.port)
		})
	}
}
