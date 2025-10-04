package header

import (
	"fmt"
	"spahttp/domain"
)

type MenuModel struct {
	Title string
	Items MenuItemSlice
	model domain.Model
}

type MenuItem struct {
	Name   string
	Title  string
	Active bool
	Desc   string
	Svg    string
}

type MenuItemSlice []*MenuItem

var _ domain.Modeler = (*MenuModel)(nil)

// создаем модель считываем ее состояние и возвращаем указатель
func NewModel(app domain.Apper) (*MenuModel, error) {
	model := &MenuModel{
		model: domain.Header,
		Title: "Меню",
		Items: make(MenuItemSlice, 0),
	}
	if err := model.ReadState(app); err != nil {
		return nil, fmt.Errorf("model MenuModel read state %w", err)
	}
	return model, nil
}

// синхронизирует с приложением в сторону приложения из модели редуктора
func (m *MenuModel) SyncToStore(app domain.Apper) (err error) {
	return err
}

// читаем состояние приложения
func (m *MenuModel) ReadState(app domain.Apper) (err error) {
	return nil
}

func (a *MenuModel) Copy() (domain.Modeler, error) {
	dst := *a
	if a.Items != nil {
		dst.Items = make(MenuItemSlice, len(a.Items))
		for i, it := range a.Items {
			if it != nil {
				v := *it
				dst.Items[i] = &v
			}
		}
	}
	return &dst, nil
}

func (a *MenuModel) Model() domain.Model {
	return a.model
}

func (a *MenuModel) Save(_ domain.Apper) (err error) {
	return nil
}
