package main

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestSanitizeTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple title",
			input:    "Simple Title",
			expected: "simple-title",
		},
		{
			name:     "title with special characters",
			input:    "Title with Special! Characters@",
			expected: "title-with-special-characters",
		},
		{
			name:     "title with multiple spaces",
			input:    "Title   with    multiple    spaces",
			expected: "title-with-multiple-spaces",
		},
		{
			name:     "title with leading/trailing whitespace",
			input:    "  Title with whitespace  ",
			expected: "title-with-whitespace",
		},
		{
			name:     "title with unicode characters",
			input:    "Café and 北京",
			expected: "café-and-北京",
		},
		{
			name:     "title with emojis",
			input:    "Title with emojis 🎉 🚀",
			expected: "title-with-emojis",
		},
		{
			name:     "title with only special characters",
			input:    "!@#$%^&*()",
			expected: "", // Will be replaced by untitled-{timestamp}
		},
		{
			name:     "empty string",
			input:    "",
			expected: "", // Will be replaced by untitled-{timestamp}
		},
		{
			name:     "title with hyphens",
			input:    "Title-with-hyphens",
			expected: "title-with-hyphens",
		},
		{
			name:     "title with multiple consecutive hyphens",
			input:    "Title---with---hyphens",
			expected: "title-with-hyphens",
		},
		{
			name:     "title with leading/trailing hyphens",
			input:    "-Title-",
			expected: "title",
		},
		{
			name:     "title with mixed case",
			input:    "UPPERCASE lowercase MixedCase",
			expected: "uppercase-lowercase-mixedcase",
		},
		{
			name:     "title with numbers",
			input:    "Title with 123 numbers",
			expected: "title-with-123-numbers",
		},
		{
			name:     "title with dots and commas",
			input:    "Title, with. dots, and commas.",
			expected: "title-with-dots-and-commas",
		},
		{
			name:     "title with parentheses",
			input:    "Title (with parentheses)",
			expected: "title-with-parentheses",
		},
		{
			name:     "title with quotes",
			input:    `Title "with quotes"`,
			expected: "title-with-quotes",
		},
		{
			name:     "title with slashes",
			input:    "Title with / slashes \\ backslashes",
			expected: "title-with-slashes-backslashes",
		},
		{
			name:     "very long title (>200 chars)",
			input:    strings.Repeat("a", 250),
			expected: "", // Will be truncated with MD5 hash
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeTitle(tt.input)

			// Special handling for empty results (timestamp fallback)
			if tt.expected == "" && strings.HasPrefix(result, "untitled-") {
				// Verify it has the correct format: untitled-{timestamp}
				parts := strings.Split(result, "-")
				if len(parts) != 2 || parts[0] != "untitled" {
					t.Errorf("expected untitled-{timestamp} format, got %q", result)
				}
				return
			}

			// Special handling for very long titles (MD5 hash appended)
			if len(tt.input) > 200 && len(result) > 200 {
				// Verify truncation and hash appended
				if !strings.Contains(result, "-") {
					t.Errorf("expected MD5 hash appended to truncated title, got %q", result)
				}
				// Check that the result starts with the expected pattern
				expectedPrefix := SanitizeTitle(tt.input[:200])
				if !strings.HasPrefix(result, expectedPrefix[:200]) {
					t.Errorf("truncated title doesn't match expected prefix")
				}
				return
			}

			if result != tt.expected {
				t.Errorf("SanitizeTitle(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeTitleTruncation(t *testing.T) {
	// Test that long titles are properly truncated with MD5 hash
	longTitle := strings.Repeat("a", 250)
	result := SanitizeTitle(longTitle)

	// Should be truncated to 200 chars + hyphen + 8 char hash = 209 chars
	expectedLength := 200 + 1 + 8 // title + hyphen + hash
	if len(result) != expectedLength {
		t.Errorf("expected length %d, got %d", expectedLength, len(result))
	}

	// Verify hash is correct
	hash := md5.Sum([]byte(strings.Repeat("a", 250)))
	expectedHash := fmt.Sprintf("%x", hash)[:8]
	if !strings.HasSuffix(result, expectedHash) {
		t.Errorf("expected hash suffix %s, got %s", expectedHash, result[len(result)-8:])
	}

	// Verify the truncated part
	if !strings.HasPrefix(result, strings.Repeat("a", 200)) {
		t.Errorf("truncated part doesn't match expected prefix")
	}
}

func TestSanitizeTitleUnicode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Chinese characters",
			input:    "文章标题",
			expected: "文章标题",
		},
		{
			name:     "Japanese characters",
			input:    "記事のタイトル",
			expected: "記事のタイトル",
		},
		{
			name:     "Arabic characters",
			input:    "عنوان المقال",
			expected: "عنوان-المقال",
		},
		{
			name:     "Mixed scripts",
			input:    "Title 标题 عنوان",
			expected: "title-标题-عنوان",
		},
		{
			name:     "Cyrillic characters",
			input:    "Заголовок статьи",
			expected: "заголовок-статьи",
		},
		{
			name:     "Greek characters",
			input:    "Τίτλος άρθρου",
			expected: "τίτλος-άρθρου",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeTitle(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeTitle(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetOutputDir(t *testing.T) {
	testDate := time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC)

	t.Run("with custom directory", func(t *testing.T) {
		// Create a temporary directory for testing
		tmpDir := t.TempDir()
		customDir := filepath.Join(tmpDir, "custom")

		dir, err := GetOutputDir(customDir, testDate)
		if err != nil {
			t.Fatalf("GetOutputDir failed: %v", err)
		}

		expected := filepath.Join(customDir, "contents", "2026-01-18")
		if dir != expected {
			t.Errorf("GetOutputDir() = %q, want %q", dir, expected)
		}

		// Verify directory was created
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("directory was not created: %s", dir)
		}
	})

	t.Run("with XDG default", func(t *testing.T) {
		dir, err := GetOutputDir("", testDate)
		if err != nil {
			t.Fatalf("GetOutputDir failed: %v", err)
		}

		// Verify directory structure contains expected components
		if !strings.Contains(dir, "recally") {
			t.Errorf("directory path should contain 'recally': %s", dir)
		}
		if !strings.Contains(dir, "contents") {
			t.Errorf("directory path should contain 'contents': %s", dir)
		}
		if !strings.Contains(dir, "2026-01-18") {
			t.Errorf("directory path should contain date '2026-01-18': %s", dir)
		}

		// Verify platform-specific path structure
		switch runtime.GOOS {
		case "linux":
			if !strings.Contains(dir, ".local/share") {
				t.Errorf("Linux path should contain '.local/share': %s", dir)
			}
		case "darwin":
			if !strings.Contains(dir, "Library/Application Support") {
				t.Errorf("macOS path should contain 'Library/Application Support': %s", dir)
			}
		case "windows":
			// Windows paths use backslashes, but filepath handles this
			if !strings.Contains(dir, "recally") {
				t.Errorf("Windows path should contain 'recally': %s", dir)
			}
		}

		// Verify directory was created (and clean up afterward)
		defer func() { _ = os.RemoveAll(dir) }()
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("directory was not created: %s", dir)
		}
	})

	t.Run("with different dates", func(t *testing.T) {
		tmpDir := t.TempDir()

		date1 := time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC)
		date2 := time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)

		dir1, err := GetOutputDir(tmpDir, date1)
		if err != nil {
			t.Fatalf("GetOutputDir failed for date1: %v", err)
		}

		dir2, err := GetOutputDir(tmpDir, date2)
		if err != nil {
			t.Fatalf("GetOutputDir failed for date2: %v", err)
		}

		if !strings.Contains(dir1, "2026-01-18") {
			t.Errorf("dir1 should contain '2026-01-18': %s", dir1)
		}
		if !strings.Contains(dir2, "2025-12-25") {
			t.Errorf("dir2 should contain '2025-12-25': %s", dir2)
		}

		// Verify both directories were created
		if _, err := os.Stat(dir1); os.IsNotExist(err) {
			t.Errorf("dir1 was not created: %s", dir1)
		}
		if _, err := os.Stat(dir2); os.IsNotExist(err) {
			t.Errorf("dir2 was not created: %s", dir2)
		}
	})

	t.Run("directory creation with nested path", func(t *testing.T) {
		tmpDir := t.TempDir()
		customDir := filepath.Join(tmpDir, "deep", "nested", "custom")

		dir, err := GetOutputDir(customDir, testDate)
		if err != nil {
			t.Fatalf("GetOutputDir failed: %v", err)
		}

		// Verify all intermediate directories were created
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("nested directory was not created: %s", dir)
		}
	})
}

