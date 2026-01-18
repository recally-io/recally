package main

import (
	"log/slog"
	"os"

	"recally/internal/pkg/webreader"
)

// cliLogger implements webreader.Logger using log/slog
type cliLogger struct {
	logger *slog.Logger
}

// NewLogger creates a new CLI logger that implements webreader.Logger
// verbose enables debug-level logging, otherwise logs at info level and above
// All output goes to stderr (proper CLI behavior - stdout for content, stderr for logs)
func NewLogger(verbose bool) webreader.Logger {
	// Determine log level based on verbose flag
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}

	// Create a text handler for human-readable output
	// Remove timestamps for cleaner CLI output
	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time attribute for cleaner CLI output
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}

	handler := slog.NewTextHandler(os.Stderr, opts)
	logger := slog.New(handler)

	return &cliLogger{
		logger: logger,
	}
}

// Info logs an informational message
func (l *cliLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// Error logs an error message
func (l *cliLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

// Compile-time check that cliLogger implements webreader.Logger
var _ webreader.Logger = (*cliLogger)(nil)
