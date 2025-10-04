package templates

import (
	"fmt"
	"io"
	"spahttp/domain"
)

// page имя страницы из массива структур шаблонов по имени каталога view
// name имя шаблона для вида страницы из имени файла в каталоге view без расширения
// Render processes and renders an HTML template with the passed data and HTTP status code to the response writer.
func (t *Templates) Render(w io.Writer, pageType domain.Model, name string, data interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic %v", r)
		}
	}()
	if name == "" {
		name = defaultTemplate
		t.Logger().Errorf("%s page template is empty", modError)
		// return fmt.Errorf("%s page template is empty", modError)
	}
	tmpl, ok := t.pages[pageType]
	if !ok {
		return fmt.Errorf("template %v not found", pageType)
	}
	return tmpl.ExecuteTemplate(w, name, data)
}

func (t *Templates) RenderDebug(w io.Writer, pageType domain.Model, name string, data interface{}) (err error) {
	t.semaphore.Acquire()
	defer t.semaphore.Release()
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic %v", r)
		}
	}()

	pages, err := t.DynLoadTemplates()
	if err != nil {
		return err
	}
	tmpl, ok := pages[pageType]
	if !ok {
		return fmt.Errorf("template %v not found (name=%s)", pageType, name)
	}

	return tmpl.ExecuteTemplate(w, name, data)
}
