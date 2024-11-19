package jinareader

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/logger"
	"vibrain/internal/pkg/tools"
)

const (
	jinaHost = "https://r.jina.ai"
)

type Tool struct {
	tools.BaseTool
	httpClient *http.Client
}

type RequestArgs struct {
	Url     string   `json:"url" jsonschema_description:"The URL to read."`
	Formats []string `json:"formats,omitempty" jsonschema_default:"markdown" jsonschema_description:"The content formats to return in the response, supported values: text, html, markdown, screenshot."`
}

type Content struct {
	Url           string `json:"url"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Description   string `json:"description,omitempty"`
	Text          string `json:"text,omitempty"`
	Html          string `json:"html,omitempty"`
	ScreenshotUrl string `json:"screenshotUrl,omitempty"`
}

type Result struct {
	Code   int     `json:"code"`
	Status float64 `json:"status"`
	Data   Content `json:"data"`
}

func New() *Tool {
	return &Tool{
		BaseTool: tools.BaseTool{
			Name:        "jinareader",
			Description: "Get LLM-friendly content from a URL.",
			Parameters:  &RequestArgs{},
		},
		httpClient: &http.Client{
			Timeout: 5 * 60 * time.Second,
		},
	}
}

func (t *Tool) Invoke(ctx context.Context, args string) (string, error) {
	var params RequestArgs
	if err := t.UnmarshalArgs(ctx, args, &params); err != nil {
		return "", err
	}

	result, err := t.ReadWithCache(ctx, params, cache.MemCache)
	if err != nil {
		return "", fmt.Errorf("failed to invoke tool: %w", err)
	}
	return t.MarshalResult(ctx, result)
}

func (t *Tool) Read(ctx context.Context, args RequestArgs) (*Content, error) {
	url := args.Url
	if url == "" {
		return nil, fmt.Errorf("jina reader: url is empty")
	}
	url = fmt.Sprintf("%s/%s", jinaHost, url)
	logger.FromContext(ctx).Debug("jina reader start", "url", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create jina reader request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Locale", "en-US")

	if len(args.Formats) == 0 {
		args.Formats = []string{"markdown"}
	}
	req.Header.Set("X-Return-Format", strings.Join(args.Formats, ","))

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send jina reader request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read jina reader response body: %w", err)
		}
		logger.FromContext(ctx).Error("jina reader read url error", "url", url, "status", resp.Status, "err", string(respData))
		return nil, fmt.Errorf("jina reader read URL '%s' Error: %s, %v", url, resp.Status, string(respData))
	}

	content := &Result{}

	if err := json.NewDecoder(resp.Body).Decode(content); err != nil {
		return nil, fmt.Errorf("failed to decode jina reader response: %w", err)
	}
	return &content.Data, nil
}

// ReadWithCache reads content from a URL with cache.
func (t *Tool) ReadWithCache(ctx context.Context, args RequestArgs, cacheService cache.Cache) (*Content, error) {
	// get result from cache
	cacheKey := cache.NewCacheKey("JinaReader", args.Url)
	if val, ok := cache.Get[Content](ctx, cacheService, cacheKey); ok {
		logger.FromContext(ctx).Info("JinaReader", "cache", "hit", "url", args.Url)
		return val, nil
	}
	data, err := t.Read(ctx, args)
	if err != nil {
		return nil, err
	}
	// set cache
	cacheService.Set(cacheKey, data, 24*time.Hour)
	logger.FromContext(ctx).Info("JinaReader", "cache", "set", "url", args.Url)
	return data, nil
}
