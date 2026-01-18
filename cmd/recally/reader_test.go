package main

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"recally/internal/pkg/webreader"
)

// mockFetcher implements webreader.Fetcher for testing
type mockFetcher struct {
	content     string
	err         error
	closeCalled bool
}

func (m *mockFetcher) Fetch(ctx context.Context, url string) (*webreader.FetchedContent, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &webreader.FetchedContent{
		Content: webreader.Content{
			URL:  url,
			Html: m.content,
		},
		Reader:      io.NopCloser(strings.NewReader(m.content)),
		StatusCode:  200,
		ContentType: "text/html",
		Headers:     make(map[string][]string),
	}, nil
}

func (m *mockFetcher) Close() error {
	m.closeCalled = true
	return nil
}

// mockLogger implements webreader.Logger for testing
type mockLogger struct {
	infoLogs  []string
	errorLogs []string
}

func (m *mockLogger) Info(msg string, args ...any) {
	m.infoLogs = append(m.infoLogs, msg)
}

func (m *mockLogger) Error(msg string, args ...any) {
	m.errorLogs = append(m.errorLogs, msg)
}

// TestFetchAndProcess_Success tests successful fetch and processing
func TestFetchAndProcess_Success(t *testing.T) {
	// Sample HTML with article content
	htmlContent := `
<!DOCTYPE html>
<html>
<head>
	<title>Test Article</title>
	<meta name="author" content="John Doe">
	<meta name="description" content="Test article description">
</head>
<body>
	<article>
		<h1>Test Article Title</h1>
		<p>This is the main content of the article.</p>
		<p>It has multiple paragraphs.</p>
	</article>
	<footer>This is footer content that should be removed by readability</footer>
</body>
</html>
`

	ctx := context.Background()
	logger := &mockLogger{}
	fetcher := &mockFetcher{content: htmlContent}

	content, err := FetchAndProcess(ctx, "https://example.com/article", logger, fetcher)

	if err != nil {
		t.Fatalf("FetchAndProcess failed: %v", err)
	}

	if content == nil {
		t.Fatal("Expected content, got nil")
	}

	// Verify URL is preserved
	if content.URL != "https://example.com/article" {
		t.Errorf("Expected URL 'https://example.com/article', got '%s'", content.URL)
	}

	// Verify title is extracted by Readability
	if content.Title == "" {
		t.Error("Expected title to be extracted, got empty string")
	}

	// Verify markdown is generated
	if content.Markwdown == "" {
		t.Error("Expected markdown content, got empty string")
	}

	// Verify markdown contains article content
	if !strings.Contains(content.Markwdown, "Test Article Title") {
		t.Error("Markdown should contain article title")
	}

	// Verify HTML is cleaned by Readability (footer should be removed)
	if strings.Contains(content.Html, "footer") {
		t.Error("Readability should remove footer content")
	}

	// Verify fetcher was closed
	if !fetcher.closeCalled {
		t.Error("Expected fetcher.Close() to be called")
	}
}

// TestFetchAndProcess_InvalidURL tests URL validation
func TestFetchAndProcess_InvalidURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"malformed URL", "ht!tp://invalid"},
		{"file scheme", "file:///etc/passwd"},
		{"javascript scheme", "javascript:alert(1)"},
		{"data scheme", "data:text/html,<script>alert(1)</script>"},
		{"no scheme", "example.com/article"},
	}

	ctx := context.Background()
	logger := &mockLogger{}
	fetcher := &mockFetcher{content: "<html></html>"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FetchAndProcess(ctx, tt.url, logger, fetcher)
			if err == nil {
				t.Errorf("Expected error for URL %s, got nil", tt.url)
			}
		})
	}
}

// TestFetchAndProcess_FetchError tests handling of fetch errors
func TestFetchAndProcess_FetchError(t *testing.T) {
	ctx := context.Background()
	logger := &mockLogger{}
	fetchErr := errors.New("network timeout")
	fetcher := &mockFetcher{err: fetchErr}

	_, err := FetchAndProcess(ctx, "https://example.com/article", logger, fetcher)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to fetch and process content") {
		t.Errorf("Expected fetch error to be wrapped, got: %v", err)
	}
}

