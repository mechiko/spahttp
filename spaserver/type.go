package spaserver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"spahttp/domain"
	"spahttp/embedded"
	"spahttp/spaserver/middleware"
	"spahttp/spaserver/session"
	"spahttp/spaserver/templates"
	"spahttp/spaserver/views"
	"spahttp/zaplog/zap4echo"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	emiddle "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

const (
	// _defaultReadTimeout     = 5 * time.Second
	// _defaultWriteTimeout    = 5 * time.Second
	_defaultAddr            = "127.0.0.1:8888"
	_defaultShutdownTimeout = 5 * time.Second
)

// Server -.
type Server struct {
	domain.Apper
	addr            string
	server          *echo.Echo
	notify          chan error
	shutdownTimeout time.Duration
	sessionManager  *scs.SessionManager
	debug           bool
	templates       *templates.Templates
	views           map[domain.Model]views.IView
	menu            []domain.Model
	activePage      domain.Model
	defaultPage     string
	htmx            *htmx.HTMX
	protected       *echo.Group // для защищенной области роутинга
	// dynamic         *echo.Group // для открытой области роутинга в данной реализации это только /login /sse
}

// var sseManager *sse.Server

func New(a domain.Apper, zl *zap.Logger, port string, debug bool) (ss *Server, err error) {
	addr := fmt.Sprintf("%s:%s", "127.0.0.1", port)
	if port == "" {
		addr = _defaultAddr
	}
	sess := scs.New()
	sess.Lifetime = 24 * time.Hour
	e := echo.New()
	e.Logger.SetOutput(io.Discard)
	e.Use(
		session.LoadAndSave(sess, a.Logger()),
		zap4echo.Logger(zl),
		zap4echo.Recover(zl),
	)
	e.Use(emiddle.CORSWithConfig(emiddle.CORSConfig{
		AllowOrigins: []string{"*"},
		// AllowHeaders:     []string{"authorization", "Content-Type"},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization, echo.HeaderXCSRFToken},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost},
		// AllowMethods:     []string{echo.OPTIONS, echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	// https://www.compilenrun.com/docs/framework/echo/echo-middleware/echo-security-middleware
	e.Use(emiddle.StaticWithConfig(emiddle.StaticConfig{
		Skipper:    func(c echo.Context) bool { return c.Path() == "/" },
		HTML5:      true,
		Root:       "root", // because files are located in `root` directory
		Filesystem: http.FS(embedded.Root),
		// Index:      "",
	}))

	e.Use(emiddle.SecureWithConfig(emiddle.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "DENY",
		ReferrerPolicy:     "no-referrer",
	}))

	// More customized rate limiting
	e.Use(emiddle.RateLimiterWithConfig(emiddle.RateLimiterConfig{
		Skipper: emiddle.DefaultSkipper,
		Store: emiddle.NewRateLimiterMemoryStoreWithConfig(
			emiddle.RateLimiterMemoryStoreConfig{
				Rate:      10,              // 10 requests per time unit
				Burst:     30,              // Allow bursts of up to 30 requests
				ExpiresIn: 1 * time.Minute, // Time unit (1 minute)
			},
		),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			// Use IP address as identifier
			return c.RealIP(), nil
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"message": "Too many requests, please try again later.",
			})
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"message": "Rate limit exceeded.",
			})
		},
	}))

	e.Use(emiddle.CSRFWithConfig(emiddle.CSRFConfig{
		// TokenLookup:    "form:_csrf",
		TokenLookup:    "header:X-CSRF-Token,form:_csrf",
		CookieName:     "_csrf",
		CookieMaxAge:   3600,
		CookieSecure:   true,
		CookieHTTPOnly: true,
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			p := c.Request().URL.Path
			if strings.HasPrefix(p, "/assets/") || strings.HasSuffix(p, ".css") || strings.HasSuffix(p, ".js") {
				c.Response().Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			return next(c)
		}
	})

	ss = &Server{
		Apper:           a,
		addr:            addr,
		server:          e,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
		debug:           debug,
		sessionManager:  sess,
		views:           make(map[domain.Model]views.IView), // массив видов по нему находим шаблоны для рендера
		menu:            make([]domain.Model, 0),
		defaultPage:     "",
		activePage:      domain.NoPage,
		htmx:            htmx.New(),
	}

	e.Renderer = ss
	if ss.templates, err = templates.New(ss); err != nil {
		return nil, fmt.Errorf("spaserver templates error %w", err)
	}
	mdl := middleware.NewMiddleware(ss)
	e.Use(mdl.Authenticate, mdl.RedirectAuthenticatedUsers)
	// append to (session.LoadAndSave(ss.sessionManager, ss.Logger()), mdl.Authenticate, mdl.RedirectAuthenticatedUsers)
	ss.protected = e.Group("site", mdl.LoginRequired)
	ss.protected.Use(emiddle.CORSWithConfig(emiddle.CORSConfig{
		AllowOrigins: []string{"*"},
		// AllowHeaders:     []string{"authorization", "Content-Type"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
		MaxAge:           3600,
	}))
	if err := ss.Routes(); err != nil {
		return nil, fmt.Errorf("spaserver new routes error %w", err)
	}
	return ss, nil
}

func (s *Server) Start() {
	go func() {
		// s.notify <- s.server.StartTLS(s.addr, "cert.pem", "key.pem")
		s.notify <- s.server.Start(s.addr)
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Handler() http.Handler {
	return s.server
}

func (s *Server) SessionManager() *scs.SessionManager {
	return s.sessionManager
}

func (s *Server) Echo() *echo.Echo {
	return s.server
}

func (s *Server) SetActivePage(p domain.Model) {
	s.activePage = p
}

func (s *Server) ActivePage() domain.Model {
	return s.activePage
}

func (s *Server) Views() map[domain.Model]views.IView {
	return s.views
}

func (s *Server) Reload() {
}

func (s *Server) Htmx() *htmx.HTMX {
	return s.htmx
}

func (s *Server) Menu() []domain.Model {
	return s.menu
}

func (s *Server) TemplateIsDebug() bool {
	if s.templates == nil {
		return false
	}
	return s.templates.IsDebug()
}

func (s *Server) RootPathTemplates() string {
	if s.templates == nil {
		return ""
	}
	return s.templates.RootPathTemplates()
}
