package login

import (
	"net/http"
	"spahttp/repo"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (t *page) Routes() error {
	t.Echo().GET("/", t.Index)
	t.Echo().GET("/login", t.Index)
	t.Echo().POST("/login", t.Login)
	t.Echo().POST("/logout", t.Logout)
	return nil
}

func (t *page) Index(c echo.Context) error {
	data, err := t.PageModel()
	if err != nil {
		return t.ServerError(c, err)
	}
	data.Csrf = c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("index", data)); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}

func (t *page) Logout(c echo.Context) error {
	// Renew session token to prevent session fixation attacks
	err := t.SessionManager().RenewToken(c.Request().Context())
	if err != nil {
		return t.ServerError(c, err)
	}
	// Remove authenticated user id
	t.SessionManager().Remove(c.Request().Context(), "authenticatedUserID")
	return c.Redirect(http.StatusSeeOther, "/")
}

func (t *page) Login(c echo.Context) error {
	// Get form values
	email := c.FormValue("email")
	password := c.FormValue("password")
	email = "kbprime@mail.ru"
	password = "APremote14!"
	hds := c.Request().Header
	t.Logger().Debug(len(hds))
	repo, err := repo.GetRepository()
	if err != nil {
		return t.ServerError(c, err)
	}
	dbLock, err := repo.LockOther()
	if err != nil {
		return t.ServerError(c, err)
	}
	defer repo.UnlockOther(dbLock)

	usr, err := dbLock.Authenticate(email, password)
	if err != nil {
		return t.ServerError(c, err)
	}

	err = t.SessionManager().RenewToken(c.Request().Context())
	if err != nil {
		return t.ServerError(c, err)
	}

	t.SessionManager().Put(c.Request().Context(), "authenticatedUserID", usr.Id)
	t.SessionManager().Put(c.Request().Context(), "authenticatedUsername", usr.Email)

	return c.Redirect(http.StatusSeeOther, "/")
}
