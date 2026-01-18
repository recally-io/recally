package main

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"

	"recally/internal/pkg/webreader"
)

// TestLoggerImplementsInterface verifies at compile time that cliLogger implements webreader.Logger
func TestLoggerImplementsInterface(t *testing.T) {
	var _ webreader.Logger = NewLogger(false)
}

// TestNewLogger verifies that NewLogger creates a valid logger
func TestNewLogger(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
	}{
		{
			name:    "normal mode",
			verbose: false,
		},
		{
			name:    "verbose mode",
			verbose: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.verbose)
			if logger == nil {
				t.Fatal("NewLogger returned nil")
			}

			// Verify it's the correct type
			if _, ok := logger.(*cliLogger); !ok {
				t.Errorf("NewLogger returned unexpected type: %T", logger)
			}
		})
	}
}

// TestLoggerOutput verifies that logger writes to stderr
func TestLoggerOutput(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w

	// Create logger and log a message
	logger := NewLogger(false)
	logger.Info("test message", "key", "value")

	// Restore stderr and read captured output
	w.Close()
	os.Stderr = oldStderr
	
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatal(err)
	}
	output := buf.String()

	// Verify output contains the message
	if !strings.Contains(output, "test message") {
		t.Errorf("expected output to contain 'test message', got: %s", output)
	}

	// Verify output contains the key-value pair
	if !strings.Contains(output, "key=value") {
		t.Errorf("expected output to contain 'key=value', got: %s", output)
	}

	// Verify output does NOT contain timestamp (we remove it for CLI)
	// slog.TextHandler formats time as "time=2006-01-02T15:04:05.000Z07:00"
	if strings.Contains(output, "time=") {
		t.Errorf("expected output to NOT contain timestamp, got: %s", output)
	}
}

// TestLoggerVerboseMode verifies that verbose mode enables debug logging
func TestLoggerVerboseMode(t *testing.T) {
	tests := []struct {
		name       string
		verbose    bool
		logFunc    func(webreader.Logger)
		shouldLog  bool
	}{
		{
			name:    "info level in normal mode",
			verbose: false,
			logFunc: func(l webreader.Logger) {
				l.Info("info message")
			},
			shouldLog: true,
		},
		{
			name:    "error level in normal mode",
			verbose: false,
			logFunc: func(l webreader.Logger) {
				l.Error("error message")
			},
			shouldLog: true,
		},
		{
			name:    "info level in verbose mode",
			verbose: true,
			logFunc: func(l webreader.Logger) {
				l.Info("info message")
			},
			shouldLog: true,
		},
		{
			name:    "error level in verbose mode",
			verbose: true,
			logFunc: func(l webreader.Logger) {
				l.Error("error message")
			},
			shouldLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr
			oldStderr := os.Stderr
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			os.Stderr = w

			// Create logger and log
			logger := NewLogger(tt.verbose)
			tt.logFunc(logger)

			// Restore stderr and read captured output
			w.Close()
			os.Stderr = oldStderr
			
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Fatal(err)
			}
			output := buf.String()

			hasOutput := len(output) > 0
			if hasOutput != tt.shouldLog {
				t.Errorf("expected shouldLog=%v, got output: %s", tt.shouldLog, output)
			}
		})
	}
}

// TestLoggerLevelFiltering verifies that normal mode filters out debug messages
// Note: Since webreader.Logger interface only has Info and Error methods,
// we can't directly test Debug logging. This test documents the expected behavior
// that normal mode uses Info level and verbose mode uses Debug level.
func TestLoggerLevelFiltering(t *testing.T) {
	// Test that we can access the underlying slog logger and verify its level
	tests := []struct {
		name          string
		verbose       bool
		expectedLevel slog.Level
	}{
		{
			name:          "normal mode uses Info level",
			verbose:       false,
			expectedLevel: slog.LevelInfo,
		},
		{
			name:          "verbose mode uses Debug level",
			verbose:       true,
			expectedLevel: slog.LevelDebug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.verbose)
			cliLog, ok := logger.(*cliLogger)
			if !ok {
				t.Fatal("expected *cliLogger type")
			}

			// Verify the logger is enabled at the expected level
			if !cliLog.logger.Enabled(nil, tt.expectedLevel) {
				t.Errorf("expected logger to be enabled at level %v", tt.expectedLevel)
			}

			// Verify behavior: if normal mode, should NOT be enabled at Debug level
			if !tt.verbose && cliLog.logger.Enabled(nil, slog.LevelDebug) {
				t.Error("expected normal mode to NOT be enabled at Debug level")
			}
		})
	}
}

// TestLoggerErrorOutput verifies that Error method works correctly
func TestLoggerErrorOutput(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w

	// Create logger and log an error
	logger := NewLogger(false)
	logger.Error("error occurred", "error", "something went wrong", "code", 500)

	// Restore stderr and read captured output
	w.Close()
	os.Stderr = oldStderr
	
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatal(err)
	}
	output := buf.String()

	// Verify output contains error level indicator
	if !strings.Contains(output, "level=ERROR") {
		t.Errorf("expected output to contain 'level=ERROR', got: %s", output)
	}

	// Verify output contains the error message
	if !strings.Contains(output, "error occurred") {
		t.Errorf("expected output to contain 'error occurred', got: %s", output)
	}

	// Verify output contains the structured attributes
	if !strings.Contains(output, "error=\"something went wrong\"") {
		t.Errorf("expected output to contain error attribute, got: %s", output)
	}

	if !strings.Contains(output, "code=500") {
		t.Errorf("expected output to contain code attribute, got: %s", output)
	}
}

// TestLoggerNoTimestamp verifies that timestamps are removed from output
func TestLoggerNoTimestamp(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w

	// Create logger and log multiple messages
	logger := NewLogger(true)
	logger.Info("first message")
	logger.Error("second message")

	// Restore stderr and read captured output
	w.Close()
	os.Stderr = oldStderr
	
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatal(err)
	}
	output := buf.String()

	// Verify no timestamp patterns exist
	// slog.TextHandler uses "time=" prefix for timestamps
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		if strings.Contains(line, "time=") {
			t.Errorf("line %d should not contain timestamp, got: %s", i+1, line)
		}
	}
}
