package restserver

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/benchkram/bobc/application"

	"github.com/benchkram/bobc/restserver/generated"
	"github.com/benchkram/errz"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/logrusorgru/aurora"
)

var (
	ErrServerAlreadyStarted  = fmt.Errorf("Server already started")
	ErrTooManyListenConfigs  = fmt.Errorf("Only http or unix config can be set, not both")
	ErrInvalidBase64Encoding = fmt.Errorf("Invalid base64 encoding")

	ErrUnauthorized     = fmt.Errorf("Unauthorized")
	ErrServerNotStarted = fmt.Errorf("Server not started yet")

	ErrInvalidProjectID = fmt.Errorf("Project ID is not in valid format")
)

var (
	defaultHostname  = "0.0.0.0"
	defaultPort      = "8100"
	DefaultUploadDir = filepath.Join(os.TempDir(), "./upload")
)

type Authenticator interface {
	Authenticate(ctx echo.Context) (err error)
}

type S struct {
	// service to perform crud operations with DB
	app application.Application

	// Address to listen to schema: ["host:port"]
	// Only one of address or socketAddress can be set.
	address string

	// hold the http server
	server *http.Server

	// Directory to store uploads.
	// Created if not exists.
	// Defaults to ["/tmp/upload"]
	uploadDir string

	authenticator Authenticator
}

func New(opts ...Option) (s *S, err error) {
	defer errz.Recover(&err)
	s = &S{
		address:   fmt.Sprintf("%s:%s", defaultHostname, defaultPort),
		uploadDir: DefaultUploadDir,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}

	return s, nil
}

var once bool

func (s *S) checkIfApplicationExists(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Path()
		if strings.HasPrefix(path, "/api/") && s.app == nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}

		return next(c)
	}
}

// Start the server. Can be called only once.
func (s *S) Start() (err error) {
	if once {
		return ErrServerAlreadyStarted
	}
	once = true

	errz.Recover(&err)
	e := echo.New()
	e.Debug = true

	// FIXME: Configure logging, NOT to stdout.
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			"authorization",
			"x-ms-useragent",
			"x-ms-client-request-id",
		},
		ExposeHeaders: []string{"Location", "Content-Disposition"},
	}))

	e.Use(s.checkIfApplicationExists)

	// register handlers generated from openapi
	generated.RegisterHandlers(e, s)

	// add cloudflare header to not cache backend routes
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()

			if strings.HasPrefix(path, "/api/") {
				c.Response().Header().Set("CF-Cache-Status", "DYNAMIC")
			}

			return next(c)
		}
	})

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		he, ok := err.(*echo.HTTPError)
		if ok {
			if he.Internal != nil {
				if herr, ok := he.Internal.(*echo.HTTPError); ok {
					he = herr
				}
			}
		} else {
			he = &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			}
		}

		// Issue #1426
		code := he.Code
		message := he.Message
		if m, ok := he.Message.(string); ok {
			if e.Debug {
				message = echo.Map{"message": m, "error": err.Error()}
			} else {
				message = echo.Map{"message": m}
			}
		}

		// Send response
		// err is always being set to not found. Needs to be inspected
		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead { // Issue #608
				err = c.NoContent(code)
			} else {
				err = c.JSON(code, message)
			}
			if err != nil {
				e.Logger.Error(err)
			}
		}
	}

	e.HideBanner = true
	e.HidePort = true

	go func() { aurora.Red(e.Start(s.address)) }()

	fmt.Printf("\n%s\n", aurora.Green("Server started at "+s.address))
	for _, route := range e.Routes() {
		fmt.Printf("[%5v] %s\n", route.Method, route.Path)
	}

	s.server = e.Server
	return nil
}

// Stop the server from the outside
func (s *S) Stop() (err error) {
	if !once {
		return ErrServerNotStarted
	}

	err = s.server.Close()
	if err != nil {
		return err
	}

	once = false
	fmt.Printf("\n%s\n", aurora.Red("Server is shutting down from "+s.address))

	return nil
}

// Address the server adderss (host&port) without the protocol http/https
func (s *S) Address() string {
	return s.address
}
