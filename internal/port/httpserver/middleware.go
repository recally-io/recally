package httpserver

import (
	"context"
	"log/slog"
	"strings"
	"time"
	"vibrain/internal/pkg/contexts"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/logger"

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
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "custom timeout error message returns to client",
		Timeout:      30 * time.Second,
	}))
	e.Use(contextMiddleWare())
	e.Use(transactionMiddleWare(pool))
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "cookie:token",
		Validator: authValidation,
		Skipper: func(c echo.Context) bool {
			// do not skip auth for /api/ paths
			return !strings.HasPrefix(c.Path(), "/api/")
		},
	}))
}

func authValidation(key string, c echo.Context) (bool, error) {
	setContext(c, contexts.ContextKeyUserID, uuid.NewString())
	// validate key
	// user, err := auth.ValidateJWT(key)
	// if err != nil {
	// 	return false, fmt.Errorf("invalid token: %w", err)
	// }
	// setContext(c, contexts.ContextKeyUserID, user.UserID)
	return true, nil
}

// contextMiddleWare is a middleware that sets logger and other context values to echo.Context
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
		logger.FromContext(c.Request().Context()).Error("recovered from panic", "err", err, "stack", string(stack))
		return err
	}
	return middleware.RecoverWithConfig(cfg)
}

func requestLoggerMiddleware() echo.MiddlewareFunc {
	logger := logger.New()
	logValuesFunc := func(c echo.Context, v middleware.RequestLoggerValues) error {
		attrs := []slog.Attr{
			slog.Time("start_time", v.StartTime),
			slog.Duration("duration", v.Latency),
			slog.String("remote_ip", v.RemoteIP),
			slog.String("method", v.Method),
			slog.String("uri", v.URI),
			slog.String("request_id", v.RequestID),
			slog.String("user_agent", v.UserAgent),
			slog.Int("status", v.Status),
			slog.Any("headers", v.Headers),
		}
		ctx := c.Request().Context()
		if v.Error == nil {
			logger.LogAttrs(ctx, slog.LevelInfo, "REQUEST", attrs...)
		} else {
			attrs = append(attrs, slog.String("error", v.Error.Error()))
			logger.LogAttrs(ctx, slog.LevelError, "REQUEST", attrs...)
		}
		return nil
	}

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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
		HandleError:   true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: logValuesFunc,
	})
}

// contextMiddleWare is a middleware that sets logger and other context values to echo.Context
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

func setContext(c echo.Context, key string, value interface{}) {
	ctx := contexts.Set(c.Request().Context(), key, value)

	// set to context.Context
	c.SetRequest(c.Request().WithContext(ctx))

	// set to echo.Context
	c.Set(key, value)
}
