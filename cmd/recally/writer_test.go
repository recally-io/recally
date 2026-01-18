package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"recally/internal/pkg/webreader"

	"gopkg.in/yaml.v3"
)

func TestWriteMarkdown(t *testing.T) {
	// Create temporary directory for test output
	tempDir := t.TempDir()

	// Create test content with all fields populated
	publishedTime := time.Date(2026, 1, 18, 10, 0, 0, 0, time.UTC)
	modifiedTime := time.Date(2026, 1, 18, 12, 0, 0, 0, time.UTC)

	content := &webreader.Content{
		URL:           "https://example.com/article",
		Title:         "Test Article Title",
		Author:        "John Doe",
		Description:   "This is a test article description",
		SiteName:      "Example Site",
		PublishedTime: &publishedTime,
		ModifiedTime:  &modifiedTime,
		Cover:         "https://example.com/cover.jpg",
		Favicon:       "https://example.com/favicon.ico",
		Markwdown:     "# Article Content\n\nThis is the article body.",
	}

	// Write markdown file
	outputPath, err := WriteMarkdown(content, tempDir)
	if err != nil {
		t.Fatalf("WriteMarkdown failed: %v", err)
	}

	// Verify output path
	expectedFilename := "test-article-title.md"
	if !strings.HasSuffix(outputPath, expectedFilename) {
		t.Errorf("Expected filename %s, got %s", expectedFilename, filepath.Base(outputPath))
	}

	// Verify file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file does not exist: %s", outputPath)
	}

	// Read file content
	fileContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	fileStr := string(fileContent)

	// Verify frontmatter structure
	if !strings.HasPrefix(fileStr, "---\n") {
		t.Error("File should start with '---\\n'")
	}

	// Split frontmatter and content
	parts := strings.SplitN(fileStr, "---\n", 3)
	if len(parts) != 3 {
		t.Fatalf("Expected 3 parts (empty, frontmatter, content), got %d", len(parts))
	}

	frontmatterYAML := parts[1]
	markdownContent := parts[2]

	// Parse and verify frontmatter
	var fm frontmatter
	if err := yaml.Unmarshal([]byte(frontmatterYAML), &fm); err != nil {
		t.Fatalf("Failed to parse YAML frontmatter: %v", err)
	}

	// Verify all fields
	if fm.URL != content.URL {
		t.Errorf("URL mismatch: expected %s, got %s", content.URL, fm.URL)
	}
	if fm.Title != content.Title {
		t.Errorf("Title mismatch: expected %s, got %s", content.Title, fm.Title)
	}
	if fm.Author != content.Author {
		t.Errorf("Author mismatch: expected %s, got %s", content.Author, fm.Author)
	}
	if fm.Description != content.Description {
		t.Errorf("Description mismatch: expected %s, got %s", content.Description, fm.Description)
	}
	if fm.SiteName != content.SiteName {
		t.Errorf("SiteName mismatch: expected %s, got %s", content.SiteName, fm.SiteName)
	}
	if fm.Cover != content.Cover {
		t.Errorf("Cover mismatch: expected %s, got %s", content.Cover, fm.Cover)
	}
	if fm.Favicon != content.Favicon {
		t.Errorf("Favicon mismatch: expected %s, got %s", content.Favicon, fm.Favicon)
	}

	// Verify time formatting (RFC3339)
	expectedPublished := publishedTime.Format(time.RFC3339)
	if fm.PublishedTime != expectedPublished {
		t.Errorf("PublishedTime mismatch: expected %s, got %s", expectedPublished, fm.PublishedTime)
	}
	expectedModified := modifiedTime.Format(time.RFC3339)
	if fm.ModifiedTime != expectedModified {
		t.Errorf("ModifiedTime mismatch: expected %s, got %s", expectedModified, fm.ModifiedTime)
	}

	// Verify saved_at is in UTC and RFC3339 format
	if fm.SavedAt == "" {
		t.Error("SavedAt should not be empty")
	}
	savedAt, err := time.Parse(time.RFC3339, fm.SavedAt)
	if err != nil {
		t.Errorf("SavedAt should be valid RFC3339: %v", err)
	}
	if savedAt.Location() != time.UTC {
		t.Errorf("SavedAt should be UTC, got %s", savedAt.Location())
	}

	// Verify markdown content
	if !strings.Contains(markdownContent, content.Markwdown) {
		t.Error("Markdown content should be present in output")
	}

	// Verify file permissions
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Failed to stat output file: %v", err)
	}
	if info.Mode().Perm() != 0644 {
		t.Errorf("Expected file permissions 0644, got %o", info.Mode().Perm())
	}
}