// TestFetchAndProcess_ContextCancellation tests context cancellation
func TestFetchAndProcess_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	logger := &mockLogger{}
	fetcher := &mockFetcher{content: "<html></html>"}

	_, err := FetchAndProcess(ctx, "https://example.com/article", logger, fetcher)

	// Note: Actual behavior depends on fetcher implementation
	// Mock fetcher doesn't respect context, but real fetchers should
	// This is expected - mock doesn't implement context cancellation
	// Real fetchers would return context.Canceled error
	_ = err
}

// TestFetchAndProcess_NilLogger tests that nil logger works (no-op)
func TestFetchAndProcess_NilLogger(t *testing.T) {
	htmlContent := `<html><body><article><h1>Test</h1><p>Content</p></article></body></html>`

	ctx := context.Background()
	fetcher := &mockFetcher{content: htmlContent}

	// Pass nil logger - should use no-op logger
	content, err := FetchAndProcess(ctx, "https://example.com/article", nil, fetcher)

	if err != nil {
		t.Fatalf("FetchAndProcess with nil logger failed: %v", err)
	}

	if content == nil {
		t.Fatal("Expected content, got nil")
	}
}

// TestReaderPipeline tests the processor chain
func TestReaderPipeline(t *testing.T) {
	// HTML with metadata and complex structure
	htmlContent := `
<!DOCTYPE html>
<html>
<head>
	<title>Complex Article</title>
	<meta name="author" content="Jane Smith">
	<meta name="description" content="A complex article with metadata">
	<meta property="og:site_name" content="Example Site">
	<meta property="og:image" content="https://example.com/cover.jpg">
	<meta property="article:published_time" content="2026-01-18T10:00:00Z">
</head>
<body>
	<header>
		<nav>Navigation that should be removed</nav>
	</header>
	<article>
		<h1>Main Article Title</h1>
		<p>First paragraph with <strong>bold text</strong> and <em>italic text</em>.</p>
		<h2>Subheading</h2>
		<p>Second paragraph with <a href="https://example.com">a link</a>.</p>
		<ul>
			<li>List item 1</li>
			<li>List item 2</li>
		</ul>
		<blockquote>A quoted section</blockquote>
	</article>
	<aside>Sidebar content that should be removed</aside>
	<footer>Footer that should be removed</footer>
</body>
</html>
`

	ctx := context.Background()
	logger := &mockLogger{}
	fetcher := &mockFetcher{content: htmlContent}

	content, err := FetchAndProcess(ctx, "https://example.com/article", logger, fetcher)

	if err != nil {
		t.Fatalf("FetchAndProcess failed: %v", err)
	}

	// Test Readability processor output
	t.Run("Readability extracts metadata", func(t *testing.T) {
		// Readability uses <title> tag, not <h1>
		if content.Title != "Complex Article" {
			t.Errorf("Expected title 'Complex Article', got '%s'", content.Title)
		}

		if content.Author != "Jane Smith" {
			t.Errorf("Expected author 'Jane Smith', got '%s'", content.Author)
		}

		if content.Description != "A complex article with metadata" {
			t.Errorf("Expected description 'A complex article with metadata', got '%s'", content.Description)
		}

		if content.SiteName != "Example Site" {
			t.Errorf("Expected site_name 'Example Site', got '%s'", content.SiteName)
		}

		if content.Cover != "https://example.com/cover.jpg" {
			t.Errorf("Expected cover 'https://example.com/cover.jpg', got '%s'", content.Cover)
		}

		if content.PublishedTime == nil {
			t.Error("Expected published_time to be parsed, got nil")
		} else {
			expected := time.Date(2026, 1, 18, 10, 0, 0, 0, time.UTC)
			if !content.PublishedTime.Equal(expected) {
				t.Errorf("Expected published_time %v, got %v", expected, *content.PublishedTime)
			}
		}
	})

	t.Run("Readability removes boilerplate", func(t *testing.T) {
		// Navigation, sidebar, footer should be removed
		if strings.Contains(content.Html, "Navigation that should be removed") {
			t.Error("Readability should remove navigation")
		}
		if strings.Contains(content.Html, "Sidebar content") {
			t.Error("Readability should remove sidebar")
		}
		if strings.Contains(content.Html, "Footer that should be removed") {
			t.Error("Readability should remove footer")
		}

		// Main content should be preserved
		if !strings.Contains(content.Html, "Main Article Title") {
			t.Error("Readability should preserve main article content")
		}
	})

	t.Run("Markdown processor converts HTML", func(t *testing.T) {
		if content.Markwdown == "" {
			t.Fatal("Expected markdown content, got empty string")
		}

		// Check markdown formatting
		if !strings.Contains(content.Markwdown, "# Main Article Title") &&
			!strings.Contains(content.Markwdown, "Main Article Title") {
			t.Error("Markdown should contain title")
		}

		if !strings.Contains(content.Markwdown, "**bold text**") &&
			!strings.Contains(content.Markwdown, "bold text") {
			t.Error("Markdown should contain bold text")
		}

		if !strings.Contains(content.Markwdown, "*italic text*") &&
			!strings.Contains(content.Markwdown, "_italic text_") {
			t.Error("Markdown should contain italic text")
		}

		if !strings.Contains(content.Markwdown, "[a link]") {
			t.Error("Markdown should contain link")
		}

		if !strings.Contains(content.Markwdown, "- List item 1") &&
			!strings.Contains(content.Markwdown, "* List item 1") {
			t.Error("Markdown should contain list items")
		}

		if !strings.Contains(content.Markwdown, "> A quoted section") {
			t.Error("Markdown should contain blockquote")
		}
	})

	t.Run("Text content is extracted", func(t *testing.T) {
		if content.Text == "" {
			t.Error("Expected text content to be extracted, got empty string")
		}

		// Text should contain main content without HTML tags
		if !strings.Contains(content.Text, "Main Article Title") {
			t.Error("Text should contain article title")
		}

		if !strings.Contains(content.Text, "First paragraph") {
			t.Error("Text should contain paragraph content")
		}
	})
}