func TestResolveOutputPath(t *testing.T) {
	t.Run("no conflict", func(t *testing.T) {
		tmpDir := t.TempDir()

		path, err := ResolveOutputPath(tmpDir, "test-article")
		if err != nil {
			t.Fatalf("ResolveOutputPath failed: %v", err)
		}

		expected := filepath.Join(tmpDir, "test-article.md")
		if path != expected {
			t.Errorf("ResolveOutputPath() = %q, want %q", path, expected)
		}
	})

	t.Run("with conflict", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create a conflicting file
		existingFile := filepath.Join(tmpDir, "test-article.md")
		if err := os.WriteFile(existingFile, []byte("existing"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		path, err := ResolveOutputPath(tmpDir, "test-article")
		if err != nil {
			t.Fatalf("ResolveOutputPath failed: %v", err)
		}

		expected := filepath.Join(tmpDir, "test-article-1.md")
		if path != expected {
			t.Errorf("ResolveOutputPath() = %q, want %q", path, expected)
		}
	})

	t.Run("with multiple conflicts", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create multiple conflicting files
		for i := range 3 {
			var filename string
			if i == 0 {
				filename = "test-article.md"
			} else {
				filename = fmt.Sprintf("test-article-%d.md", i)
			}
			path := filepath.Join(tmpDir, filename)
			if err := os.WriteFile(path, []byte("existing"), 0644); err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}
		}

		path, err := ResolveOutputPath(tmpDir, "test-article")
		if err != nil {
			t.Fatalf("ResolveOutputPath failed: %v", err)
		}

		expected := filepath.Join(tmpDir, "test-article-3.md")
		if path != expected {
			t.Errorf("ResolveOutputPath() = %q, want %q", path, expected)
		}
	})

	t.Run("with title needing sanitization", func(t *testing.T) {
		tmpDir := t.TempDir()

		path, err := ResolveOutputPath(tmpDir, "Test Article!@#")
		if err != nil {
			t.Fatalf("ResolveOutputPath failed: %v", err)
		}

		expected := filepath.Join(tmpDir, "test-article.md")
		if path != expected {
			t.Errorf("ResolveOutputPath() = %q, want %q", path, expected)
		}
	})

	t.Run("with empty title", func(t *testing.T) {
		tmpDir := t.TempDir()

		path, err := ResolveOutputPath(tmpDir, "")
		if err != nil {
			t.Fatalf("ResolveOutputPath failed: %v", err)
		}

		// Should use untitled-{timestamp} format
		filename := filepath.Base(path)
		if !strings.HasPrefix(filename, "untitled-") {
			t.Errorf("expected filename to start with 'untitled-', got %q", filename)
		}
		if !strings.HasSuffix(filename, ".md") {
			t.Errorf("expected filename to end with '.md', got %q", filename)
		}
	})

	t.Run("with unicode title", func(t *testing.T) {
		tmpDir := t.TempDir()

		path, err := ResolveOutputPath(tmpDir, "文章标题")
		if err != nil {
			t.Fatalf("ResolveOutputPath failed: %v", err)
		}

		expected := filepath.Join(tmpDir, "文章标题.md")
		if path != expected {
			t.Errorf("ResolveOutputPath() = %q, want %q", path, expected)
		}
	})

	t.Run("with very long title", func(t *testing.T) {
		tmpDir := t.TempDir()

		longTitle := strings.Repeat("a", 250)
		path, err := ResolveOutputPath(tmpDir, longTitle)
		if err != nil {
			t.Fatalf("ResolveOutputPath failed: %v", err)
		}

		// Should be truncated with MD5 hash
		filename := filepath.Base(path)
		if len(filename) > 220 { // 200 + hyphen + 8 char hash + .md
			t.Errorf("filename too long: %d chars", len(filename))
		}
		if !strings.HasSuffix(filename, ".md") {
			t.Errorf("expected filename to end with '.md', got %q", filename)
		}
	})
}