func TestWriteMarkdownWithNilTimes(t *testing.T) {
	tempDir := t.TempDir()

	// Create content with nil times
	content := &webreader.Content{
		URL:           "https://example.com/article",
		Title:         "Article Without Times",
		PublishedTime: nil, // Nil time
		ModifiedTime:  nil, // Nil time
		Markwdown:     "# Content",
	}

	outputPath, err := WriteMarkdown(content, tempDir)
	if err != nil {
		t.Fatalf("WriteMarkdown failed: %v", err)
	}

	// Read file content
	fileContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	fileStr := string(fileContent)

	// Split frontmatter
	parts := strings.SplitN(fileStr, "---\n", 3)
	if len(parts) != 3 {
		t.Fatalf("Expected 3 parts, got %d", len(parts))
	}

	frontmatterYAML := parts[1]

	// Parse frontmatter
	var fm frontmatter
	if err := yaml.Unmarshal([]byte(frontmatterYAML), &fm); err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}

	// Verify nil times render as empty string, not "null"
	if fm.PublishedTime != "" {
		t.Errorf("PublishedTime should be empty string for nil time, got: %s", fm.PublishedTime)
	}
	if fm.ModifiedTime != "" {
		t.Errorf("ModifiedTime should be empty string for nil time, got: %s", fm.ModifiedTime)
	}

	// Verify YAML doesn't contain "null" strings
	if strings.Contains(frontmatterYAML, "null") {
		t.Error("YAML frontmatter should not contain 'null' for nil times")
	}

	// Verify optional fields are omitted when empty (omitempty tag)
	if strings.Contains(frontmatterYAML, "published_time:") && fm.PublishedTime == "" {
		t.Error("Empty published_time should be omitted from YAML (omitempty)")
	}
	if strings.Contains(frontmatterYAML, "modified_time:") && fm.ModifiedTime == "" {
		t.Error("Empty modified_time should be omitted from YAML (omitempty)")
	}
}

func TestWriteMarkdownConflict(t *testing.T) {
	tempDir := t.TempDir()

	content := &webreader.Content{
		URL:       "https://example.com/article",
		Title:     "Duplicate Article",
		Markwdown: "# Content",
	}

	// Write first file
	path1, err := WriteMarkdown(content, tempDir)
	if err != nil {
		t.Fatalf("First WriteMarkdown failed: %v", err)
	}

	// Verify first file
	expectedFilename1 := "duplicate-article.md"
	if filepath.Base(path1) != expectedFilename1 {
		t.Errorf("Expected first filename %s, got %s", expectedFilename1, filepath.Base(path1))
	}

	// Write second file with same title
	path2, err := WriteMarkdown(content, tempDir)
	if err != nil {
		t.Fatalf("Second WriteMarkdown failed: %v", err)
	}

	// Verify conflict resolution with -1 suffix
	expectedFilename2 := "duplicate-article-1.md"
	if filepath.Base(path2) != expectedFilename2 {
		t.Errorf("Expected second filename %s, got %s", expectedFilename2, filepath.Base(path2))
	}

	// Write third file
	path3, err := WriteMarkdown(content, tempDir)
	if err != nil {
		t.Fatalf("Third WriteMarkdown failed: %v", err)
	}

	// Verify -2 suffix
	expectedFilename3 := "duplicate-article-2.md"
	if filepath.Base(path3) != expectedFilename3 {
		t.Errorf("Expected third filename %s, got %s", expectedFilename3, filepath.Base(path3))
	}

	// Verify all three files exist
	for _, path := range []string{path1, path2, path3} {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file to exist: %s", path)
		}
	}

	// Verify files are different (have different paths)
	if path1 == path2 || path2 == path3 || path1 == path3 {
		t.Error("All three paths should be unique")
	}
}

