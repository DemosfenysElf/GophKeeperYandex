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

	"PasManagerGophKeeper/internal/service"
	"PasManagerGophKeeper/internal/testsService"
)

func TestAutPost(t *testing.T) {
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
		jwt        string
		result     driver.Result
		errBD1     error
		errBD2     error
	}{
		{
			name: "TestAutPost1",
			want: want{
				codePost: 500,
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
			jwt:        "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJMb2dpbiI6ImxvZzEiLCJleHAiOjE2ODQ3NzUxMjIsImlhdCI6MTY4NDc3NTA2Mn0.2aeKc72ppFeWwta_4JBx8lH3Ex82UywV2yBv2iVBJi4",
		},
		{
			name: "TestAutPost2",
			want: want{
				codePost: 500,
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
			jwt:        "2Mn0.2aeKc72ppFeWwta_4JBx8lH3Ex82UywV2yBv2iVBJi4",
		},
		{
			name: "TestAutPost3",
			want: want{
				codePost: 401,
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
			jwt:        "",
		},
		{
			name: "TestAutPost4",
			want: want{
				codePost: 500,
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
			jwt:        "Bearer fewJ9.eyJMb2dpbiI6ImxvZzEiLCJleHAiOjE2efwwefODQ3NzUxMjIsImlhdCI6MTY4NDc3NTA2Mn0.2aeKc72ppfwex8lH3Ex82UywV2yBv2iVBJi4",
		},
		{
			name: "TestAutPost5errBD1",
			want: want{
				codePost: 500,
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
			jwt:        "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJMb2dpbiI6ImxvZzEiLCJleHAiOjE2ODQ3NzUxMjIsImlhdCI6MTY4NDc3NTA2Mn0.2aeKc72ppFeWwta_4JBx8lH3Ex82UywV2yBv2iVBJi4",
			errBD1:     errors.New("Просто ошибка"),
		},
		{
			name: "TestAutPost5errBD2",
			want: want{
				codePost: 500,
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
			jwt:        "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJMb2dpbiI6ImxvZzEiLCJleHAiOjE2ODQ3NzUxMjIsImlhdCI6MTY4NDc3NTA2Mn0.2aeKc72ppFeWwta_4JBx8lH3Ex82UywV2yBv2iVBJi4",
			errBD2:     errors.New("Просто ошибка"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := testsService.DBGormMockOnTests()
			if err != nil {
				return
			}

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
				WithArgs(tt.data.userID, tt.data.typeData, testbodyToString).WillReturnError(tt.errBD2).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			mock.ExpectCommit()
			//

			request := httptest.NewRequest(http.MethodPost, tt.postAdress, bytes.NewReader(tt.testBody))
			request.Header.Add(service.Authorization, tt.jwt)

			rout.InitRouter(e)

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
