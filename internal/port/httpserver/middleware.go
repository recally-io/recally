package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"strings"
	"time"

	"recally/internal/pkg/config"
	"recally/internal/pkg/contexts"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Service) registerMiddlewares() {
	e := s.Server
	pool := s.pool

	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.Must(uuid.NewV7()).String()
		},
		RequestIDHandler: func(c echo.Context, id string) {
			setContext(c, contexts.ContextKeyRequestID, id)
		},
	}))
	e.Use(requestLoggerMiddleware())
	e.Use(recoverMiddleware())
	e.Use(middleware.CORS())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper: func(c echo.Context) bool {
			return false
		},
		ErrorMessage: "custom timeout error message returns to client",
		Timeout:      60 * time.Second,
	}))
	e.Use(middleware.CORS())
	e.Use(contextMiddleWare())
	e.Use(transactionMiddleWare(pool))
}

// contextMiddleWare is a middleware that sets logger and other context values to echo.Context.
func contextMiddleWare() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			setContext(c, contexts.ContextKeyLogger, logger.FromContext(c.Request().Context()))

			return next(c)
		}
	}
}

func recoverMiddleware() echo.MiddlewareFunc {
	cfg := middleware.DefaultRecoverConfig
	cfg.LogErrorFunc = func(c echo.Context, err error, stack []byte) error {
		if config.Settings.Debug {
			debug.PrintStack()
		}

		logger.FromContext(c.Request().Context()).Error("http request recovered from panic", "err", err, "stack", string(stack))

		return err
	}

	return middleware.RecoverWithConfig(cfg)
}

func requestLoggerMiddleware() echo.MiddlewareFunc {
	logValuesFunc := func(c echo.Context, v middleware.RequestLoggerValues) error {
		attrs := []slog.Attr{
			slog.Time("start_time", v.StartTime),
			slog.Duration("duration", time.Duration(v.Latency.Milliseconds())),
			slog.String("remote_ip", v.RemoteIP),
			slog.String("method", v.Method),
			slog.String("uri", v.URI),
			slog.String("request_id", v.RequestID),
			slog.String("user_agent", v.UserAgent),
			slog.Int("status", v.Status),
			slog.Any("headers", v.Headers),
		}
		ctx := c.Request().Context()
		msg := fmt.Sprintf("HTTP Request - %s %s, Status: %d, Duration: %dms", v.Method, v.URI, v.Status, v.Latency.Milliseconds())

		if v.Error == nil {
			logger.Default.LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
		} else {
			attrs = append(attrs, slog.String("error", v.Error.Error()))
			logger.Default.LogAttrs(ctx, slog.LevelError, msg, attrs...)
		}

		return v.Error
	}

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Request().RequestURI, "/assets/")
		},
		LogStatus:  true,
		LogURI:     true,
		LogError:   true,
		LogLatency: true,
		// LogProtocol:  true,
		LogRemoteIP: true,
		// LogHost:      true,
		LogMethod:     true,
		LogUserAgent:  true,
		LogRequestID:  true,
		LogReferer:    true,
		LogHeaders:    []string{},
		HandleError:   false,
		LogValuesFunc: logValuesFunc,
	})
}

// contextMiddleWare is a middleware that sets logger and other context values to echo.Context.
func transactionMiddleWare(pool *db.Pool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			tx, err := pool.Begin(ctx)
			if err != nil {
				return err
			}

			setContext(c, contexts.ContextKeyTx, tx)

			defer func() {
				if r := recover(); r != nil {
					if err := tx.Rollback(context.Background()); err != nil {
						logger.FromContext(ctx).Error("failed to rollback transaction", "err", err)
					}

					panic(r)
				}
			}()

			if err := next(c); err != nil {
				if err := tx.Rollback(context.Background()); err != nil {
					logger.FromContext(ctx).Error("failed to rollback transaction", "err", err)
				}

				return err
			}

			return tx.Commit(context.Background())
		}
	}
}

func setContext(c echo.Context, key string, value any) {
	ctx := contexts.Set(c.Request().Context(), key, value)

	// set to context.Context
	c.SetRequest(c.Request().WithContext(ctx))

	// set to echo.Context
	c.Set(key, value)
}
