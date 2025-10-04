package domain

import (
	"spahttp/config"

	"go.uber.org/zap"
)

type Apper interface {
	Options() *config.Configuration
	SetOptions(key string, value interface{}) error
	SaveOptions() error
	Logger() *zap.SugaredLogger
	ConfigPath() string
	DbPath() string
	LogPath() string
	BaseUrl() string
}
