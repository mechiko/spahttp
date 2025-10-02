package header

import (
	"fmt"
	"net/http"
	"spahttp/domain"

	"github.com/labstack/echo/v4"
)

func (t *page) Routes() error {
	// Serve static and media files under /static/ and /uploads/ path.
	t.Echo().GET("/header", t.Index)
	t.Echo().GET("/header/modal", t.modal)
	t.Echo().GET("/header/:page", t.pager)
	return nil
}

func (t *page) Index(c echo.Context) error {
	data, err := t.PageData()
	if err != nil {
		return t.ServerError(c, err)
	}
	if err := c.Render(http.StatusOK, t.Name(), map[string]interface{}{
		"template": "content",
		"data":     data,
	}); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}

func (t *page) modal(c echo.Context) error {
	model, err := t.InitData()
	if err != nil {
		return t.ServerError(c, err)
	}
	if err := c.Render(http.StatusOK, t.Name(), map[string]interface{}{
		"template": "modal",
		"data":     model,
	}); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}

func (t *page) pager(c echo.Context) error {
	page := c.Param("page")
	if page == "" {
		return t.ServerError(c, fmt.Errorf("pager page empty"))
	}
	model, err := domain.ModelFromString(page)
	if err != nil {
		return t.ServerError(c, err)
	}
	t.SetActivePage(model)
	t.Reload()
	return c.NoContent(http.StatusNoContent)
}
