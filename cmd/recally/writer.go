package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"recally/internal/pkg/webreader"

	"gopkg.in/yaml.v3"
)

const (
	// MinDiskSpace is the minimum required disk space before writing (100MB)
	MinDiskSpace = 100 * 1024 * 1024
)

// frontmatter represents the YAML frontmatter structure
type frontmatter struct {
	URL           string `yaml:"url"`
	Title         string `yaml:"title"`
	Author        string `yaml:"author,omitempty"`
	Description   string `yaml:"description,omitempty"`
	SiteName      string `yaml:"site_name,omitempty"`
	PublishedTime string `yaml:"published_time,omitempty"`
	ModifiedTime  string `yaml:"modified_time,omitempty"`
	Cover         string `yaml:"cover,omitempty"`
	Favicon       string `yaml:"favicon,omitempty"`
	SavedAt       string `yaml:"saved_at"`
}

// WriteMarkdown writes Content to markdown file with YAML frontmatter.
// Returns the full path to the written file.
//
// File format:
//
//	---
//	url: https://example.com/article
//	title: Article Title
//	author: Author Name
//	...
//	---
//
//	Markdown content here...
//
// Disk space check:
//   - Requires at least 100MB free disk space
//   - Returns error if insufficient space
//
// Conflict resolution:
//   - Uses ResolveOutputPath to handle filename conflicts
//   - Appends -N counter if file exists (article-1.md, article-2.md, etc.)
func WriteMarkdown(content *webreader.Content, outputDir string) (string, error) {
	// Check disk space before writing
	if err := checkDiskSpace(outputDir); err != nil {
		return "", fmt.Errorf("disk space check failed: %w", err)
	}

	// Resolve output path (handles conflicts with -N counter)
	outputPath, err := ResolveOutputPath(outputDir, content.Title)
	if err != nil {
		return "", fmt.Errorf("failed to resolve output path: %w", err)
	}

	// Generate markdown content with frontmatter
	markdown, err := formatMarkdown(content)
	if err != nil {
		return "", fmt.Errorf("failed to format markdown: %w", err)
	}

	// Write to file with 0644 permissions (rw-r--r--)
	if err := os.WriteFile(outputPath, []byte(markdown), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return outputPath, nil
}

// formatMarkdown generates the complete markdown file content
// with YAML frontmatter header and markdown body
func formatMarkdown(content *webreader.Content) (string, error) {
	// Create frontmatter struct
	fm := frontmatter{
		URL:         content.URL,
		Title:       content.Title,
		Author:      content.Author,
		Description: content.Description,
		SiteName:    content.SiteName,
		Cover:       content.Cover,
		Favicon:     content.Favicon,
		SavedAt:     time.Now().UTC().Format(time.RFC3339),
	}

	// Handle nil times: render as empty string, not "null"
	if content.PublishedTime != nil {
		fm.PublishedTime = content.PublishedTime.Format(time.RFC3339)
	}
	if content.ModifiedTime != nil {
		fm.ModifiedTime = content.ModifiedTime.Format(time.RFC3339)
	}

	// Marshal frontmatter to YAML
	yamlBytes, err := yaml.Marshal(&fm)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}

	// Build complete markdown file
	var builder strings.Builder
	builder.WriteString("---\n")
	builder.Write(yamlBytes)
	builder.WriteString("---\n\n")
	builder.WriteString(content.Markwdown)

	return builder.String(), nil
}

// checkDiskSpace verifies that at least MinDiskSpace bytes are available
// on the filesystem containing the output directory
func checkDiskSpace(dir string) error {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(dir, &stat); err != nil {
		return fmt.Errorf("failed to get filesystem stats: %w", err)
	}

	// Calculate available space
	// Bavail is available blocks for unprivileged users
	// Bsize is the filesystem block size
	availableSpace := uint64(stat.Bavail) * uint64(stat.Bsize)

	if availableSpace < MinDiskSpace {
		return fmt.Errorf("insufficient disk space: %d bytes available, %d bytes required",
			availableSpace, MinDiskSpace)
	}

	return nil
}
