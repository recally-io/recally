package logger

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"vibrain/internal/pkg/config"
)

var defaultLogAttrs = []string{config.ContextKeyRequestID, config.ContextKeyUserID}

func New() *slog.Logger {
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		debug = false
	}
	var logger *slog.Logger
	if debug {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return logger
}

func FromContext(ctx context.Context) *slog.Logger {
	logger := getValFromContext(ctx, config.ContextKeyLogger)
	if logger != nil {
		return logger.(*slog.Logger)
	}
	newLogger := New()
	handler := newLogger.Handler().WithAttrs(buildLogAttrs(ctx))
	newLogger = slog.New(handler)
	return newLogger
}

func getValFromContext(ctx context.Context, key string) interface{} {
	return ctx.Value(config.ContextKey(key))
}

func buildLogAttrs(ctx context.Context) []slog.Attr {
	attrs := make([]slog.Attr, 0)
	for _, key := range defaultLogAttrs {
		if val := getValFromContext(ctx, key); val != nil {
			attrs = append(attrs, slog.Any(key, val))
		}
	}
	return attrs
}
