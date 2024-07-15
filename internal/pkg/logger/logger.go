package logger

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"vibrain/internal/pkg/constant"
)

var defaultLogAttrs = []string{constant.ContextKeyRequestID, constant.ContextKeyUserID}

// Default logger
var Default = New()

// Logger is a wrapper around slog.Logger
type Logger struct {
	*slog.Logger
}

// Debug logs a message at level Fatal on the standard logger.
// it will exit the program after logging
func (l Logger) Fatal(msg string, args ...interface{}) {
	l.Error(msg, args...)
	os.Exit(1)
}

// New creates a new logger
func New() Logger {
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
	return Logger{
		Logger: logger,
	}
}

// FromContext returns a logger from context
func FromContext(ctx context.Context) Logger {
	logger := getValFromContext(ctx, constant.ContextKeyLogger)
	if logger != nil {
		return logger.(Logger)
	}
	sLogger := New()
	handler := sLogger.Handler().WithAttrs(buildLogAttrs(ctx))
	return Logger{
		Logger: slog.New(handler),
	}
}

func getValFromContext(ctx context.Context, key string) interface{} {
	return ctx.Value(constant.ContextKey(key))
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
