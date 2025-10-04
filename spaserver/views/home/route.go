package home

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
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
	// тригер ставим в запросах внутри htmx в индексе мы загружаем всю страницу и htmx еще молчит
	// if len(data.errors) > 0 {
	// 	c.Response().Header().Set("HX-Trigger", "alert")
	// }
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("index", data)); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}

func (t *page) HomePage(c echo.Context) error {
	data, err := t.PageModel()
	if err != nil {
		// return t.ServerError(c, err)
		data.errors = append(data.errors, err)
	}
	data.errors = append(data.errors, fmt.Errorf("ошибочка вышлась проверка ..."))
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("homepage", data)); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}