func TestWriteMarkdownEmptyContent(t *testing.T) {
	tempDir := t.TempDir()

	// Content with empty markdown
	content := &webreader.Content{
		URL:       "https://example.com/empty",
		Title:     "Empty Article",
		Markwdown: "", // Empty markdown
	}

	outputPath, err := WriteMarkdown(content, tempDir)
	if err != nil {
		t.Fatalf("WriteMarkdown should handle empty content: %v", err)
	}

	// Verify file exists
	fileContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Should still have frontmatter
	if !strings.HasPrefix(string(fileContent), "---\n") {
		t.Error("File should have frontmatter even with empty content")
	}

	// Verify frontmatter is parseable
	parts := strings.SplitN(string(fileContent), "---\n", 3)
	if len(parts) != 3 {
		t.Fatalf("Expected 3 parts, got %d", len(parts))
	}

	var fm frontmatter
	if err := yaml.Unmarshal([]byte(parts[1]), &fm); err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}

	if fm.URL != content.URL {
		t.Errorf("URL should be preserved: expected %s, got %s", content.URL, fm.URL)
	}
}

func TestWriteMarkdownEmptyTitle(t *testing.T) {
	tempDir := t.TempDir()

	// Content with empty title (should use untitled-{timestamp})
	content := &webreader.Content{
		URL:       "https://example.com/notitle",
		Title:     "", // Empty title
		Markwdown: "# Content without title",
	}

	outputPath, err := WriteMarkdown(content, tempDir)
	if err != nil {
		t.Fatalf("WriteMarkdown failed: %v", err)
	}

	// Verify filename uses untitled-{timestamp} pattern
	filename := filepath.Base(outputPath)
	if !strings.HasPrefix(filename, "untitled-") {
		t.Errorf("Empty title should use 'untitled-{timestamp}' pattern, got: %s", filename)
	}
	if !strings.HasSuffix(filename, ".md") {
		t.Errorf("Filename should end with .md, got: %s", filename)
	}
}

func TestWriteMarkdownInvalidDirectory(t *testing.T) {
	// Use a non-existent directory that cannot be created
	invalidDir := "/this/path/definitely/does/not/exist/and/cannot/be/created"

	content := &webreader.Content{
		URL:       "https://example.com/article",
		Title:     "Test Article",
		Markwdown: "# Content",
	}

	// Should fail because directory doesn't exist
	// Could fail at disk space check, path resolution, or file writing
	_, err := WriteMarkdown(content, invalidDir)
	if err == nil {
		t.Error("WriteMarkdown should fail with invalid directory")
	}

	// Error should mention disk space check, path resolution, or file writing
	errStr := err.Error()
	if !strings.Contains(errStr, "disk space check") &&
		!strings.Contains(errStr, "resolve output path") &&
		!strings.Contains(errStr, "write file") {
		t.Errorf("Error should mention disk/path/file issue, got: %s", errStr)
	}
}

