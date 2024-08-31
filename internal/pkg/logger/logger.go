package logger

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"vibrain/internal/pkg/contexts"
)

var defaultLogAttrs = []string{contexts.ContextKeyRequestID, contexts.ContextKeyUserID, contexts.ContextKeyUserName}

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
		handler := NewColorHandler(&slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}, WithDestinationWriter(os.Stdout), WithColor(), WithOutputEmptyAttrs())
		logger = slog.New(handler)
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: false,
			Level:     slog.LevelInfo,
		}))
	}
	return Logger{
		Logger: logger,
	}
}

// FromContext returns a logger from context
func FromContext(ctx context.Context, attrs ...slog.Attr) Logger {
	logger, ok := contexts.Get[Logger](ctx, contexts.ContextKeyLogger)
	if ok {
		return logger
	}
	sLogger := New()

	defaultAttrs := buildLogAttrs(ctx)
	handler := sLogger.Handler().WithAttrs(append(defaultAttrs, attrs...))
	return Logger{
		Logger: slog.New(handler),
	}
}

func buildLogAttrs(ctx context.Context) []slog.Attr {
	attrs := make([]slog.Attr, 0)
	for _, key := range defaultLogAttrs {

		val, ok := contexts.Get[any](ctx, key)
		if ok {
			attrs = append(attrs, slog.Any(key, val))
		}
	}
	return attrs
}
