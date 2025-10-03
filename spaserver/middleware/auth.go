package middleware

import (
	"context"
	"net/http"
	"spahttp/domain"
	"spahttp/repo"
	"spahttp/repo/selfdb/models"

	"github.com/labstack/echo/v4"
)

func (m *Middleware) LoginRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// m.Logger().Debugf("middleware LoginRequired %v", c.Request().URL.Path)
		// If the user is not authenticated, redirect them to the login page and return
		// from the middleware chain so that no subsequent handlers in the chain are
		// executed.
		if !m.IsAuthenticated(c.Request()) {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache (or
		// other intermediary cache).
		c.Response().Header().Add("Cache-Control", "no-store")
		return next(c)
	}
}

// RedirectAuthenticatedUsers is middleware that redirects authenticated users away from login or register pages to /events.
func (m *Middleware) RedirectAuthenticatedUsers(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// m.Logger().Debugf("middleware RedirectAuthenticatedUsers %v", c.Request().URL.Path)
		// Check if the user is authenticated using app's existing method
		if m.IsAuthenticated(c.Request()) {
			// Check if the path is either the login or register page
			if c.Request().URL.Path == "/login" || c.Request().URL.Path == "/" {
				// Redirect authenticated users to the events page
				return c.Redirect(http.StatusSeeOther, "/site/index")
			}
		}
		// If not authenticated or not targeting the restricted pages, continue
		return next(c)
	}
}

// Authenticate is middleware that ensures incoming requests are associated with a valid authenticated user.
func (m *Middleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// m.Logger().Debugf("middleware Authenticate %v", c.Request().URL.Path)
		// Retrieve the authenticatedUserID value from the session using the GetInt()
		// method.
		//This will return the zero value for an int (0) if no
		// "authenticatedUserID" value is in the session -- in which case we call the
		// next handler in the chain as normal and return.
		id := m.SessionManager().GetInt(c.Request().Context(), "authenticatedUserID")
		if id == 0 {
			return next(c)
		}

		// Otherwise, we check to see if a user with that ID exists in our
		// database.
		exists, err := m.userExists(id)
		if err != nil {
			return m.ServerError(c, err)
		}

		// If a matching user is found, we know that the request is coming from an
		// authenticated user who exists in our database.
		//We create a new copy of the
		// request (with an isAuthenticatedContextKey value of true in the request data)
		// and assign it to r.
		// if exists {
		// 	ctx := context.WithValue(c.Request().Context(), domain.IsAuthenticatedContextKey, true)
		// 	c.Request().WithContext(ctx)
		// }
		if exists {
			ctx := context.WithValue(c.Request().Context(), domain.IsAuthenticatedContextKey, true)
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)
		}
		// Call the next handler in the chain.
		return next(c)
	}
}

func (m *Middleware) userExists(id int) (bool, error) {
	repo, err := repo.GetRepository()
	if err != nil {
		return false, err
	}
	dbLock, err := repo.LockOther()
	if err != nil {
		return false, err
	}
	defer repo.UnlockOther(dbLock)
	exists, err := dbLock.UserExists(models.ID, id)
	if err != nil {
		return false, err
	}
	return exists, nil
}
