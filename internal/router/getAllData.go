package router

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

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
	typeData := strings.Replace(path, "/read/", "", -1)
	data, err = s.DB.ReadAllDataType(c.Request().Context(), userID, typeData)

	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}
	var dataByte [][]byte
	for _, datum := range data {
		decodeString, err := hex.DecodeString(datum)
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}
		dataByte = append(dataByte, decodeString)
	}

	marshalData, err := json.Marshal(dataByte)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	c.Response().WriteHeader(http.StatusAccepted)
	c.Response().Write(marshalData)
	return nil
}