func TestResolveOutputPathConflictResolution(t *testing.T) {
	tmpDir := t.TempDir()

	// Resolve paths multiple times and verify they're all unique
	paths := make([]string, 5)
	for i := range 5 {
		// Create the previously resolved path
		if i > 0 {
			if err := os.WriteFile(paths[i-1], []byte("test"), 0644); err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}
		}

		path, err := ResolveOutputPath(tmpDir, "test-article")
		if err != nil {
			t.Fatalf("ResolveOutputPath failed on iteration %d: %v", i, err)
		}
		paths[i] = path

		// Verify the path is unique
		for j := 0; j < i; j++ {
			if paths[i] == paths[j] {
				t.Errorf("duplicate path generated: %s (iteration %d and %d)", paths[i], i, j)
			}
		}
	}

	// Verify expected filenames
	expectedFilenames := []string{
		"test-article.md",
		"test-article-1.md",
		"test-article-2.md",
		"test-article-3.md",
		"test-article-4.md",
	}

	for i, path := range paths {
		filename := filepath.Base(path)
		if filename != expectedFilenames[i] {
			t.Errorf("iteration %d: expected filename %q, got %q", i, expectedFilenames[i], filename)
		}
	}
}

func TestGetOutputDirXDGDataHome(t *testing.T) {
	// Skip on non-Linux platforms as XDG_DATA_HOME is Linux-specific
	if runtime.GOOS != "linux" {
		t.Skip("XDG_DATA_HOME test is Linux-specific")
	}

	testDate := time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC)

	t.Run("honors XDG_DATA_HOME when set", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Save original value and restore after test
		originalXDG := os.Getenv("XDG_DATA_HOME")
		defer func() { _ = os.Setenv("XDG_DATA_HOME", originalXDG) }()

		// Set custom XDG_DATA_HOME
		_ = os.Setenv("XDG_DATA_HOME", tmpDir)

		dir, err := GetOutputDir("", testDate)
		if err != nil {
			t.Fatalf("GetOutputDir failed: %v", err)
		}

		expected := filepath.Join(tmpDir, "recally", "contents", "2026-01-18")
		if dir != expected {
			t.Errorf("GetOutputDir() = %q, want %q", dir, expected)
		}

		// Verify directory was created
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("directory was not created: %s", dir)
		}
	})

	t.Run("falls back to default when XDG_DATA_HOME is empty", func(t *testing.T) {
		// Save original value and restore after test
		originalXDG := os.Getenv("XDG_DATA_HOME")
		defer func() { _ = os.Setenv("XDG_DATA_HOME", originalXDG) }()

		// Unset XDG_DATA_HOME
		_ = os.Unsetenv("XDG_DATA_HOME")

		dir, err := GetOutputDir("", testDate)
		if err != nil {
			t.Fatalf("GetOutputDir failed: %v", err)
		}

		// Should use default ~/.local/share
		if !strings.Contains(dir, ".local/share") {
			t.Errorf("expected path to contain '.local/share' when XDG_DATA_HOME is unset: %s", dir)
		}

		// Clean up created directory
		defer func() { _ = os.RemoveAll(dir) }()
	})

	t.Run("custom dir takes precedence over XDG_DATA_HOME", func(t *testing.T) {
		tmpDir := t.TempDir()
		customDir := filepath.Join(tmpDir, "custom")
		xdgDir := filepath.Join(tmpDir, "xdg")

		// Save original value and restore after test
		originalXDG := os.Getenv("XDG_DATA_HOME")
		defer func() { _ = os.Setenv("XDG_DATA_HOME", originalXDG) }()

		// Set XDG_DATA_HOME
		_ = os.Setenv("XDG_DATA_HOME", xdgDir)

		dir, err := GetOutputDir(customDir, testDate)
		if err != nil {
			t.Fatalf("GetOutputDir failed: %v", err)
		}

		// Should use custom dir, not XDG_DATA_HOME
		expected := filepath.Join(customDir, "contents", "2026-01-18")
		if dir != expected {
			t.Errorf("GetOutputDir() = %q, want %q", dir, expected)
		}

		// XDG dir should not be created
		if _, err := os.Stat(xdgDir); !os.IsNotExist(err) {
			t.Errorf("XDG directory should not be created when custom dir is specified")
		}
	})
}

