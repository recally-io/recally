package workers

import (
	"context"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/tools/jinareader"
	"recally/internal/pkg/tools/jinasearcher"
)

const (
	webReaderCacheDomian   = "webreader"
	webSearcherCacheDomian = "websearcher"
	WebSummaryCacheDomian  = "websummary"
)

func (w *Worker) WebReader(ctx context.Context, url string) (*jinareader.Content, error) {
	reader := jinareader.New()
	if w.cache != nil {
		return reader.ReadWithCache(ctx, jinareader.RequestArgs{Url: url}, w.cache)
	}
	return reader.ReadWithCache(ctx, jinareader.RequestArgs{Url: url}, cache.MemCache)
}

func (w *Worker) WebSearcher(ctx context.Context, query string) ([]*jinasearcher.Content, error) {
	searcher := jinasearcher.New()
	if w.cache != nil {
		return searcher.SearchWithCache(ctx, jinasearcher.RequestArgs{Query: query}, w.cache)
	}
	return searcher.SearchWithCache(ctx, jinasearcher.RequestArgs{Query: query}, cache.MemCache)
}
