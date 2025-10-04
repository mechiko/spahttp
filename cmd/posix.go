//go:build linux || darwin || freebsd

package main

import (
	"os"
	"spahttp/domain"

	"go.uber.org/zap"
)

func openUrl(_ domain.Apper) {
}

func errMessageExit(loger *zap.SugaredLogger, title string, err error) {
	if loger != nil {
		loger.Errorf("%s %v", title, err)
	}
	os.Exit(-1)
}
