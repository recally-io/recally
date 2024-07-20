package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/logger"
)

const (
	webReaderCacheDomian   = "webreader"
	webSearcherCacheDomian = "websearcher"
	WebSummaryCacheDomian  = "websummary"
)

type WebReaderContent struct {
	Url         string `json:"url"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
}

type jinaReaderResponse struct {
	Code   int              `json:"code"`
	Status float64          `json:"status"`
	Data   WebReaderContent `json:"data"`
}

type jinaSearcherResponse struct {
	Code   int                 `json:"code"`
	Status float64             `json:"status"`
	Data   []*WebReaderContent `json:"data"`
}

const (
	jinaReaderHost   = "https://r.jina.ai"
	jinaSearcherHost = "https://s.jina.ai"
)

func newHttpClient() *http.Client {
	return &http.Client{
		Timeout: 5 * 60 * time.Second,
	}
}

func (w *Worker) WebReader(ctx context.Context, url string) (*WebReaderContent, error) {
	// get result from cache
	cacheKey := cache.NewCacheKey(webReaderCacheDomian, url)
	if w.cache != nil {
		if val, ok := w.cache.GetWithContext(ctx, cacheKey); ok {
			logger.FromContext(ctx).Info("WebReader", "cache", "hit", "url", url)
			return val.(*WebReaderContent), nil
		}
	}
	readerUrl := fmt.Sprintf("%s/%s", jinaReaderHost, url)
	logger.FromContext(ctx).Info("WebReader", "url", readerUrl)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, readerUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := newHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		slog.Error("web reader read url error", "url", url, "status", resp.Status, "err", string(respData))
		return nil, fmt.Errorf("Read URL %s Error: %s, %v", readerUrl, resp.Status, string(respData))
	}

	content := &jinaReaderResponse{}

	if err := json.NewDecoder(resp.Body).Decode(content); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	// set cache
	if w.cache != nil {
		w.cache.Set(cacheKey, &content.Data, 24*time.Hour)
		logger.FromContext(ctx).Info("WebReader", "cache", "set", "url", url)
	}
	return &content.Data, nil
}

func (w *Worker) WebSearcher(ctx context.Context, query string) ([]*WebReaderContent, error) {
	// get result from cache
	cacheKey := cache.NewCacheKey(webSearcherCacheDomian, query)
	if w.cache != nil {
		if val, ok := w.cache.GetWithContext(ctx, cacheKey); ok {
			logger.FromContext(ctx).Info("WebSearcher", "cache", "hit", "query", query)
			var content []*WebReaderContent
			if err := json.Unmarshal(val.([]byte), &content); err != nil {
				return nil, fmt.Errorf("failed to unmarshal cache value: %w", err)
			}
			return content, nil
		}
	}
	searcherUrl := fmt.Sprintf("%s/%s", jinaSearcherHost, query)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searcherUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := newHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		slog.Error("web searcher error", "query", query, "status", resp.Status, "err", string(respData))
		return nil, fmt.Errorf("Search web Error: %s, %v", resp.Status, string(respData))
	}

	content := &jinaSearcherResponse{}

	if err := json.NewDecoder(resp.Body).Decode(content); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if w.cache != nil {
		w.cache.SetWithContext(ctx, cacheKey, content.Data, 24*time.Hour)
		logger.FromContext(ctx).Info("WebSearcher", "cache", "set", "query", query)
	}
	return content.Data, nil
}
