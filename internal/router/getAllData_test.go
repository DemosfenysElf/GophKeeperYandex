package router

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"io"
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

func TestAutGet(t *testing.T) {
	type want struct {
		codeGet int
	}
	type user struct {
		id       int
		login    string
		password string
	}
	tests := []struct {
		name      string
		want      want
		user      user
		data      [][]byte
		getAdress string
		getType   string
		result    driver.Result
		errBD1    error
		errBD2    error
	}{
		{
			name: "TestGet1",
			want: want{
				codeGet: 202,
			},
			user: user{
				id:       5,
				login:    "log1",
				password: "pass1",
			},
			data:      [][]byte{[]byte("d11"), []byte("d1521"), []byte("gsdg231"), []byte("11r111")},
			getAdress: "/read/card",
			getType:   "card",
		},
		{
			name: "TestGet2",
			want: want{
				codeGet: 202,
			},
			user: user{
				id:       5,
				login:    "log1",
				password: "pass1",
			},
			data:      [][]byte{[]byte("d11"), []byte("d1521"), []byte("gsdg231"), []byte("gsdg2231"), []byte("gsdg2131"), []byte("gsd1g231"), []byte("gsd4124g231"), []byte("11r111")},
			getAdress: "/read/text",
			getType:   "text",
		},
		{
			name: "TestGet3",
			want: want{
				codeGet: 202,
			},
			user: user{
				id:       5,
				login:    "log1",
				password: "pass1",
			},
			data:      [][]byte{[]byte("d11"), []byte("d1521"), []byte("gsdg231"), []byte("11r111")},
			getAdress: "/read/text",
			getType:   "text",
		},
		{
			name: "TestGet4errBD1",
			want: want{
				codeGet: 500,
			},
			user: user{
				id:       5,
				login:    "log1",
				password: "pass1",
			},
			data:      [][]byte{[]byte("d11"), []byte("d1521"), []byte("gsdg231"), []byte("11r111")},
			getAdress: "/read/text",
			getType:   "text",
			errBD1:    errors.New("Просто ошибка"),
		},
		{
			name: "TestGet4errBD2",
			want: want{
				codeGet: 500,
			},
			user: user{
				id:       5,
				login:    "log1",
				password: "pass1",
			},
			data:      [][]byte{[]byte("d11"), []byte("d1521"), []byte("gsdg231"), []byte("11r111")},
			getAdress: "/read/text",
			getType:   "text",
			errBD2:    errors.New("Просто ошибка"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenJWT, err := service.EncodeJWT(tt.user.login)
			if err != nil {
				t.Errorf("Ошибка формирование JWT")
			}

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

			row1 := sqlmock.NewRows([]string{"data"})
			for i := range tt.data {
				testbodyToString := hex.EncodeToString(tt.data[i])
				row1.AddRow(testbodyToString)

			}

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT data FROM "data" WHERE`)).
				WithArgs(tt.user.id, tt.getType).WillReturnRows(row1).WillReturnError(tt.errBD2)
			//

			request := httptest.NewRequest(http.MethodGet, tt.getAdress, nil)
			request.Header.Add(service.Authorization, service.Bearer+" "+tokenJWT)

			rout.InitRouter(e)

			responseRecorder := httptest.NewRecorder()
			e.ServeHTTP(responseRecorder, request)
			//

			response := responseRecorder.Result()
			defer response.Body.Close()
			if response.StatusCode != tt.want.codeGet {
				t.Errorf("Expected status code %d, got %d", tt.want.codeGet, responseRecorder.Code)
			}

			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Errorf("Неполучилось проверить ответ")
			}
			defer response.Body.Close()

			var mstr [][]byte
			json.Unmarshal(body, &mstr)
			for i := range mstr {
				if string(mstr[i]) != string(tt.data[i]) {
					t.Errorf("Полученные данные не совпадают с исходными")
				}
			}
		})
	}
}
