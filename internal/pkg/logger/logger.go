package logger

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"recally/internal/pkg/contexts"

	slogbetterstack "github.com/samber/slog-betterstack"
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

	handlers := make([]slog.Handler, 0)

	if debug {
		handler := NewColorHandler(&slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: debug,
		}, WithDestinationWriter(os.Stdout), WithColor(), WithOutputEmptyAttrs())

		handlers = append(handlers, handler)
	} else {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: debug,
			Level:     slog.LevelInfo,
		})
		handlers = append(handlers, handler)
	}

	if os.Getenv("BETTER_STACK_SOURCE_TOKEN") != "" {
		handler := slogbetterstack.Option{
			Token:     os.Getenv("BETTER_STACK_SOURCE_TOKEN"),
			AddSource: debug,
		}.NewBetterstackHandler()
		handlers = append(handlers, handler)
	}

	logger = slog.New(NewMultiHandler(handlers...))

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

// CopyContext creates a new context with copied values from the original context
// for the default logging attributes.
//
// Parameters:
//   - ctx: The source context
//
// Returns:
//   - context.Context: A new context with copied values
func CopyContext(ctx context.Context) context.Context {
	newCtx := context.Background()

	for _, key := range defaultLogAttrs {
		val, ok := contexts.Get[any](ctx, key)
		if ok {
			newCtx = contexts.Set(newCtx, key, val)
		}
	}

	return newCtx
}
