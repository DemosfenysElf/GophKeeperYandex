package router

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"

	"PasManagerGophKeeper/internal/service"
)

func (s *serverKeeper) getReadALL(c echo.Context) error {
	user := c.Get(service.User)
	userID, err := s.DB.GetUserID(c.Request().Context(), user.(string))
	if err != nil {
		return err
	}

	var data interface{}
	getType := c.Request().Header.Get(service.Type)
	switch getType {
	case service.Card:
		data, err = s.DB.ReadAllCard(c.Request().Context(), userID)
	case service.Password:
		data, err = s.DB.ReadAllPassword(c.Request().Context(), userID)
	case service.Text:
		data, err = s.DB.ReadAllText(c.Request().Context(), userID)
	default:
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}

	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}

	marshalData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	c.Response().WriteHeader(http.StatusAccepted)
	c.Response().Write(marshalData)
	return nil
}