//go:build windows

package main

import (
	"os"
	"spahttp/domain"

	"github.com/mechiko/utility"
	"go.uber.org/zap"
)

func openUrl(app domain.Apper) {
	// utility.OpenHttpLinkInShell(app.BaseUrl())
	utility.OpenHttpBrowser(app.BaseUrl(), utility.Chrome)
}

func errMessageExit(loger *zap.SugaredLogger, title string, err error) {
	if loger != nil {
		loger.Errorf("%s %v", title, err)
	}
	utility.MessageBox(title, err.Error())
	os.Exit(-1)
}
