package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"vibrain/internal/pkg/session"
	"vibrain/internal/pkg/webreader"

	utls "github.com/refraction-networking/utls"
)

// HTTPConfig extends the base Config with HTTP-specific options
type HTTPConfig struct {
	Timeout         int               `json:"timeout"`          // Timeout in seconds
	MaxBodySize     int64             `json:"max_body_size"`    // Maximum body size in bytes
	MaxRedirects    int               `json:"max_redirects"`    // Maximum number of redirects to follow
	RetryCount      int               `json:"retry_count"`      // Number of times to retry failed requests
	RetryDelay      time.Duration     `json:"retry_delay"`      // Delay between retries
	ExtraHeaders    map[string]string `json:"extra_headers"`    // Additional HTTP headers to include
	FollowRedirects bool              `json:"follow_redirects"` // Whether to follow redirects
}

// HTTPFetcher implements the Fetcher interface using net/http
type HTTPFetcher struct {
	client *http.Client
	config HTTPConfig
	closed bool
}

// NewHTTPFetcher creates a new HTTPFetcher with the given configuration
func NewHTTPFetcher(config HTTPConfig) *HTTPFetcher {
	if config.Timeout == 0 {
		config.Timeout = 30 // Default 30 seconds timeout
	}
	if config.MaxBodySize == 0 {
		config.MaxBodySize = 10 * 1024 * 1024 // Default 10MB
	}

	sess := session.New(session.WithClientHelloID(utls.HelloChrome_100_PSK))

	client := sess.Client
	client.Timeout = time.Duration(config.Timeout) * time.Second

	// Configure redirect policy if needed
	if !config.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else if config.MaxRedirects > 0 {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= config.MaxRedirects {
				return fmt.Errorf("stopped after %d redirects", config.MaxRedirects)
			}
			return nil
		}
	}

	return &HTTPFetcher{
		client: client,
		config: config,
	}
}

// Fetch implements the Fetcher interface
func (f *HTTPFetcher) Fetch(ctx context.Context, url string) (*webreader.Content, error) {
	if f.closed {
		return nil, fmt.Errorf("fetcher is closed")
	}

	var lastErr error
	retries := f.config.RetryCount + 1

	for attempt := 0; attempt < retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(f.config.RetryDelay):
			}
		}

		result, err := f.doFetch(ctx, url)
		if err == nil {
			return result, nil
		}

		lastErr = err
		// Don't retry on context cancellation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("all fetch attempts failed: %w", lastErr)
}

// doFetch performs the actual HTTP fetch
func (f *HTTPFetcher) doFetch(ctx context.Context, url string) (*webreader.Content, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}

	// Check if the response status code indicates an error
	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, fmt.Errorf("http status %d: %s", resp.StatusCode, resp.Status)
	}

	return &webreader.Content{
		URL:         url,
		Content:     resp.Body,
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		Headers:     resp.Header,
	}, nil
}

// Close implements the Fetcher interface
func (f *HTTPFetcher) Close() error {
	if f.closed {
		return nil
	}
	f.closed = true
	f.client.CloseIdleConnections()
	return nil
}
