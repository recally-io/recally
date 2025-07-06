package loader

import (
	"bytes"
	"context"
	"io"
	"recally/internal/pkg/rag/document"
	"recally/internal/pkg/rag/document/transformer"
)

// TextLoader loads text data from an io.Reader.
type TextLoader struct {
	io.Reader
}

// NewText creates a new text loader with an io.Reader.
func NewTextLoader(r io.Reader) TextLoader {
	return TextLoader{
		Reader: r,
	}
}

// Load reads from the io.Reader and returns documents with the data.
func (r TextLoader) Load(ctx context.Context, transformers ...transformer.Transformer) ([]document.Document, error) {
	buf := new(bytes.Buffer)

	_, err := io.Copy(buf, r.Reader)
	if err != nil {
		return nil, err
	}

	docs := []document.Document{
		{
			Content:  buf.String(),
			Metadata: map[string]any{},
		},
	}

	return transformerPipeline(docs, transformers...)
}
