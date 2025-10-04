package spaserver

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// вызывает рендеринг по имени страницы вида в s.activePage
func (s *Server) Page(c echo.Context) error {
	ap := s.activePage
	if _, ok := s.views[ap]; !ok {
		return s.ServerError(c, fmt.Errorf("not found active page %s", ap))
	}
	view := s.views[ap]
	return view.Index(c)
}
