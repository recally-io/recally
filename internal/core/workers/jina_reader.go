package workers

import (
	"context"
	"time"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/logger"
	"vibrain/internal/pkg/tools/jinareader"
	"vibrain/internal/pkg/tools/jinasearcher"
)

const (
	webReaderCacheDomian   = "webreader"
	webSearcherCacheDomian = "websearcher"
	WebSummaryCacheDomian  = "websummary"
)

func (w *Worker) WebReader(ctx context.Context, url string) (*jinareader.Content, error) {
	// get result from cache
	cacheKey := cache.NewCacheKey(webReaderCacheDomian, url)
	if w.cache != nil {
		if val, ok := cache.Get[jinareader.Content](ctx, w.cache, cacheKey); ok {
			logger.FromContext(ctx).Info("WebReader", "cache", "hit", "url", url)
			return val, nil
		}
	}
	reader := jinareader.New()
	data, err := reader.Read(ctx, jinareader.RequestArgs{
		Url: url,
	})
	if err != nil {
		return nil, err
	}
	// set cache
	if w.cache != nil {
		w.cache.Set(cacheKey, data, 24*time.Hour)
		logger.FromContext(ctx).Info("WebReader", "cache", "set", "url", url)
	}
	return data, nil
}

func (w *Worker) WebSearcher(ctx context.Context, query string) ([]*jinasearcher.Content, error) {
	// get result from cache
	cacheKey := cache.NewCacheKey(webSearcherCacheDomian, query)
	if w.cache != nil {
		if val, ok := cache.Get[[]*jinasearcher.Content](ctx, w.cache, cacheKey); ok {
			logger.FromContext(ctx).Info("WebSearcher", "cache", "hit", "query", query)
			return *val, nil
		}
	}

	searcher := jinasearcher.New()
	data, err := searcher.Search(ctx, jinasearcher.RequestArgs{
		Query: query,
	})
	if err != nil {
		return nil, err
	}

	if w.cache != nil {
		w.cache.SetWithContext(ctx, cacheKey, data, 24*time.Hour)
		logger.FromContext(ctx).Info("WebSearcher", "cache", "set", "query", query)
	}

	return data, nil
}
