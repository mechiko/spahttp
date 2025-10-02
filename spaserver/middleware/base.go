package middleware

import (
	"net/http"
	"spahttp/domain"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
)

type IApp interface {
	domain.Apper
	// ServerError(w http.ResponseWriter, r *http.Request, err error)
	ServerError(echo.Context, error) error
	IsAuthenticated(r *http.Request) bool
	GetAuthenticatedUserId(r *http.Request) int
	GetAuthenticatedUserName(r *http.Request) string
	SessionManager() *scs.SessionManager
}

type Middleware struct {
	IApp
}

// NewMiddleware initializes and returns a new Middleware instance, associating it with the provided app.App instance.
func NewMiddleware(app IApp) *Middleware {
	return &Middleware{IApp: app}
}
