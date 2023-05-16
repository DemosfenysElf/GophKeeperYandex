package router

import (
	"encoding/hex"
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

	user := c.Get(service.User)
	userID, err := s.DB.GetUserID(c.Request().Context(), user.(string))
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}
	bodyToString := hex.EncodeToString(body)

	path := c.Request().URL.Path
	switch path {
	case service.Write + service.Card:
		err = s.DB.WriteCard(c.Request().Context(), bodyToString, userID)
	case service.Write + service.Password:
		err = s.DB.WritePassword(c.Request().Context(), bodyToString, userID)
	case service.Write + service.Text:
		err = s.DB.WriteText(c.Request().Context(), bodyToString, userID)
	case service.Write + service.Bin:
		err = s.DB.WriteBin(c.Request().Context(), bodyToString, userID)
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
