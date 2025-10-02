package login

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (t *page) Routes() error {
	t.Echo().GET("", t.Index)
	t.Echo().GET("/login", t.Index)
	return nil
}

func (t *page) Index(c echo.Context) error {
	data, err := t.PageData()
	if err != nil {
		return t.ServerError(c, err)
	}
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("index", data)); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}
