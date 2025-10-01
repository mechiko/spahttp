package application

import (
	"fmt"
	"spahttp/config"
	"spahttp/domain"
	"spahttp/repo"
)

type Application struct {
	model domain.Model
	Title string
	Debug bool
	Host  string
	Port  string
}

var _ domain.Modeler = (*Application)(nil)

// создаем модель считываем ее состояние и возвращаем указатель
func New(app domain.Apper) (*Application, error) {
	rp, err := repo.GetRepository()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	model := &Application{
		model: domain.Application,
		Title: "Application Title",
	}
	if err := model.ReadState(app, rp); err != nil {
		return nil, fmt.Errorf("model application read state %w", err)
	}
	return model, nil
}

// синхронизирует с приложением в сторону приложения из модели редуктора
func (m *Application) SyncToStore(app domain.Apper) (err error) {
	return nil
}

// читаем состояние приложения
func (m *Application) ReadState(app domain.Apper, _ *repo.Repository) (err error) {
	m.Host = app.Options().Hostname
	m.Port = app.Options().HostPort
	m.Debug = config.Mode == "development"
	return nil
}

func (a *Application) Copy() (interface{}, error) {
	// shallow copy that`s why fields is simple
	dst := *a
	return &dst, nil
}

func (a *Application) Model() domain.Model {
	return a.model
}

func (m *Application) Save(app domain.Apper) (err error) {
	if err := app.SaveOptions(); err != nil {
		return fmt.Errorf("application: save options failed: %w", err)
	}
	return nil
}
