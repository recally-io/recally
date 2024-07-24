package transformer

import (
	"context"
	"vibrain/internal/pkg/rag/document"
	"vibrain/internal/pkg/rag/embedddings"
)

func WithEmbedder(embedder embedddings.Embedder) Transformer {
	return func(docs []document.Document) ([]document.Document, error) {
		return embedder.EmbedDocuments(context.TODO(), docs)
	}
}
