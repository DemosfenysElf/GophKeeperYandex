package menu

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"

	"PasManagerGophKeeper/internal/router"
	"PasManagerGophKeeper/internal/service"
	"PasManagerGophKeeper/internal/testsService"
)

func TestClientPostLogin(t *testing.T) {
	tests := []struct {
		name  string
		user  router.User
		ad    AllDataTest
		port  string
		errBD error
	}{
		{
			name: "TestLogin1",
			ad: AllDataTest{
				login:    "login1",
				password: "Password",

				serverAddress: "http://localhost:8081",
			},
			port:  ":8081",
			errBD: nil,
		},
		{
			name: "TestRegistr1",
			ad: AllDataTest{
				login:    "login1",
				password: "Password",

				serverAddress: "http://localhost:8081",
			},
			port:  ":8081",
			errBD: errors.New("Просто ошибка"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataC := initData()
			testDataCFG(dataC, tt.ad)
			logpass, passMock := testLogPass(dataC.login, dataC.password)

			mockDB, mock, err := testsService.DBGormMockOnTests()
			if err != nil {
				return
			}

			mock.ExpectBegin()
			mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" `)).
				WithArgs(dataC.login, passMock).
				WillReturnError(tt.errBD).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			if tt.errBD != nil {
				mock.ExpectRollback()
			}
			mock.ExpectCommit()

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

				// само тестирование
				err = dataC.postLogin(logpass)
				if (err != nil) && (err != errDuplicateLogin) {
					t.Errorf("Ошибка postRegistration  %s", err)
				}

				if dataC.tokenJWT == "" {
					t.Errorf("JWT отсутствует")
				}
				headerParts := strings.Split(dataC.tokenJWT, " ")
				if len(headerParts) != 2 || headerParts[0] != service.Bearer {
					t.Errorf("JWT отсутствует")
				}
			}()
			e.Start(tt.port)
		})
	}
}
