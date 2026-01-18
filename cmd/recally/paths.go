package main

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode"
)

var (
	// Regex for replacing multiple whitespace with single hyphen
	whitespaceRegex = regexp.MustCompile(`\s+`)
)

// GetOutputDir returns the XDG-compliant output directory for the given date.
// If customDir is not empty, uses that instead of XDG default.
// Creates the directory if it doesn't exist.
//
// Directory structure:
//   - Linux/macOS: ~/.local/share/recally/contents/2026-01-18/
//   - macOS (alt): ~/Library/Application Support/recally/contents/2026-01-18/
//   - Windows:     %LOCALAPPDATA%\recally\contents\2026-01-18\
func GetOutputDir(customDir string, date time.Time) (string, error) {
	var baseDir string

	if customDir != "" {
		baseDir = customDir
	} else {
		// Get user config directory (cross-platform)
		configDir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user config directory: %w", err)
		}

		// Convert config directory to data directory based on OS
		switch runtime.GOOS {
		case "linux":
			// On Linux, UserConfigDir returns ~/.config
			// We want ~/.local/share instead
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get home directory: %w", err)
			}
			baseDir = filepath.Join(homeDir, ".local", "share", "recally")
		case "darwin":
			// On macOS, UserConfigDir returns ~/Library/Application Support
			// This is actually what we want for data as well
			baseDir = filepath.Join(configDir, "recally")
		case "windows":
			// On Windows, UserConfigDir returns %LOCALAPPDATA%
			// This is correct for our purposes
			baseDir = filepath.Join(configDir, "recally")
		default:
			// Fallback for other platforms
			baseDir = filepath.Join(configDir, "recally")
		}
	}

	// Format date as YYYY-MM-DD
	dateStr := date.Format("2006-01-02")

	// Build full directory path
	fullPath := filepath.Join(baseDir, "contents", dateStr)

	// Create directory with 0755 permissions
	// MkdirAll is safe against symlink attacks (creates all intermediate dirs)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", fullPath, err)
	}

	return fullPath, nil
}

// SanitizeTitle converts article title to a safe filename.
// Sanitization rules (applied in order):
//  1. Trim leading/trailing whitespace
//  2. Convert to lowercase
//  3. Replace multiple spaces/whitespace with single hyphen
//  4. Remove non-alphanumeric characters except hyphens (preserve Unicode letters/numbers)
//  5. Truncate to 200 chars, add MD5 hash suffix if truncated
//  6. If empty after sanitization: use "untitled-{unix-timestamp}"
func SanitizeTitle(title string) string {
	// 1. Trim whitespace
	title = strings.TrimSpace(title)

	// 2. Convert to lowercase
	title = strings.ToLower(title)

	// 3. Replace multiple whitespace with single hyphen
	title = whitespaceRegex.ReplaceAllString(title, "-")

	// 4. Remove non-alphanumeric except hyphens (preserve Unicode letters/numbers)
	var builder strings.Builder
	for _, r := range title {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' {
			builder.WriteRune(r)
		}
	}
	title = builder.String()

	// Remove leading/trailing hyphens and collapse multiple consecutive hyphens
	title = strings.Trim(title, "-")
	title = strings.Join(strings.FieldsFunc(title, func(r rune) bool {
		return r == '-'
	}), "-")

	// 6. If empty, use timestamp fallback
	if title == "" {
		return fmt.Sprintf("untitled-%d", time.Now().Unix())
	}

	// 5. Truncate to 200 chars with MD5 hash if needed
	if len(title) > 200 {
		// Calculate MD5 hash of the full title
		hash := md5.Sum([]byte(title))
		hashStr := fmt.Sprintf("%x", hash)[:8] // Use first 8 chars of hash

		// Truncate title to 200 chars and append hash
		title = title[:200] + "-" + hashStr
	}

	return title
}

// ResolveOutputPath generates a unique output file path.
// Handles filename conflicts by appending a counter (-N).
// Returns the full path to a non-existent file.
//
// Examples:
//   - First save:  /path/to/article.md
//   - Conflict 1:  /path/to/article-1.md
//   - Conflict 2:  /path/to/article-2.md
func ResolveOutputPath(dir, title string) (string, error) {
	// Sanitize the title to create base filename
	baseFilename := SanitizeTitle(title)

	// Try the base filename first
	basePath := filepath.Join(dir, baseFilename+".md")

	// Check if file exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		// File doesn't exist, we can use this path
		return basePath, nil
	}

	// File exists, start appending counters
	counter := 1
	for {
		filename := fmt.Sprintf("%s-%d.md", baseFilename, counter)
		fullPath := filepath.Join(dir, filename)

		// Check if this path exists
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			// Found a unique filename
			return fullPath, nil
		}

		counter++

		// Safety check: prevent infinite loop (extremely unlikely)
		if counter > 10000 {
			return "", fmt.Errorf("failed to find unique filename after 10000 attempts")
		}
	}
}
