package home

import (
	"fmt"
	"spahttp/reductor"
)

func (t *page) InitData() (interface{}, error) {
	model, err := NewModel(t)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	reductor.Instance().SetModel(model, false)
	return model, nil
}

func (t *page) PageData() (interface{}, error) {
	model, err := reductor.Instance().Model(t.model)
	if mdl, ok := model.(*HomeModel); ok {
		mdl.errors = make([]error, 0)
		return mdl, nil
	}
	return model, err
}

// с преобразованием
func (t *page) PageModel() (*HomeModel, error) {
	model, err := reductor.Instance().Model(t.model)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if mdl, ok := model.(*HomeModel); ok {
		mdl.errors = make([]error, 0)
		return mdl, nil
	}

	return nil, fmt.Errorf("pagemodel ProdUtilModel wrong type %T", model)
}

// сброс модели редуктора для страницы
func (t *page) ResetData() {
}
