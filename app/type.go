package app

import (
	"fmt"
	"net/url"
	"strings"

	"spahttp/config"
	"spahttp/domain"
	"time"

	"go.uber.org/zap"
)

type app struct {
	config    *config.Config
	options   *config.Configuration // копия config.Configuration
	logger    *zap.SugaredLogger
	pwd       string
	startTime time.Time
	endTime   time.Time
}

var _ domain.Apper = (*app)(nil)

func New(cfg *config.Config, logger *zap.SugaredLogger, pwd string) *app {
	if logger == nil {
		logger = zap.NewNop().Sugar()
	}
	if cfg == nil {
		panic("nil *config.Config passed to app.New")
	}
	newApp := &app{}
	newApp.pwd = pwd
	newApp.logger = logger
	newApp.config = cfg
	newApp.options = cfg.Configuration()
	newApp.initDateMn()
	return newApp
}

func (a *app) initDateMn() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		a.logger.Warnf("failed to load timezone, using UTC: %v", err)
		loc = time.UTC
	}
	t := time.Now().In(loc)
	year, month, _ := t.Date()
	a.startTime = time.Date(year, month, 1, 0, 0, 0, 0, loc)
	a.endTime = time.Date(year, month+1, 1, 0, 0, 0, 0, loc).Add(-time.Nanosecond)
}

func (a *app) NowDateString() string {
	n := time.Now().Local()
	return fmt.Sprintf("%4d.%02d.%02d %02d:%02d:%02d", n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second())
}

func (a *app) StartDateString() string {
	return fmt.Sprintf("%4d.%02d.%02d", a.startTime.Local().Year(), a.startTime.Local().Month(), a.startTime.Local().Day())
}

func (a *app) EndDateString() string {
	return fmt.Sprintf("%4d.%02d.%02d", a.endTime.Local().Year(), a.endTime.Local().Month(), a.endTime.Local().Day())
}

func (a *app) SetStartDate(d time.Time) {
	a.startTime = d
}

func (a *app) SetEndDate(d time.Time) {
	a.endTime = d
}

func (a *app) StartDate() time.Time {
	return a.startTime
}

func (a *app) EndDate() time.Time {
	return a.endTime
}

func (a *app) Pwd() string {
	return a.pwd
}

func (a *app) Config() *config.Config {
	return a.config
}

func (a *app) Logger() *zap.SugaredLogger {
	return a.logger
}

// выдаем адрес структуры опций программы чтобы править по месту
func (a *app) Options() *config.Configuration {
	return a.options
}

// записываем ключ и его значение только в пакет config
// и Options
// изменения не записываются в файл конфигурации
func (a *app) SetOptions(key string, value any) error {
	a.config.SetInConfig(key, value)
	a.options = a.config.Configuration()
	return nil
}

// записываем файл конфигурации состояние конфигурации
func (a *app) SaveOptions() error {
	if err := a.config.Save(); err != nil {
		return fmt.Errorf("save all in config error %w", err)
	}
	return nil
}

// создаем по необходимости пути программы
func (a *app) CreatePath() error {
	// создаем папку вывода если не пустое значение
	// в папке запуска программы только или если она задана абсолютным значением пути
	// if a.options == nil {
	// 	return fmt.Errorf("опции программы не инициализированы")
	// }
	// if a.options.Output != "" {
	// 	if output, err := createPath(a.options.Output, ""); err != nil {
	// 		return fmt.Errorf("ошибка создания каталога %w", err)
	// 	} else {
	// 		a.options.Output = output
	// 	}
	// 	a.loger.Infof("путь output приложения %s", a.options.Output)
	// }
	return nil
}

func (a *app) ConfigPath() string {
	if a.config != nil {
		return a.config.ConfigPath()
	}
	return ""
}

func (a *app) DbPath() string {
	if a.config != nil {
		return a.config.DbPath()
	}
	return ""
}

func (a *app) LogPath() string {
	if a.config != nil {
		return a.config.LogPath()
	}
	return ""
}

func (a *app) BaseUrl() string {
	host := a.options.Hostname
	if host == "" {
		host = "127.0.0.1"
	}
	port := a.options.HostPort
	u := &url.URL{Scheme: "http"}
	if port != "" {
		u.Host = fmt.Sprintf("%s:%s", host, port)
	} else {
		u.Host = host
	}
	return u.String()
}

func (a *app) Debug() bool {
	mode := strings.ToLower(config.Mode)
	if mode == "development" {
		return true
	}
	return false
}
