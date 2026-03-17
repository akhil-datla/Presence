package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/akhil-datla/Presence/internal/auth"
	"github.com/akhil-datla/Presence/internal/config"
	"github.com/akhil-datla/Presence/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server wraps the HTTP server and dependencies.
type Server struct {
	echo  *echo.Echo
	store *store.Store
	cfg   *config.Config
}

// New creates and configures a new server.
func New(cfg *config.Config, s *store.Store, staticFiles embed.FS) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderAuthorization, echo.HeaderContentType},
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${latency_human}\n",
	}))

	jwt := auth.NewJWTService(cfg.JWTSecret)
	registerRoutes(e, s, jwt)

	// Serve embedded static files
	staticFS, _ := fs.Sub(staticFiles, "web/static")
	staticHandler := http.FileServer(http.FS(staticFS))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", staticHandler)))

	// SPA catch-all: serve index.html for any non-API, non-static route
	e.GET("/*", func(c echo.Context) error {
		data, err := staticFiles.ReadFile("web/static/index.html")
		if err != nil {
			return echo.NewHTTPError(500, "frontend not found")
		}
		return c.HTMLBlob(http.StatusOK, data)
	})

	return &Server{echo: e, store: s, cfg: cfg}
}

// Start begins listening for HTTP requests.
func (s *Server) Start() error {
	return s.echo.Start(fmt.Sprintf(":%d", s.cfg.Port))
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.echo.Shutdown(shutdownCtx)
}
