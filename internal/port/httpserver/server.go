package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/logger"
	"vibrain/internal/port/httpserver/handlers"
	"vibrain/web"

	"github.com/labstack/echo/v4"
)

type Service struct {
	Server  *echo.Echo
	Handler *handlers.Handler
}

type Option func(*Service)

func WithCache(c *cache.DbCache) Option {
	return func(s *Service) {
		s.Handler.Cache = c
	}
}

func New(pool *db.Pool, opts ...Option) (*Service, error) {
	handler := handlers.New(pool)

	service := &Service{
		Server: newServer(handler),
	}
	for _, opt := range opts {
		opt(service)
	}
	return service, nil
}

func (s *Service) Start(ctx context.Context) {
	addr := fmt.Sprintf("%s:%d", config.Settings.Service.Host, config.Settings.Service.Port)

	if err := s.Server.Start(addr); err != nil {
		logger.Default.Fatal("failed to start", "service", s.Name(), "addr", addr, "error", err)
	}
}

func (s *Service) Stop(ctx context.Context) {
	if err := s.Server.Shutdown(ctx); err != nil {
		logger.Default.Fatal("failed to stop", "service", s.Name(), "error", err)
	}
}

func (s *Service) Name() string {
	return "http server"
}

func newServer(handler *handlers.Handler) *echo.Echo {
	e := echo.New()
	registerMiddlewares(e)
	registerRouters(e, handler)

	// Health check
	e.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	// static files
	e.GET("/*", echo.WrapHandler(http.FileServer(web.StaticHttpFS)))
	return e
}
