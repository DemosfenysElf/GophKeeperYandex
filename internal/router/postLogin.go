package router

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo"
	"github.com/pkg/errors"

	"PasManagerGophKeeper/internal/service"
)

type registration struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *serverKeeper) postAPIUserLogin(c echo.Context) error {
	var userLog registration
	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}
	if err = json.Unmarshal(body, &userLog); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}

	tokenJWT, err := s.DB.LoginUser(c.Request().Context(), userLog.Login, userLog.Password)
	if (tokenJWT == "") && (err == nil) {
		c.Response().WriteHeader(http.StatusUnauthorized)
		return nil
	}

	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}

	c.Response().Header().Set(service.Authorization, service.Bearer+" "+tokenJWT)
	c.Response().WriteHeader(http.StatusOK)
	return nil
}

func (s *serverKeeper) postAPIUserRegister(c echo.Context) error {
	var userLog registration
	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}
	if err = json.Unmarshal(body, &userLog); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}
	var pgErr *pgconn.PgError

	tokenJWT, err := s.DB.RegisterUser(c.Request().Context(), userLog.Login, userLog.Password)
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation: // дубликат
			c.Response().WriteHeader(http.StatusConflict)
			return nil
		default:
			c.Response().WriteHeader(http.StatusInternalServerError)
			return nil
		}
	}
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}

	c.Response().Header().Set(service.Authorization, service.Bearer+" "+tokenJWT)
	c.Response().WriteHeader(http.StatusOK)
	return nil
}
