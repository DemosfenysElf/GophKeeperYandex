package menu

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"

	"PasManagerGophKeeper/internal/router"
	"PasManagerGophKeeper/internal/testsService"
)

func TestRead(t *testing.T) {
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
			text: saveText{
				TextName: "text1",
				Text:     "1222222fewfgwwew222222",
			},
			typeOp:   "/read/text",
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
			mockReturn, err := testDataMockReturn(tt.text, dataC.password)
			if err != nil {
				return
			}

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
				WithArgs(dataC.login).WillReturnRows(row).WillReturnError(nil)

			row1 := sqlmock.NewRows([]string{"data"}).AddRow(mockReturn).AddRow(mockReturn)

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT data FROM "data" WHERE`)).
				WithArgs(tt.idUser, tt.typeData).WillReturnRows(row1).WillReturnError(nil)

			// server
			rout := router.InitServer()
			e := echo.New()
			rout.DB = mockDB
			rout.InitRouter(e)
			rout.Cfg.ServerAddress = tt.ad.serverAddress
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
				defer cancel()
				defer func() {
					time.Sleep(time.Millisecond * 500)
					e.Shutdown(ctx)
				}()

				read, err := dataC.getRead(tt.typeOp)
				if err != nil {
					t.Errorf("Я уже не понимаю что происходит 2 %s", err)
				}

				var sT saveText
				var sTs []saveText
				for _, bytes := range read {
					err = json.Unmarshal(bytes, &sT)
					if err != nil {
						t.Errorf("Я уже не понимаю что происходит 3 %s", err)
					}
					sTs = append(sTs, sT)
				}
				for i := range sTs {
					if (sTs[i].Text != tt.text.Text) || (sTs[i].TextName != tt.text.TextName) {
						t.Errorf("Данные не соответствуют ожидаемым")
					}
				}
				fmt.Println("всё ок")
			}()
			e.Start(tt.port)
		})
	}
}
