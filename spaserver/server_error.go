package spaserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// обработка ошибки
func (s *Server) ServerError(c echo.Context, err error) error {
	s.Logger().Errorf("%s server error %v", c.Request().RequestURI, err)
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}
