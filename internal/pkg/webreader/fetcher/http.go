package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"recally/internal/pkg/session"
	"recally/internal/pkg/webreader"
	"time"

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

// HTTPOption is a function type that modifies HTTPFetcher options
type HTTPOption func(*HTTPConfig)

// HTTPFetcher implements the Fetcher interface using net/http
type HTTPFetcher struct {
	client *http.Client
	config HTTPConfig
}

func DefaultHTTPConfig() HTTPConfig {
	return HTTPConfig{
		Timeout:         30,               // Default 30 seconds timeout
		MaxBodySize:     10 * 1024 * 1024, // Default 10MB
		FollowRedirects: true,             // Default to following redirects
	}
}

// NewHTTPFetcher creates a new HTTPFetcher with the given options
func NewHTTPFetcher(opts ...HTTPOption) (*HTTPFetcher, error) {
	config := DefaultHTTPConfig()

	// Apply all options
	for _, opt := range opts {
		opt(&config)
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
	}, nil
}

// Fetch implements the Fetcher interface
func (f *HTTPFetcher) Fetch(ctx context.Context, url string) (*webreader.FetchedContent, error) {
	var lastErr error
	retries := f.config.RetryCount + 1

	for attempt := range retries {
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
func (f *HTTPFetcher) doFetch(ctx context.Context, url string) (*webreader.FetchedContent, error) {
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
		_ = resp.Body.Close()
		return nil, fmt.Errorf("http status %d: %s", resp.StatusCode, resp.Status)
	}

	return &webreader.FetchedContent{
		Reader:      resp.Body,
		StatusCode:  http.StatusOK,
		ContentType: resp.Header.Get("Content-Type"),
		Headers:     resp.Header,

		Content: webreader.Content{
			URL: url,
		},
	}, nil
}

// Close implements the Fetcher interface
func (f *HTTPFetcher) Close() error {
	return nil
}

// WithHTTPOptionTimeout sets the timeout for HTTP requests
func WithHTTPOptionTimeout(timeout int) HTTPOption {
	return func(config *HTTPConfig) {
		config.Timeout = timeout
	}
}

// WithHTTPOptionMaxBodySize sets the maximum body size for HTTP responses
func WithHTTPOptionMaxBodySize(size int64) HTTPOption {
	return func(config *HTTPConfig) {
		config.MaxBodySize = size
	}
}

// WithHTTPOptionMaxRedirects sets the maximum number of redirects to follow
func WithHTTPOptionMaxRedirects(maxRedirects int) HTTPOption {
	return func(config *HTTPConfig) {
		config.MaxRedirects = maxRedirects
	}
}

// WithHTTPOptionRetryCount sets the number of times to retry failed requests
func WithHTTPOptionRetryCount(retryCount int) HTTPOption {
	return func(config *HTTPConfig) {
		config.RetryCount = retryCount
	}
}

// WithRetryDelay sets the delay between retries
func WithRetryDelay(delay time.Duration) HTTPOption {
	return func(config *HTTPConfig) {
		config.RetryDelay = delay
	}
}

// WithHTTPOptionExtraHeaders sets additional HTTP headers
func WithHTTPOptionExtraHeaders(headers map[string]string) HTTPOption {
	return func(config *HTTPConfig) {
		config.ExtraHeaders = headers
	}
}

// WithHTTPOptionFollowRedirects sets whether to follow redirects
func WithHTTPOptionFollowRedirects(follow bool) HTTPOption {
	return func(config *HTTPConfig) {
		config.FollowRedirects = follow
	}
}
