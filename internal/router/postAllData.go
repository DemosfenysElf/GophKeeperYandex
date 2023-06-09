package router

import (
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo"

	"PasManagerGophKeeper/internal/service"
)

// postWrite сохранение в бд данных в зависимости от пути роута
func (s *serverKeeper) postWrite(c echo.Context) error {
	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}

	user := c.Get(service.User)
	userID, err := s.DB.GetUserID(c.Request().Context(), user.(string))
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}
	bodyToString := hex.EncodeToString(body)

	path := c.Request().URL.Path
	typeData := strings.Replace(path, "/write/", "", -1)
	err = s.DB.WriteData(c.Request().Context(), bodyToString, userID, typeData)

	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}

	c.Response().WriteHeader(http.StatusAccepted)
	return nil
}
