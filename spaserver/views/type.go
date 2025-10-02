package views

import (
	"spahttp/domain"

	"github.com/labstack/echo/v4"
)

type IView interface {
	PageData() (interface{}, error)
	Routes() error
	Index(c echo.Context) error
	// имя подшаблона вида по умолчанию
	DefaultTemplate() string
	// имя подшаблона вида текущий
	CurrentTemplate() string
	// строковое значение domain.Model
	Name() string
	// заголовок страницы
	Title() string
	InitData() (interface{}, error)
	Model() domain.Model
	Svg() string
	Desc() string
}
