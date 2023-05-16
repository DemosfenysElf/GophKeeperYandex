package router

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"

	"PasManagerGophKeeper/internal/service"
)

func (s *serverKeeper) getReadALL(c echo.Context) error {
	user := c.Get(service.User)
	userID, err := s.DB.GetUserID(c.Request().Context(), user.(string))
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	var data []string
	path := c.Request().URL.Path
	switch path {
	case service.Read + service.Card:
		data, err = s.DB.ReadAllCard(c.Request().Context(), userID)
	case service.Read + service.Password:
		data, err = s.DB.ReadAllPassword(c.Request().Context(), userID)
	case service.Read + service.Text:
		data, err = s.DB.ReadAllText(c.Request().Context(), userID)
	case service.Read + service.Bin:
		data, err = s.DB.ReadAllBin(c.Request().Context(), userID)
	default:
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}

	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}
	var dataByte [][]byte
	for _, datum := range data {
		decodeString, err := hex.DecodeString(datum)
		if err != nil {
			return err
		}
		dataByte = append(dataByte, decodeString)

	}

	marshalData, err := json.Marshal(dataByte)
	if err != nil {
		return err
	}

	c.Response().WriteHeader(http.StatusAccepted)
	c.Response().Write(marshalData)
	return nil
}
