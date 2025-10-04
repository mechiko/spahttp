package reductor

import (
	"fmt"
	"strings"
)

// этот тип для ведения списка всех моделей
// просто запоминать по строке сложно вычислить ошибки в программе
// если строка вдруг окажется не такой как планировалось
type ModelType int

const (
	TrueClient ModelType = iota
	Home
	Application
	Header
	Footer
	Setup
	Index
)

// имена модели используются так же в роутинге там они выступают в качестве имен вида
// должны в роутере приводится к нижнему регистру
func (s ModelType) String() string {
	switch s {
	case Home:
		return "home"
	case TrueClient:
		return "trueclient"
	case Application:
		return "application"
	case Header:
		return "header"
	case Footer:
		return "footer"
	case Setup:
		return "setup"
	case Index:
		return "index"
	default:
		return "неизвестная"
	}
}

// строка приводится в нижний регистр потом сравнивается
func ModelTypeFromString(s string) (ModelType, error) {
	s = strings.ToLower(s)
	switch s {
	case "home":
		return Home, nil
	case "trueclient":
		return TrueClient, nil
	case "application":
		return Application, nil
	case "header":
		return Header, nil
	case "footer":
		return Footer, nil
	case "setup":
		return Setup, nil
	case "index":
		return Index, nil
	}
	return 0, fmt.Errorf("unknown model type: %s", s)
}