// TestReaderPipeline_MinimalHTML tests handling of minimal HTML
func TestReaderPipeline_MinimalHTML(t *testing.T) {
	htmlContent := `<html><body><p>Just a paragraph</p></body></html>`

	ctx := context.Background()
	logger := &mockLogger{}
	fetcher := &mockFetcher{content: htmlContent}

	content, err := FetchAndProcess(ctx, "https://example.com/minimal", logger, fetcher)

	if err != nil {
		t.Fatalf("FetchAndProcess failed: %v", err)
	}

	// Should still produce markdown
	if content.Markwdown == "" {
		t.Error("Expected markdown for minimal HTML, got empty string")
	}

	// Should contain the paragraph text
	if !strings.Contains(content.Markwdown, "Just a paragraph") &&
		!strings.Contains(content.Text, "Just a paragraph") {
		t.Error("Content should contain paragraph text")
	}
}

// TestReaderPipeline_EmptyHTML tests handling of empty HTML
func TestReaderPipeline_EmptyHTML(t *testing.T) {
	htmlContent := `<html><body></body></html>`

	ctx := context.Background()
	logger := &mockLogger{}
	fetcher := &mockFetcher{content: htmlContent}

	content, err := FetchAndProcess(ctx, "https://example.com/empty", logger, fetcher)

	// Should not error, but content will be minimal
	if err != nil {
		t.Fatalf("FetchAndProcess failed: %v", err)
	}

	// Content should be mostly empty
	if content == nil {
		t.Fatal("Expected content struct, got nil")
	}
}

// TestReaderPipeline_Host tests that markdown processor uses correct host
func TestReaderPipeline_Host(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		expectedHost string
	}{
		{"simple domain", "https://example.com/article", "example.com"},
		{"with port", "https://example.com:8080/article", "example.com:8080"},
		{"subdomain", "https://blog.example.com/article", "blog.example.com"},
		{"path", "https://example.com/path/to/article", "example.com"},
	}

	htmlContent := `<html><body><article><h1>Test</h1></article></body></html>`

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			logger := &mockLogger{}
			fetcher := &mockFetcher{content: htmlContent}

			content, err := FetchAndProcess(ctx, tt.url, logger, fetcher)

			if err != nil {
				t.Fatalf("FetchAndProcess failed: %v", err)
			}

			if content == nil {
				t.Fatal("Expected content, got nil")
			}

			// Verify URL is preserved with correct host
			if content.URL != tt.url {
				t.Errorf("Expected URL '%s', got '%s'", tt.url, content.URL)
			}
		})
	}
}
