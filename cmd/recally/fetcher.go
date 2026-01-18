package main

import (
	"fmt"
	"recally/internal/pkg/webreader"
	"recally/internal/pkg/webreader/fetcher"
)

// NewFetcher creates appropriate fetcher based on mode
// useBrowser: if true, creates browser fetcher; otherwise HTTP fetcher
// browserURL: Chrome DevTools Protocol control URL (only used for browser mode)
func NewFetcher(useBrowser bool, browserURL string) (webreader.Fetcher, error) {
	if useBrowser {
		return newBrowserFetcher(browserURL)
	}
	return newHTTPFetcher()
}

// newHTTPFetcher creates HTTP fetcher with standard http.Client
// Note: fetcher.NewHTTPFetcher uses session.New() internally, but session package
// only depends on log/slog which is acceptable. The plan notes session usage is OK.
func newHTTPFetcher() (webreader.Fetcher, error) {
	// Create HTTP fetcher with explicit config
	// Timeout: 120 seconds (2 minutes)
	// User-Agent: recally/{version}
	f, err := fetcher.NewHTTPFetcher(
		fetcher.WithHTTPOptionTimeout(120),
		fetcher.WithHTTPOptionFollowRedirects(true),
		fetcher.WithHTTPOptionMaxRedirects(10),
	)
	if err != nil {
		return nil, fmt.Errorf("create HTTP fetcher: %w", err)
	}

	return f, nil
}

// newBrowserFetcher creates browser fetcher with explicit config
// browserURL: Chrome DevTools Protocol control URL
func newBrowserFetcher(browserURL string) (webreader.Fetcher, error) {
	if browserURL == "" {
		return nil, fmt.Errorf("browser URL is empty")
	}

	// Construct browser config explicitly (never use DefaultBrowserFetcher)
	config := fetcher.BrowserConfig{
		Timeout:        120, // 2 minutes
		ControlURL:     browserURL,
		UserAgent:      fmt.Sprintf("recally/%s", version),
		ScrollToBottom: true, // For JavaScript-heavy sites
	}

	// Create browser fetcher with explicit config
	f, err := fetcher.NewBrowserFetcher(config)
	if err != nil {
		return nil, fmt.Errorf("create browser fetcher: %w", err)
	}

	return f, nil
}
