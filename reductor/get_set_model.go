package reductor

import (
	"fmt"
	"reflect"
	"spahttp/domain"

	"github.com/mechiko/utility"
)

// вернет модель из мап или nil если запрошенной модели нет
// возвращает указатель модели
func (rdc *Reductor) Model(page domain.Model) (interface{}, error) {
	rdc.logger.Debugf("reductor model() %v", page.String())
	rdc.mutex.RLock()
	defer rdc.mutex.RUnlock()

	if pageModel, ok := rdc.models[page]; ok {
		if !utility.IsPointer(pageModel) {
			return nil, fmt.Errorf("reductor internal error model not pointer")
		}
		// prevent returning a typed-nil pointer
		if v := reflect.ValueOf(pageModel); v.Kind() == reflect.Ptr && v.IsNil() {
			return nil, fmt.Errorf("reductor internal error: model pointer is nil for page %s", page)
		}
		return pageModel, nil
	}
	return nil, fmt.Errorf("reductor запрошенной модели нет")
}

// записываем модель по типу енум моделей
// модель должна быть указателем!
// в редукторе модели храним тоже по указателям
// send - извещать в канал о смене состояния (это когда смена состояния в форме которой незачем обновлятся)
func (rdc *Reductor) SetModel(model domain.Modeler, send bool) error {
	rdc.logger.Debugf("reductor setmodel() %v", model.Model())
	rdc.mutex.Lock()
	defer rdc.mutex.Unlock()
	if !utility.IsPointer(model) {
		return fmt.Errorf("reductor: model must be a pointer")
	}
	page := model.Model()
	if !domain.IsValidModel(string(page)) {
		return fmt.Errorf("reductor: model type is invalid")
	}
	// disallow typed-nil pointers
	if v := reflect.ValueOf(model); v.Kind() == reflect.Ptr && v.IsNil() {
		return fmt.Errorf("reductor: model pointer is nil")
	}

	storeModel, err := model.Copy()
	if err != nil {
		return fmt.Errorf("reductor: само копирования модели %w", err)
	}
	if !utility.IsPointer(storeModel) {
		return fmt.Errorf("reductor: model copy must be a pointer")
	}
	// also guard typed-nil copies
	if v := reflect.ValueOf(storeModel); v.Kind() == reflect.Ptr && v.IsNil() {
		return fmt.Errorf("reductor: model copy pointer is nil")
	}
	if rdc.models == nil {
		rdc.models = make(ModelList)
	}
	rdc.models[page] = storeModel
	if !send {
		return nil
	}
	// select-based non-blocking send
	if rdc.outStateChan != nil {
		select {
		case rdc.outStateChan <- page:
		default:
			// channel full—drop this update
		}
	}
	return nil
}
