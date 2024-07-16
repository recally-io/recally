package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()
	registerMiddlewares(e)
	registerRouters(e)

	// Health check
	e.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	return e
}



type HttpService struct {
	e *echo.Echo
}


func NewServer() (*HttpService, error) {
	return &HttpService{
		e: New(),
	}, nil
}

func (s *HttpService) Start(ctx context.Context) {
	addr := fmt.Sprintf(":%d", config.Settings.Port)
	if err := s.e.Start(addr); err != nil {
		logger.Default.Fatal("failed to start", "service", s.Name(), "addr", addr, "error", err)
	}
}

func (s *HttpService) Stop(ctx context.Context) {
	if err := s.e.Shutdown(ctx); err != nil {
		logger.Default.Fatal("failed to stop", "service", s.Name(), "error", err)
	}
}

func (s *HttpService) Name() string {
	return "http server"
}
