package router

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"gorm.io/gorm"

	"PasManagerGophKeeper/internal/service"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// postAPIUserRegister регистрация нового пользователя
func (s *serverKeeper) postAPIUserRegister(c echo.Context) error {
	var userLog User
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
	tokenJWT, err := s.DB.RegisterUser(c.Request().Context(), userLog.Login, userLog.Password)
	if err == gorm.ErrDuplicatedKey {
		c.Response().WriteHeader(http.StatusConflict)
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

// postAPIUserLogin вход в существующую учетную запись
func (s *serverKeeper) postAPIUserLogin(c echo.Context) error {
	var userLog User
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