func TestFormatMarkdown(t *testing.T) {
	publishedTime := time.Date(2026, 1, 18, 10, 0, 0, 0, time.UTC)
	modifiedTime := time.Date(2026, 1, 18, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		content *webreader.Content
		checks  []func(t *testing.T, markdown string)
	}{
		{
			name: "all fields populated",
			content: &webreader.Content{
				URL:           "https://example.com/article",
				Title:         "Test Article",
				Author:        "John Doe",
				Description:   "Test description",
				SiteName:      "Example",
				PublishedTime: &publishedTime,
				ModifiedTime:  &modifiedTime,
				Cover:         "https://example.com/cover.jpg",
				Favicon:       "https://example.com/favicon.ico",
				Markwdown:     "# Article\n\nContent here.",
			},
			checks: []func(t *testing.T, markdown string){
				func(t *testing.T, markdown string) {
					if !strings.HasPrefix(markdown, "---\n") {
						t.Error("Should start with ---")
					}
				},
				func(t *testing.T, markdown string) {
					if !strings.Contains(markdown, "url: https://example.com/article") {
						t.Error("Should contain URL")
					}
				},
				func(t *testing.T, markdown string) {
					if !strings.Contains(markdown, "title: Test Article") {
						t.Error("Should contain title")
					}
				},
				func(t *testing.T, markdown string) {
					if !strings.Contains(markdown, "author: John Doe") {
						t.Error("Should contain author")
					}
				},
				func(t *testing.T, markdown string) {
					if !strings.Contains(markdown, "# Article") {
						t.Error("Should contain markdown content")
					}
				},
				func(t *testing.T, markdown string) {
					expectedPublished := publishedTime.Format(time.RFC3339)
					// YAML may format the timestamp value without the explicit prefix
					// Check if the expected timestamp is present anywhere in the markdown
					if !strings.Contains(markdown, expectedPublished) {
						t.Errorf("Should contain published_time in RFC3339 format (%s), got markdown:\n%s",
							expectedPublished, markdown)
					}
				},
			},
		},
		{
			name: "minimal fields",
			content: &webreader.Content{
				URL:       "https://example.com/minimal",
				Title:     "Minimal",
				Markwdown: "Content",
			},
			checks: []func(t *testing.T, markdown string){
				func(t *testing.T, markdown string) {
					if !strings.HasPrefix(markdown, "---\n") {
						t.Error("Should have frontmatter")
					}
				},
				func(t *testing.T, markdown string) {
					if !strings.Contains(markdown, "url: https://example.com/minimal") {
						t.Error("Should contain URL")
					}
				},
				func(t *testing.T, markdown string) {
					if !strings.Contains(markdown, "saved_at:") {
						t.Error("Should contain saved_at timestamp")
					}
				},
			},
		},
		{
			name: "nil times",
			content: &webreader.Content{
				URL:           "https://example.com/notimes",
				Title:         "No Times",
				PublishedTime: nil,
				ModifiedTime:  nil,
				Markwdown:     "Content",
			},
			checks: []func(t *testing.T, markdown string){
				func(t *testing.T, markdown string) {
					// Should not contain "null" for nil times
					if strings.Contains(markdown, "null") {
						t.Error("Should not contain 'null' string")
					}
				},
				func(t *testing.T, markdown string) {
					// Parse the YAML to verify empty strings
					parts := strings.SplitN(markdown, "---\n", 3)
					if len(parts) == 3 {
						var fm frontmatter
						if err := yaml.Unmarshal([]byte(parts[1]), &fm); err == nil {
							if fm.PublishedTime != "" {
								t.Error("PublishedTime should be empty string")
							}
							if fm.ModifiedTime != "" {
								t.Error("ModifiedTime should be empty string")
							}
						}
					}
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			markdown, err := formatMarkdown(tt.content)
			if err != nil {
				t.Fatalf("formatMarkdown failed: %v", err)
			}

			for _, check := range tt.checks {
				check(t, markdown)
			}
		})
	}
}

func TestFormatMarkdownYAMLValidity(t *testing.T) {
	content := &webreader.Content{
		URL:       "https://example.com/test",
		Title:     "Test",
		Markwdown: "Content",
	}

	markdown, err := formatMarkdown(content)
	if err != nil {
		t.Fatalf("formatMarkdown failed: %v", err)
	}

	// Extract YAML frontmatter
	parts := strings.SplitN(markdown, "---\n", 3)
	if len(parts) != 3 {
		t.Fatalf("Expected 3 parts, got %d", len(parts))
	}

	// Verify YAML is valid
	var fm frontmatter
	if err := yaml.Unmarshal([]byte(parts[1]), &fm); err != nil {
		t.Fatalf("YAML frontmatter is invalid: %v", err)
	}

	// Verify parsed values
	if fm.URL != content.URL {
		t.Errorf("URL mismatch: expected %s, got %s", content.URL, fm.URL)
	}
	if fm.Title != content.Title {
		t.Errorf("Title mismatch: expected %s, got %s", content.Title, fm.Title)
	}
	if fm.SavedAt == "" {
		t.Error("SavedAt should not be empty")
	}
}

func TestCheckDiskSpace(t *testing.T) {
	// Test with current directory (should have space)
	tempDir := t.TempDir()
	err := checkDiskSpace(tempDir)
	if err != nil {
		t.Errorf("checkDiskSpace should succeed for temp directory: %v", err)
	}

	// Test with invalid directory
	err = checkDiskSpace("/this/path/does/not/exist")
	if err == nil {
		t.Error("checkDiskSpace should fail for non-existent directory")
	}
	if !strings.Contains(err.Error(), "failed to get filesystem stats") {
		t.Errorf("Expected filesystem stats error, got: %v", err)
	}
}
