package jinasearcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/logger"
	"vibrain/internal/pkg/tools"
)

const (
	jinaHost = "https://s.jina.ai"
)

type Tool struct {
	tools.BaseTool
	httpClient *http.Client
}

type RequestArgs struct {
	Query string `json:"query" jsonschema_description:"The query to search."`
}

type Content struct {
	Url         string `json:"url"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
}

type Result struct {
	Code   int        `json:"code"`
	Status float64    `json:"status"`
	Data   []*Content `json:"data"`
}

func New() *Tool {
	return &Tool{
		BaseTool: tools.BaseTool{
			Name:        "jinasearcher",
			Description: "Get LLM-friendly content from a web search.",
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

	result, err := t.SearchWithCache(ctx, params, cache.MemCache)
	if err != nil {
		return "", fmt.Errorf("failed to invoke tool: %w", err)
	}
	return t.MarshalResult(ctx, result)
}

func (t *Tool) Search(ctx context.Context, args RequestArgs) ([]*Content, error) {
	url := fmt.Sprintf("%s/%s", jinaHost, args.Query)
	logger.FromContext(ctx).Debug("jina searcher start", "query", args.Query)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create jina searcher request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send jina searcher request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read jina searcher response body: %w", err)
		}
		logger.FromContext(ctx).Error("jina searcher error", "query", args.Query, "status", resp.Status, "err", string(respData))
		return nil, fmt.Errorf("jina searcher for '%s' error: %s, %v", args.Query, resp.Status, string(respData))
	}

	content := &Result{}

	if err := json.NewDecoder(resp.Body).Decode(content); err != nil {
		return nil, fmt.Errorf("failed to decode jina searcher response: %w", err)
	}
	return content.Data, nil
}

// ReadWithCache reads content from a URL with cache.
func (t *Tool) SearchWithCache(ctx context.Context, args RequestArgs, cacheService cache.Cache) ([]*Content, error) {
	// get result from cache
	cacheKey := cache.NewCacheKey("JinaSearcher", args.Query)
	if val, ok := cache.Get[[]*Content](ctx, cacheService, cacheKey); ok {
		logger.FromContext(ctx).Info("JinaSearcher", "cache", "hit", "query", args.Query)
		return *val, nil
	}
	data, err := t.Search(ctx, args)
	if err != nil {
		return nil, err
	}
	// set cache
	cacheService.Set(cacheKey, data, 24*time.Hour)
	logger.FromContext(ctx).Info("JinaSearcher", "cache", "set", "query", args.Query)
	return data, nil
}
