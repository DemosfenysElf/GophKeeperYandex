package router

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"

	"PasManagerGophKeeper/internal/service"
)

func (s *serverKeeper) mwAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		headerAuth := c.Request().Header.Get(service.Authorization)
		if headerAuth == "" {
			c.Response().WriteHeader(http.StatusUnauthorized)
			return nil
		}
		headerParts := strings.Split(headerAuth, " ")
		if len(headerParts) != 2 || headerParts[0] != service.Bearer {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return nil
		}

		if len(headerParts[1]) == 0 {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return nil
		}

		claims, err := service.DecodeJWT(headerParts[1])
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return nil
		}

		c.Set("user", claims.Login)
		return next(c)
	}
}
