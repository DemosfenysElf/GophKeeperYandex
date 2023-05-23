package router

import (
	"bytes"
	"crypto/md5"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"PasManagerGophKeeper/internal/service"
	"PasManagerGophKeeper/internal/storage"
)

func TestPostRegistr(t *testing.T) {
	type want struct {
		codePost int
	}
	type user struct {
		Login    string
		Password string
	}
	tests := []struct {
		name     string
		want     want
		testBody []byte
		user     user
		result   driver.Result
		errBD    error
	}{
		{
			name: "TestPostRegistr1",
			want: want{
				codePost: 200,
			},
			user: user{
				Login:    "log1",
				Password: "pass1",
			},
			errBD: nil,
		},
		{
			name: "TestPostRegistr2",
			want: want{
				codePost: 409,
			},
			user: user{
				Login:    "log1",
				Password: "pass1",
			},
			errBD: gorm.ErrDuplicatedKey,
		},
		{
			name: "TestPostRegistr3",
			want: want{
				codePost: 500,
			},
			user: user{
				Login:    "log1",
				Password: "pass1",
			},
			errBD: errors.New("Просто ошибка"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var key = "u0283tyuhgjfn"
			newUser := tt.user
			hexPassword := []byte(newUser.Password)
			crypt, err := service.EnCrypt(hexPassword, key)
			if err != nil {
				t.Fatalf("error %s", err)
			}
			newUser.Password = hex.EncodeToString(crypt)
			marshalUser, err := json.Marshal(newUser)
			if err != nil {
				t.Fatalf("error %s", err)
			}
			h := md5.New()
			h.Write([]byte(newUser.Password))
			newUser.Password = hex.EncodeToString(h.Sum(nil))

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
			rout := InitServer()
			e := echo.New()
			rout.DB = mockDB

			mock.ExpectBegin()
			mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" `)).
				WithArgs(newUser.Login, newUser.Password).
				WillReturnError(tt.errBD).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			if tt.errBD != nil {
				mock.ExpectRollback()
			}
			mock.ExpectCommit()

			request := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(marshalUser))
			rout.initRouter(e)
			responseRecorder := httptest.NewRecorder()
			e.ServeHTTP(responseRecorder, request)
			response := responseRecorder.Result()
			defer response.Body.Close()
			if response.StatusCode != tt.want.codePost {
				t.Errorf("Expected status code %d, got %d", tt.want.codePost, responseRecorder.Code)
			}
			if response.StatusCode == 200 {
				headerAuth := response.Header.Get(service.Authorization)
				if headerAuth == "" {
					t.Errorf("JWT отсутствует")
				}
				headerParts := strings.Split(headerAuth, " ")
				if len(headerParts) != 2 || headerParts[0] != service.Bearer {
					t.Errorf("JWT отсутствует")
				}
			}
		})
	}
}

func TestPostLogin(t *testing.T) {
	type want struct {
		codePost int
	}
	type user struct {
		Login    string
		Password string
	}
	tests := []struct {
		name      string
		want      want
		testBody  []byte
		user      user
		result    driver.Result
		userFalse string
		errBD     error
	}{
		{
			name: "TestPostLogin1",
			want: want{
				codePost: 200,
			},
			user: user{
				Login:    "log1",
				Password: "pass1",
			},
			userFalse: "log1",
		},
		{
			name: "TestPostLogin2",
			want: want{
				codePost: 500,
			},
			user: user{
				Login:    "log1",
				Password: "pass1",
			},
			userFalse: "llg2",
		},
		{
			name: "TestPostLogin2",
			want: want{
				codePost: 500,
			},
			user: user{
				Login:    "log1",
				Password: "pass1",
			},
			userFalse: "llg2",
			errBD:     errors.New("Просто ошибка"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var key = "u0283tyuhgjfn"
			newUser := tt.user
			hexPassword := []byte(newUser.Password)
			crypt, err := service.EnCrypt(hexPassword, key)
			if err != nil {
				t.Fatalf("error %s", err)
			}
			newUser.Password = hex.EncodeToString(crypt)
			marshalUser, err := json.Marshal(newUser)
			if err != nil {
				t.Fatalf("error %s", err)
			}
			h := md5.New()
			h.Write([]byte(newUser.Password))
			newUser.Password = hex.EncodeToString(h.Sum(nil))

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
			rout := InitServer()
			e := echo.New()
			rout.DB = mockDB

			row := sqlmock.NewRows([]string{"ID", "Login", "Password"}).
				AddRow("1", newUser.Login, newUser.Password)
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE `)).
				WithArgs(tt.userFalse).WillReturnRows(row).WillReturnError(tt.errBD)

			request := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(marshalUser))
			rout.initRouter(e)
			responseRecorder := httptest.NewRecorder()
			e.ServeHTTP(responseRecorder, request)
			response := responseRecorder.Result()
			defer response.Body.Close()
			if response.StatusCode != tt.want.codePost {
				t.Errorf("Expected status code %d, got %d", tt.want.codePost, responseRecorder.Code)
			}
			if response.StatusCode == 200 {
				headerAuth := response.Header.Get(service.Authorization)
				if headerAuth == "" {
					t.Errorf("JWT отсутствует")
				}
				headerParts := strings.Split(headerAuth, " ")
				if len(headerParts) != 2 || headerParts[0] != service.Bearer {
					t.Errorf("JWT отсутствует")
				}
			}
		})
	}
}
