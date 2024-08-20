package httpserver

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
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
	uiCmd  *exec.Cmd
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
	if config.Settings.DebugUI {
		logger.Default.Info("debug ui enabled")
		s.uiCmd = exec.Command("bun", "run", "dev")
		s.uiCmd.Dir = "web"
		s.uiCmd.Stdout = os.Stdout
		s.uiCmd.Stderr = os.Stderr
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
	if s.uiCmd != nil {
		go func() {
			// run vite server
			logger.Default.Info("starting vite server")
			if err := s.uiCmd.Run(); err != nil {
				logger.Default.Fatal("failed to start vite server", "error", err)
			}
		}()
	}
	if err := s.Server.Start(addr); err != nil {
		logger.Default.Fatal("failed to start", "service", s.Name(), "addr", addr, "error", err)
	}
}

func (s *Service) Stop(ctx context.Context) {
	if s.uiCmd != nil {
		if err := s.uiCmd.Process.Signal(syscall.SIGINT); err != nil {
			logger.Default.Fatal("failed to stop vite server", "error", err)
		}
	}
	if err := s.Server.Shutdown(ctx); err != nil {
		logger.Default.Fatal("failed to stop", "service", s.Name(), "error", err)
	}
}

func (s *Service) Name() string {
	return "http server"
}
