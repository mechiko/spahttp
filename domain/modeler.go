package domain

import (
	"fmt"
	"strings"
)

type Modeler interface {
	Save(Apper) error
	Copy() (Modeler, error) // структура копирует себя и выдает ссылку на копию с массивами и другими данными
	Model() Model           // возвращает тип модели
}

type Model string

const (
	Application Model = "application"
	NoPage      Model = "nopage"
	Header      Model = "header"
	Footer      Model = "footer"
	Index       Model = "index"
	Home        Model = "home"
	Login       Model = "login"
	LogView     Model = "logview"
)

func IsValidModel(s string) bool {
	switch Model(s) {
	case Application, NoPage, Header, Footer, Index, Home, Login, LogView:
		return true
	default:
		return false
	}
}

// строка приводится в нижний регистр потом сравнивается
func ModelFromString(s string) (Model, error) {
	s = strings.ToLower(s)
	switch s {
	case string(Application):
		return Application, nil
	case string(NoPage):
		return NoPage, nil
	case string(Header):
		return Header, nil
	case string(Footer):
		return Footer, nil
	case string(Index):
		return Index, nil
	case string(Home):
		return Home, nil
	case string(Login):
		return Login, nil
	case string(LogView):
		return LogView, nil
	}
	return "", fmt.Errorf("%s ошибочная модель domain.Model", s)
}

func (s Model) String() string {
	return string(s)
}