func TestResolveOutputPathStatErrors(t *testing.T) {
	t.Run("surfaces permission denied error on base path", func(t *testing.T) {
		// Skip on Windows as permission handling differs
		if runtime.GOOS == "windows" {
			t.Skip("Permission test not reliable on Windows")
		}

		tmpDir := t.TempDir()

		// Create a file with no read permissions
		noReadFile := filepath.Join(tmpDir, "test-article.md")
		if err := os.WriteFile(noReadFile, []byte("test"), 0000); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		defer func() { _ = os.Chmod(noReadFile, 0644) }() // Restore permissions for cleanup

		// Also remove read permission from directory to prevent stat
		// Note: On most systems, we can still stat a file even without read permission
		// So we need to create a directory without execute permission instead
		noExecDir := filepath.Join(tmpDir, "noexec")
		if err := os.Mkdir(noExecDir, 0755); err != nil {
			t.Fatalf("failed to create test dir: %v", err)
		}

		// Create a file inside
		innerFile := filepath.Join(noExecDir, "test-article.md")
		if err := os.WriteFile(innerFile, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create inner file: %v", err)
		}

		// Remove execute permission from directory (prevents stat on contents)
		if err := os.Chmod(noExecDir, 0600); err != nil {
			t.Fatalf("failed to chmod: %v", err)
		}
		defer func() { _ = os.Chmod(noExecDir, 0755) }() // Restore for cleanup

		_, err := ResolveOutputPath(noExecDir, "test-article")
		if err == nil {
			t.Error("expected error when stat fails with permission denied")
		}
		if !strings.Contains(err.Error(), "failed to check file") {
			t.Errorf("expected 'failed to check file' error, got: %v", err)
		}
	})

	t.Run("returns error for non-existent parent directory", func(t *testing.T) {
		// Try to resolve path in a directory that doesn't exist
		nonExistentDir := "/non/existent/path/that/should/not/exist"

		_, err := ResolveOutputPath(nonExistentDir, "test-article")
		// This should succeed since we're just building a path
		// The file doesn't exist, so os.IsNotExist returns true
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
