package loader

import (
	"context"
	"recally/internal/pkg/rag/document"
	"recally/internal/pkg/rag/document/transformer"
	"recally/internal/pkg/tools/jinareader"
)

// Text loads text data from an io.Reader.
type WebLoader struct {
	url    string
	reader *jinareader.Tool
}

// NewText creates a new text loader with an io.Reader.
func NewWebLoader(url string) WebLoader {
	return WebLoader{
		url:    url,
		reader: jinareader.New(),
	}
}

// Load reads from the io.Reader and returns documents with the data.
func (r WebLoader) Load(ctx context.Context, transformers ...transformer.Transformer) ([]document.Document, error) {
	content, err := r.reader.Read(ctx, jinareader.RequestArgs{Url: r.url})
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
