package main

import (
	"context"
	"fmt"
	"net/url"

	"recally/internal/pkg/webreader"
	"recally/internal/pkg/webreader/processor"
)

// FetchAndProcess fetches URL and processes it through the reader pipeline
// Returns the processed Content or error
//
// Pipeline stages:
// 1. Fetch raw HTML content using provided fetcher
// 2. Extract main content with Readability processor
// 3. Convert HTML to Markdown with Markdown processor
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - url: Target URL to fetch and process
//   - logger: Logger for pipeline operations (can be nil for no-op logger)
//   - fetcher: Fetcher implementation (HTTP or Browser)
//
// Returns:
//   - *webreader.Content: Processed content with all fields populated
//   - error: Fetch or processing errors
func FetchAndProcess(ctx context.Context, urlStr string, logger webreader.Logger, fetcher webreader.Fetcher) (*webreader.Content, error) {
	// Validate URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Ensure URL has http/https scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("unsupported URL scheme: %s (only http/https allowed)", parsedURL.Scheme)
	}

	// Extract host for markdown processor
	host := parsedURL.Host

	// Create processor chain
	// 1. Readability: Extracts main content, removes boilerplate
	//    - No hooks requiring S3/DB (hooks are commented out in processor)
	// 2. Markdown: Converts HTML to markdown
	//    - Uses builtin hooks only (e.g., WeChat hook)
	//    - No S3/DB dependencies (ImageHook not used)
	processors := []webreader.Processor{
		processor.NewReadabilityProcessor(),
		processor.NewMarkdownProcessor(host), // Only builtin hooks, no S3/DB
	}

	// Create reader with fetcher, logger, and processors
	reader := webreader.New(fetcher, logger, processors...)

	// Fetch and process content through pipeline
	content, err := reader.Read(ctx, urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch and process content: %w", err)
	}

	return content, nil
}
