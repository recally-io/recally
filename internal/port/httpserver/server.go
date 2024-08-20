package httpserver

import (
	"context"
	"fmt"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/llms"
	"vibrain/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

type Service struct {
	Server *echo.Echo
	pool   *db.Pool
	llm    *llms.LLM
	cache  cache.Cache
}

type Option func(*Service)

func WithCache(c cache.Cache) Option {
	return func(s *Service) {
		s.cache = c
	}
}

func New(pool *db.Pool, llm *llms.LLM, opts ...Option) (*Service, error) {
	s := &Service{
		Server: echo.New(),
		pool:   pool,
		llm:    llm,
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.cache == nil {
		s.cache = cache.MemCache
	}

	s.registerMiddlewares()
	s.registerRouters()
	return s, nil
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
