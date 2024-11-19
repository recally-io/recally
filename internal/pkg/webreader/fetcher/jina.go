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
)

const (
	jinaHost = "https://r.jina.ai"
)

// JinaConfig extends the base Config with Jina-specific options
type JinaConfig struct {
	Timeout int `json:"timeout"` // Timeout in seconds
}

// jinaContent represents the content returned by Jina API
type jinaContent struct {
	Url           string `json:"url"`
	Title         string `json:"title"`
	Content       string `json:"content"` // markdown content
	Description   string `json:"description,omitempty"`
	Text          string `json:"text,omitempty"`
	Html          string `json:"html,omitempty"`
	ScreenshotUrl string `json:"screenshotUrl,omitempty"`
}

// jinaResponse represents the API response from Jina
type jinaResponse struct {
	Code   int         `json:"code"`
	Status float64     `json:"status"`
	Data   jinaContent `json:"data"`
}

// JinaFetcher implements the Fetcher interface using Jina.ai reader
type JinaFetcher struct {
	client *http.Client
	config JinaConfig
	closed bool
}

func DefaultJinaConfig() JinaConfig {
	return JinaConfig{
		Timeout: 30, // Default 30 Seconds
	}
}

type JinaOption func(*JinaConfig)

func WithJinaOptionTimeout(timeout int) JinaOption {
	return func(config *JinaConfig) {
		config.Timeout = timeout
	}
}

// NewJinaFetcher creates a new JinaFetcher with the given configuration
func NewJinaFetcher(opts ...JinaOption) *JinaFetcher {
	config := DefaultJinaConfig()

	// Apply all options
	for _, opt := range opts {
		opt(&config)
	}

	return &JinaFetcher{
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
		config: config,
	}
}

// Fetch implements the Fetcher interface
func (f *JinaFetcher) Fetch(ctx context.Context, url string) (*webreader.Content, error) {
	if f.closed {
		return nil, fmt.Errorf("fetcher is closed")
	}

	// Prepare the Jina API URL
	jinaURL := fmt.Sprintf("%s/%s", jinaHost, url)

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jinaURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Locale", "en-US")
	req.Header.Set("X-Return-Format", "html")

	// Send request
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("jina API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	// Parse response
	var jinaResp jinaResponse
	if err := json.NewDecoder(resp.Body).Decode(&jinaResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &webreader.Content{
		URL:         url,
		Title:       jinaResp.Data.Title,
		Description: jinaResp.Data.Description,
		Text:        jinaResp.Data.Text,
		Html:        jinaResp.Data.Html,
		Image:       jinaResp.Data.ScreenshotUrl,
		Content:     io.NopCloser(strings.NewReader(jinaResp.Data.Html)),
	}, nil
}

// Close implements the Fetcher interface
func (f *JinaFetcher) Close() error {
	if f.closed {
		return nil
	}
	f.closed = true
	f.client.CloseIdleConnections()
	return nil
}
