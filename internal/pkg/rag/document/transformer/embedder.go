package transformer

import (
	"context"

	"recally/internal/pkg/rag/document"
	"recally/internal/pkg/rag/embedddings"
)

func WithEmbedder(embedder embedddings.Embedder) Transformer {
	return func(docs []document.Document) ([]document.Document, error) {
		return embedder.EmbedDocuments(context.TODO(), docs)
	}
}
