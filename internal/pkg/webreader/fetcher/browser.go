package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"vibrain/internal/pkg/webreader"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

const (
	defaultControlURL = "http://localhost:9222"
)

// BrowserConfig extends the base Config with browser-specific options
type BrowserConfig struct {
	Timeout        int    `json:"timeout"`          // Timeout in seconds
	ControlURL     string `json:"control_url"`      // Chrome DevTools Protocol control URL
	UserAgent      string `json:"user_agent"`       // User agent string
	ScrollToBottom bool   `json:"scroll_to_bottom"` // Scroll to bottom before extracting content
}

// BrowserFetcher implements the Fetcher interface using Chrome via go-rod
type BrowserFetcher struct {
	config  BrowserConfig
	browser *rod.Browser
}

func NewDefaultBrowserConfig() BrowserConfig {
	return BrowserConfig{
		Timeout:        30,
		ControlURL:     defaultControlURL,
		ScrollToBottom: true,
	}
}

// NewBrowserFetcher creates a new BrowserFetcher with the given options
func NewBrowserFetcher(opts ...BroswerOption) (*BrowserFetcher, error) {
	// Start with default configuration
	config := NewDefaultBrowserConfig()

	// Apply all options
	for _, opt := range opts {
		opt(&config)
	}
	return newBrowserFetcher(config)
}

func newBrowserFetcher(cfg BrowserConfig) (*BrowserFetcher, error) {
	// Get WebSocket debugger URL
	wsURL, err := getWebSocketDebuggerURL(cfg.ControlURL)
	if err != nil {
		return nil, fmt.Errorf("get debugger URL: %w", err)
	}

	// Connect to browser
	browser := rod.New().ControlURL(wsURL).MustConnect()

	return &BrowserFetcher{
		config:  cfg,
		browser: browser,
	}, nil
}

// getWebSocketDebuggerURL retrieves the WebSocket debugger URL from Chrome
func getWebSocketDebuggerURL(controlURL string) (string, error) {
	resp, err := http.Get(controlURL + "/json/version")
	if err != nil {
		return "", fmt.Errorf("get version info: %w", err)
	}
	defer resp.Body.Close()

	var info versionInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", fmt.Errorf("decode version info: %w", err)
	}

	if info.WebSocketDebuggerUrl == "" {
		return "", fmt.Errorf("no WebSocket debugger URL found")
	}

	return info.WebSocketDebuggerUrl, nil
}

// Fetch implements the Fetcher interface
func (f *BrowserFetcher) Fetch(ctx context.Context, url string) (*webreader.Content, error) {
	var err error
	f, err = newBrowserFetcher(f.config)
	if err != nil {
		return nil, fmt.Errorf("create new fetcher: %w", err)
	}

	// Create new page
	page := f.browser.MustPage(url)
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
	return &webreader.Content{
		URL:         url,
		Content:     io.NopCloser(strings.NewReader(html)),
		StatusCode:  http.StatusOK,
		ContentType: "text/html",
	}, nil
}

// Close implements the Fetcher interface
func (f *BrowserFetcher) Close() error {
	return f.browser.Close()
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

// versionInfo represents Chrome DevTools Protocol version info
type versionInfo struct {
	Browser              string `json:"Browser"`
	WebSocketDebuggerUrl string `json:"webSocketDebuggerUrl"`
}
