package templates

import (
	"embed"
)

//go:embed footer/*.html header/*.html index/*.html login/*.html
var root embed.FS
