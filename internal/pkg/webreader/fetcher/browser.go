package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"recally/internal/pkg/config"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

var DefaultBrowserFetcher *BrowserFetcher

func init() {
	var err error
	DefaultBrowserFetcher, err = defaultBrowserFetcher()
	if err != nil {
		logger.Default.Error("create default browser fetcher", "err", err)
	}
}

// BrowserConfig extends the base Config with browser-specific options
type BrowserConfig struct {
	Timeout        int    `json:"timeout"`          // Timeout in seconds
	ControlURL     string `json:"control_url"`      // Chrome DevTools Protocol control URL
	UserAgent      string `json:"user_agent"`       // User agent string
	ScrollToBottom bool   `json:"scroll_to_bottom"` // Scroll to bottom before extracting content
}

// BrowserFetcher implements the Fetcher interface using Chrome via go-rod
type BrowserFetcher struct {
	mux    sync.Mutex
	config BrowserConfig
}

func (f *BrowserFetcher) loadBrowser() (*rod.Browser, error) {
	// https://go-rod.github.io/#/custom-launch?id=remotely-manage-the-launcher
	l, err := launcher.NewManaged(f.config.ControlURL)
	if err != nil {
		return nil, fmt.Errorf("create new launcher: %w", err)
	}
	l.Headless(true).Set("disable-gpu").Set("no-sandbox").Set("disable-dev-shm-usage")
	browser := rod.New().Client(l.MustClient()).MustConnect()
	return browser, nil
}

func NewDefaultBrowserConfig() BrowserConfig {
	return BrowserConfig{
		Timeout:        120,
		ControlURL:     config.Settings.BrowserControlUrl,
		ScrollToBottom: true,
	}
}

func defaultBrowserFetcher(opts ...BroswerOption) (*BrowserFetcher, error) {
	config := NewDefaultBrowserConfig()
	return NewBrowserFetcher(config, opts...)
}

// NewBrowserFetcher creates a new BrowserFetcher with the given options
func NewBrowserFetcher(config BrowserConfig, opts ...BroswerOption) (*BrowserFetcher, error) {
	// Apply all options
	for _, opt := range opts {
		opt(&config)
	}
	return &BrowserFetcher{
		config: config,
	}, nil
}

// Fetch implements the Fetcher interface
func (f *BrowserFetcher) Fetch(ctx context.Context, url string) (*webreader.FetchedContent, error) {
	f.mux.Lock()
	defer f.mux.Unlock()
	browser, err := f.loadBrowser()
	if err != nil {
		return nil, fmt.Errorf("create new fetcher: %w", err)
	}
	defer func() { _ = browser.Close() }()

	// Create new page
	page := browser.MustPage(url)
	defer page.MustClose()

	if f.config.UserAgent != "" {
		_ = page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: f.config.UserAgent,
		})
	}

	// Wait for page to be ready
	page.MustWaitLoad()

	// Scroll to bottom if configured
	if f.config.ScrollToBottom {
		page.Mouse.MustScroll(0, 9999) // Large value to ensure bottom
		time.Sleep(time.Second)        // Wait for dynamic content
	}

	html, err := page.HTML()
	if err != nil {
		return nil, fmt.Errorf("get HTML: %w", err)
	}

	// Create result
	return &webreader.FetchedContent{
		Reader:      io.NopCloser(strings.NewReader(html)),
		StatusCode:  http.StatusOK,
		ContentType: "text/html",

		Content: webreader.Content{
			URL:  url,
			Html: html,
		},
	}, nil
}

// Close implements the Fetcher interface
func (f *BrowserFetcher) Close() error {
	return nil
}

// BroswerOption defines a function type for configuring BrowserFetcher
type BroswerOption func(*BrowserConfig)

// WithBroswerOptionTimeout sets the timeout for browser operations
func WithBroswerOptionTimeout(timeout int) BroswerOption {
	return func(c *BrowserConfig) {
		c.Timeout = timeout
	}
}

// WithBroswerOptionControlURL sets the Chrome DevTools Protocol control URL
func WithBroswerOptionControlURL(url string) BroswerOption {
	return func(c *BrowserConfig) {
		c.ControlURL = url
	}
}

// WithBroswerOptionUserAgent sets the user agent string
func WithBroswerOptionUserAgent(userAgent string) BroswerOption {
	return func(c *BrowserConfig) {
		c.UserAgent = userAgent
	}
}

// WithBroswerOptionScrollToBottom sets whether to scroll to bottom before extracting content
func WithBroswerOptionScrollToBottom(scroll bool) BroswerOption {
	return func(c *BrowserConfig) {
		c.ScrollToBottom = scroll
	}
}
