package home

import (
	"spahttp/domain"
	"strings"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

type IServer interface {
	domain.Apper
	Echo() *echo.Echo
	ServerError(c echo.Context, err error) error
	SetActivePage(domain.Model)
	SetFlush(string, string)
	RenderString(name string, data interface{}) (str string, err error)
	Htmx() *htmx.HTMX
}

type page struct {
	IServer
	group           *echo.Group
	model           domain.Model
	defaultTemplate string
	currentTemplate string
	title           string
	name            string
}

func New(app IServer, grp *echo.Group) *page {
	t := &page{
		IServer:         app,
		group:           grp,
		model:           domain.Home,
		defaultTemplate: "index",
		currentTemplate: "index",
		title:           "Нанесение сегодня",
	}
	t.name = strings.ToLower(string(t.model))
	return t
}

// шаблон по умолчанию это на будущее
func (p *page) DefaultTemplate() string {
	return p.defaultTemplate
}

// текущий шаблон это на будущее
func (p *page) CurrentTemplate() string {
	return p.currentTemplate
}

// low caps name
func (p *page) Name() string {
	return p.name
}

func (p *page) Model() domain.Model {
	return p.model
}

// формируем мап для рендера map[string]interface{}{template": .., "data"...}
func (p *page) RenderPageModel(tmpl string, model interface{}) map[string]interface{} {
	return map[string]interface{}{
		"template": tmpl,
		"data":     model,
	}
}

func (p *page) Title() string {
	return p.title
}

// описание вида для меню
func (p *page) Desc() string {
	return "отчеты нанесения"
}

func (p *page) ShowInMenu() bool {
	return true
}
