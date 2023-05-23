package router

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"

	"github.com/labstack/echo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"PasManagerGophKeeper/internal/service"
	"PasManagerGophKeeper/internal/storage"
)

func TestPost(t *testing.T) {
	type want struct {
		codePost int
	}
	type user struct {
		id       int
		login    string
		password string
	}
	type data struct {
		id       int
		userID   int
		typeData string
		data     string
	}
	tests := []struct {
		name       string
		want       want
		wantErr    bool
		testBody   []byte
		user       user
		data       data
		postAdress string
		result     driver.Result
		errBD1     error
		errBD2     error
	}{
		{
			name: "TestPost1",
			want: want{
				codePost: 202,
			},
			user: user{
				id:       1,
				login:    "log1",
				password: "pass1",
			},
			data: data{
				id:       4,
				userID:   1,
				typeData: "card",
				data:     "d11",
			},
			postAdress: "/write/card",
			testBody:   []byte("234tpnDA"),
		},
		{
			name: "TestPost2",
			want: want{
				codePost: 202,
			},
			user: user{
				id:       2,
				login:    "log2",
				password: "pas342s1",
			},
			data: data{
				id:       5,
				userID:   2,
				typeData: "text",
				data:     "d11",
			},
			postAdress: "/write/text",
			testBody:   []byte("234tpnDA"),
		},
		{
			name: "TestPost3errBD1",
			want: want{
				codePost: 500,
			},
			user: user{
				id:       2,
				login:    "log2",
				password: "pas342s1",
			},
			data: data{
				id:       5,
				userID:   2,
				typeData: "text",
				data:     "d11",
			},
			postAdress: "/write/text",
			testBody:   []byte("234tpnDA"),
			errBD1:     errors.New("Просто ошибка"),
		},
		{
			name: "TestPost4errBD2",
			want: want{
				codePost: 500,
			},
			user: user{
				id:       2,
				login:    "log2",
				password: "pas342s1",
			},
			data: data{
				id:       5,
				userID:   2,
				typeData: "text",
				data:     "d11",
			},
			postAdress: "/write/text",
			testBody:   []byte("234tpnDA"),
			errBD2:     errors.New("Просто ошибка"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenJWT, err := service.EncodeJWT(tt.user.login)
			if err != nil {
				t.Errorf("Ошибка формирование JWT")
			}

			//
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			dialector := postgres.New(postgres.Config{
				PreferSimpleProtocol: false,
				DriverName:           "postgres",
				Conn:                 db,
			})
			DB, err := gorm.Open(dialector)
			if err != nil {
				t.Fatalf("error Gorm: %s", err)
			}
			mockDB := &storage.Database{}
			mockDB.SetConnection(DB)
			//

			rout := InitServer()
			e := echo.New()
			rout.DB = mockDB

			row := sqlmock.NewRows([]string{"ID", "Login", "Password"}).
				AddRow(tt.user.id, tt.user.login, tt.user.password)

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE `)).
				WithArgs(tt.user.login).WillReturnRows(row).WillReturnError(tt.errBD1)

			testbodyToString := hex.EncodeToString(tt.testBody)
			mock.ExpectBegin()
			mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "data" `)).
				WithArgs(tt.data.userID, tt.data.typeData, testbodyToString).
				WillReturnError(tt.errBD2).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			mock.ExpectCommit()
			//

			request := httptest.NewRequest(http.MethodPost, tt.postAdress, bytes.NewReader(tt.testBody))
			request.Header.Add(service.Authorization, service.Bearer+" "+tokenJWT)

			rout.initRouter(e)

			responseRecorder := httptest.NewRecorder()
			e.ServeHTTP(responseRecorder, request)
			//

			response := responseRecorder.Result()
			defer response.Body.Close()
			if response.StatusCode != tt.want.codePost {
				t.Errorf("Expected status code %d, got %d", tt.want.codePost, responseRecorder.Code)
			}

		})
	}
}
