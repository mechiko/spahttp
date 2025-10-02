package header

import (
	"fmt"
	"spahttp/reductor"
)

func (t *page) InitData() (interface{}, error) {
	model, err := NewModel(t)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	for _, m := range t.Menu() {
		page, ok := t.Views()[m]
		if !ok {
			t.Logger().Errorf("menu %s not found", m)
			continue
		}
		menuItem := &MenuItem{
			Name:   page.Name(),
			Title:  page.Title(),
			Active: t.ActivePage() == page.Model(),
			Desc:   page.Desc(),
			Svg:    page.Svg(),
		}
		model.Items = append(model.Items, menuItem)
	}
	reductor.Instance().SetModel(model, false)
	return model, nil
}

func (t *page) PageData() (interface{}, error) {
	model, err := reductor.Instance().Model(t.model)
	return model, err
}

// с преобразованием
func (t *page) PageModel() (*MenuModel, error) {
	model, err := reductor.Instance().Model(t.model)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if mdl, ok := model.(*MenuModel); ok {
		return mdl, nil
	}
	return nil, fmt.Errorf("pagemodel wrong type %T", model)
}

// сброс модели редуктора для страницы
func (t *page) ResetData() {
}
