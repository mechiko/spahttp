package login

import (
	"fmt"
	"spahttp/domain"
)

type LoginModel struct {
	Title  string
	model  domain.Model
	errors []error
}

var _ domain.Modeler = (*LoginModel)(nil)

// создаем модель считываем ее состояние и возвращаем указатель
func NewModel(app domain.Apper) (*LoginModel, error) {
	model := &LoginModel{
		Title:  "Авторизация",
		errors: make([]error, 0),
		model:  domain.Login,
	}
	if err := model.ReadState(app); err != nil {
		return nil, fmt.Errorf("model prodtools read state %w", err)
	}
	return model, nil
}

// синхронизирует с приложением в сторону приложения из модели редуктора
func (m *LoginModel) SyncToStore(app domain.Apper) (err error) {
	return err
}

// читаем состояние приложения
func (m *LoginModel) ReadState(app domain.Apper) (err error) {
	return nil
}

func (a *LoginModel) Copy() (domain.Modeler, error) {
	dst := *a
	if a.errors != nil {
		dst.errors = append([]error(nil), a.errors...)
	}
	return &dst, nil
}

func (a *LoginModel) Model() domain.Model {
	return a.model
}

func (a *LoginModel) Save(_ domain.Apper) (err error) {
	return nil
}

func (a *LoginModel) Errors() []error {
	out := make([]error, len(a.errors))
	copy(out, a.errors)
	return out
}
