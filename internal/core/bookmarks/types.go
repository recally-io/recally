package bookmarks

import (
	"context"
	"fmt"
	"vibrain/internal/pkg/tools/jinareader"

	"github.com/pgvector/pgvector-go"
)

// Embedder defines the interface for generating embeddings
type Embedder interface {
	CreateEmbeddings(ctx context.Context, text string) ([]float32, error)
}

// URLFetcher defines the interface for fetching URL content
type URLFetcher interface {
	Fetch(ctx context.Context, url string) (*jinareader.Content, error)
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

type JinaFetcherService struct {
	reader *jinareader.Tool
}

func NewJinaFetcher() *JinaFetcherService {
	return &JinaFetcherService{reader: jinareader.New()}
}

func (j *JinaFetcherService) Fetch(ctx context.Context, url string) (*jinareader.Content, error) {
	args := jinareader.RequestArgs{
		Url: url,
		Formats: []string{
			// "text",
			// "html",
			"markdown",
			// "screenshot",
		},
	}
	result, err := j.reader.Read(ctx, args)
	if err != nil {
		return nil, err
	}
	return result, nil
}
