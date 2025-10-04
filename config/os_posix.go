//go:build linux || darwin || freebsd

package config

var (
	dbPath               = "/var/local/spahttp"
	logPath              = "/var/log/spahttp"
	configPath           = "/etc/spahttp"
	Supported            = true
	Linux                = true
	Windows              = false
	PosixUserUIDGUID int = 1002
	PosixChownPath   int = 0755
	PosixChownFile   int = 0644
)
