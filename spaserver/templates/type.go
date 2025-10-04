package templates

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"spahttp/domain"

	"github.com/mechiko/utility"
)

const modError = "http:templates"

const defaultTemplate = "index"

// если каталог ../spaserver/templates существует, то прописываем его в переменную
// для поиска шаблонов динамической обработки для отладки
// вычисляется абсолютный путь относительно каталога запуска в cmd под отладкой, потому ..
var pathTemplates = "../spaserver/templates"

func rootPathTemplates() (out string) {
	defer func() {
		if r := recover(); r != nil {
			out = ""
		}
	}()
	out, err := filepath.Abs(pathTemplates)
	if err != nil {
		return ""
	}
	if utility.PathOrFileExists(out) {
		return out
	}
	return ""
}

type ITemplateUI interface {
	LoadTemplates() (err error)
	Render(w io.Writer, page domain.Model, name string, data interface{}) error
	RenderDebug(w io.Writer, page domain.Model, name string, data interface{}) error
}

type Templates struct {
	domain.Apper
	debug                    bool
	pages                    map[domain.Model]*template.Template
	fs                       fs.FS
	rootPathTemplateGinDebug string
	semaphore                Semaphore
}

var _ ITemplateUI = &Templates{}

// panic if error
func New(app domain.Apper) (*Templates, error) {
	t := &Templates{
		Apper:                    app,
		pages:                    nil,
		rootPathTemplateGinDebug: rootPathTemplates(),
		semaphore:                NewSemaphore(1),
	}
	if t.rootPathTemplateGinDebug != "" {
		// отладка возможна когда путь до шаблонов существует
		t.debug = true
	}
	if err := t.LoadTemplates(); err != nil {
		return nil, fmt.Errorf("%s %w", modError, err)
	}
	return t, nil
}

func (t *Templates) IsDebug() bool {
	return t.debug
}

func (t *Templates) RootPathTemplates() string {
	return t.rootPathTemplateGinDebug
}
