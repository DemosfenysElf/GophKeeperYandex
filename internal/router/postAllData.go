package router

import (
	"io"
	"net/http"

	"github.com/labstack/echo"

	"PasManagerGophKeeper/internal/service"
)

func (s *serverKeeper) postWrite(c echo.Context) error {
	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}

	user := c.Get("user")
	userID, err := s.DB.GetUserID(c.Request().Context(), user.(string))
	if err != nil {
		return err
	}
	bodyOrder := string(body)

	getType := c.Request().Header.Get(service.Type)
	switch getType {
	case service.Card:
		err = s.DB.WriteCard(c.Request().Context(), bodyOrder, userID)
	case service.Password:
		err = s.DB.WritePassword(c.Request().Context(), bodyOrder, userID)
	case service.Text:
		err = s.DB.WriteText(c.Request().Context(), bodyOrder, userID)
	default:
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}

	c.Response().WriteHeader(http.StatusAccepted)
	return nil
}
