package httpserver

import (
	"time"
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/logger"

	"log/slog"
	"vibrain/internal/port/httpserver/contexts"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func registerMiddlewares(e *echo.Echo) {
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.Must(uuid.NewV7()).String()
		},
		RequestIDHandler: func(c echo.Context, id string) {
			contexts.Set(c, config.ContextKeyRequestID, id)
		},
	}))
	e.Use(requestLoggerMiddleware())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "custom timeout error message returns to client",
		Timeout:      30 * time.Second,
	}))
	e.Use(contextMiddleWare())
}

// contextMiddleWare is a middleware that sets logger and other context values to echo.Context
func contextMiddleWare() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			contexts.Set(c, config.ContextKeyLogger, logger.FromContext(c.Request().Context()))
			return next(c)
		}
	}
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
