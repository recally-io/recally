package httpserver

import (
	"context"
	"fmt"
	"os/exec"
	"recally/internal/core/queue"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/config"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/s3"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Service struct {
	Server *echo.Echo
	pool   *db.Pool
	llm    *llms.LLM
	cache  cache.Cache
	uiCmd  *exec.Cmd
	s3     *s3.Client
	queue  *queue.Queue
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

type Option func(*Service)

func WithCache(c cache.Cache) Option {
	return func(s *Service) {
		s.cache = c
	}
}

func WithS3(s3 *s3.Client) Option {
	return func(s *Service) {
		s.s3 = s3
	}
}

func New(pool *db.Pool, llm *llms.LLM, queue *queue.Queue, opts ...Option) (*Service, error) {
	s := &Service{
		Server: echo.New(),
		pool:   pool,
		llm:    llm,
		queue:  queue,
	}
	s.Server.Validator = &CustomValidator{validator: validator.New()}
	s.Server.HTTPErrorHandler = customHTTPErrorHandler
	// if config.Settings.DebugUI {
	// 	logger.Default.Info("debug ui enabled")
	// 	s.uiCmd = exec.Command("bun", "run", "dev")
	// 	s.uiCmd.Dir = "web"
	// 	s.uiCmd.Stdout = os.Stdout
	// 	s.uiCmd.Stderr = os.Stderr
	// }
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
				logger.Default.Fatal("failed to start vite server", "err", err)
			}
		}()
	}

	if err := s.Server.Start(addr); err != nil {
		logger.Default.Fatal("failed to start", "service", s.Name(), "addr", addr, "err", err)
	}
}

func (s *Service) Stop(ctx context.Context) {
	if s.uiCmd != nil {
		if err := s.uiCmd.Process.Signal(syscall.SIGINT); err != nil {
			logger.Default.Fatal("failed to stop vite server", "err", err)
		}
	}

	if err := s.Server.Shutdown(ctx); err != nil {
		logger.Default.Fatal("failed to stop", "service", s.Name(), "err", err)
	}
}

func (s *Service) Name() string {
	return "http server"
}
