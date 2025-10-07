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
	csrf, ok := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	if !ok {
		return t.ServerError(c, fmt.Errorf("CSRF token not found in context"))
	}
	data.Csrf = csrf
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
	csrf, ok := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	if !ok {
		return t.ServerError(c, fmt.Errorf("CSRF token not found in context"))
	}
	data.Csrf = csrf
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("homepage", data)); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}
