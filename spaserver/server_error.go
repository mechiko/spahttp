package spaserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// выводил лог ошибки и c.NoContent(204)
// в поток событий error sse
func (s *Server) ServerError(c echo.Context, err error) error {
	s.Logger().Errorf("%s server error %м", c.Request().RequestURI, err)
	c.NoContent(204)
	// s.SetFlush(err.Error(), "error")
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	// return err
}
