package home

import (
	"fmt"
	"spahttp/domain"
	"spahttp/spaserver/views"
)

type HomeModel struct {
	views.Secure
	Title  string
	model  domain.Model
	errors []error
}

var _ domain.Modeler = (*HomeModel)(nil)

// создаем модель считываем ее состояние и возвращаем указатель
func NewModel(app domain.Apper) (*HomeModel, error) {
	model := &HomeModel{
		Title: "Домашка",
		model: domain.Home,
		// errors: make([]error, 0),
	}
	if err := model.ReadState(app); err != nil {
		return nil, fmt.Errorf("model prodtools read state %w", err)
	}
	return model, nil
}

// синхронизирует с приложением в сторону приложения из модели редуктора
func (m *HomeModel) SyncToStore(app domain.Apper) (err error) {
	return err
}

// читаем состояние приложения
func (m *HomeModel) ReadState(app domain.Apper) (err error) {
	return nil
}

func (a *HomeModel) Copy() (domain.Modeler, error) {
	dst := *a
	dst.errors = nil
	return &dst, nil
}

func (a *HomeModel) Model() domain.Model {
	return a.model
}

func (a *HomeModel) Save(_ domain.Apper) (err error) {
	return nil
}

func (a *HomeModel) Errors() []error {
	out := make([]error, len(a.errors))
	copy(out, a.errors)
	return out
}
