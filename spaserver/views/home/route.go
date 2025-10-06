package home

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (t *page) Routes() error {
	t.group.GET("/", t.Index)
	t.group.GET("/index", t.Index)
	t.group.GET("/homepage", t.HomePage)
	return nil
}

func (t *page) Index(c echo.Context) error {
	data, err := t.PageModel()
	if err != nil {
		// return t.ServerError(c, err)
		data.errors = append(data.errors, err)
	}
	data.Csrf = c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("index", data)); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}

func (t *page) HomePage(c echo.Context) error {
	data, err := t.PageModel()
	if err != nil {
		data.errors = append(data.errors, err)
	}
	data.errors = append(data.errors, fmt.Errorf("ошибочка вышлась проверка ..."))
	data.Csrf = c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("homepage", data)); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}
