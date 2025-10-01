//go:build linux || darwin || freebsd

package config

// import "github.com/mechiko/telebot_v4/internal/entity"
// if !entity.supported
var (
	dbPath               = "/var/local/telebot"
	logPath              = "/var/log/telebot"
	configPath           = "/etc/telebot"
	Supported            = true
	Linux                = true
	Windows              = false
	PosixUserUIDGUID int = 1002
	PosixChownPath   int = 0755
	PosixChownFile   int = 0644
)
