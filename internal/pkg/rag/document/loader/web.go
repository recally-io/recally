package loader

import (
	"context"
	"vibrain/internal/core/workers"
	"vibrain/internal/pkg/rag/document"
	"vibrain/internal/pkg/rag/document/transformer"
)

// Text loads text data from an io.Reader.
type WebLoader struct {
	url string
}

// NewText creates a new text loader with an io.Reader.
func NewWebLoader(url string) WebLoader {
	return WebLoader{
		url: url,
	}
}

// Load reads from the io.Reader and returns documents with the data.
func (r WebLoader) Load(ctx context.Context, transformers ...transformer.Transformer) ([]document.Document, error) {
	content, err := workers.New().WebReader(ctx, r.url)
	if err != nil {
		return nil, err
	}

	docs := []document.Document{
		{
			Content: content.Content,
			Metadata: map[string]any{
				"url":         content.Url,
				"title":       content.Title,
				"description": content.Description,
			},
		},
	}
	return transformerPipeline(docs, transformers...)
}
