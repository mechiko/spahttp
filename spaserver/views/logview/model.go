package logview

import (
	"fmt"
	"path/filepath"
	"spahttp/domain"
	"spahttp/spaserver/views"
)

type LogLineEcho struct {
	Lvl          string `json:"lvl"`
	Ts           string `json:"ts"`
	Message      string `json:"message"`
	Proto        string `json:"proto"`
	Host         string `json:"host"`
	Method       string `json:"method"`
	Status       int    `json:"status"`
	ResponseSize int    `json:"response_size"`
	Latency      string `json:"latency"`
	StatusText   string `json:"status_text"`
	ClientIP     string `json:"client_ip"`
	UserAgent    string `json:"user_agent"`
	Path         string `json:"path"`
}

type LogViewModel struct {
	views.Secure
	Title    string
	Start    int
	FileName string
	Lines    []*LogLineEcho
	PerPage  int
	model    domain.Model
	errors   []error
}

var _ domain.Modeler = (*LogViewModel)(nil)

// создаем модель считываем ее состояние и возвращаем указатель
func NewModel(app domain.Apper) (*LogViewModel, error) {
	model := &LogViewModel{
		Title:    "Домашка",
		Start:    0,
		PerPage:  20,
		model:    domain.LogView,
		FileName: filepath.Join(app.LogPath(), "reductor"),
	}
	if err := model.ReadState(app); err != nil {
		return nil, fmt.Errorf("model logview read state %w", err)
	}
	return model, nil
}

// синхронизирует с приложением в сторону приложения из модели редуктора
func (m *LogViewModel) SyncToStore(app domain.Apper) (err error) {
	return err
}

// читаем состояние приложения
func (m *LogViewModel) ReadState(app domain.Apper) (err error) {
	return nil
}

func (a *LogViewModel) Copy() (domain.Modeler, error) {
	dst := *a
	dst.Lines = nil
	dst.errors = nil
	return &dst, nil
}

func (a *LogViewModel) Model() domain.Model {
	return a.model
}

func (a *LogViewModel) Save(_ domain.Apper) (err error) {
	return nil
}

func (a *LogViewModel) Errors() []error {
	out := make([]error, len(a.errors))
	copy(out, a.errors)
	return out
}
