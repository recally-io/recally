package bookmarks

import (
	"context"
	"fmt"
	"vibrain/internal/pkg/llms"
	"vibrain/internal/pkg/webreader"
	"vibrain/internal/pkg/webreader/fetcher"
	"vibrain/internal/pkg/webreader/processor"

	"github.com/pgvector/pgvector-go"
)

// LLM defines the interface for generating embeddings
type LLM interface {
	CreateEmbeddings(ctx context.Context, text string) ([]float32, error)
	TextCompletion(ctx context.Context, prompt string, options ...llms.Option) (string, error)
}

// UrlReader defines the interface for fetching URL content
type UrlReader interface {
	Read(ctx context.Context, url string) (*webreader.Content, error)
}

// Common errors
var (
	ErrNotFound     = fmt.Errorf("bookmark not found")
	ErrDuplicate    = fmt.Errorf("bookmark already exists")
	ErrInvalidInput = fmt.Errorf("invalid input")
	ErrUnauthorized = fmt.Errorf("unauthorized access")
)

// SearchParams encapsulates search parameters
type SearchParams struct {
	UserID    int32
	Query     string
	Embedding pgvector.Vector
	Limit     int32
	Offset    int32
}

func NewWebReader(llm *llms.LLM) (UrlReader, error) {
	fetcher, err := fetcher.NewHTTPFetcher()
	if err != nil {
		return nil, fmt.Errorf("create browser fetcher error: %w", err)
	}

	processors := []webreader.Processor{
		processor.NewMarkdownProcessor(),
		processor.NewSummaryProcessor(llm),
	}

	return webreader.New(fetcher, processors...), nil
}
