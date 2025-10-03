package login

import (
	"net/http"
	"spahttp/repo"

	"github.com/labstack/echo/v4"
)

func (t *page) Routes() error {
	t.Echo().GET("/", t.Index)
	t.Echo().GET("/login", t.Index)
	t.Echo().POST("/login", t.Login)
	return nil
}

func (t *page) Index(c echo.Context) error {
	data, err := t.PageData()
	if err != nil {
		return t.ServerError(c, err)
	}
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("index", data)); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}

func (t *page) Login(c echo.Context) error {
	// Get form values
	email := c.FormValue("email")
	password := c.FormValue("password")
	email = "kbprime@mail.ru"
	password = "APremote14!"

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
