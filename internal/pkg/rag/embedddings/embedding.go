package embedddings

import (
	"context"
	"fmt"
	"recally/internal/pkg/rag/document"

	"github.com/sashabaranov/go-openai"
)

type Embedder struct {
	*openai.Client
}

func NewEmbeder(client *openai.Client) *Embedder {
	embeder := &Embedder{Client: client}

	return embeder
}

func (e *Embedder) EmbedTexts(ctx context.Context, texts []string, opts ...Option) ([][]float32, error) {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	req := options.ToEmbeddingRequest()
	req.Input = texts

	resp, err := e.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create embeddings: %w", err)
	}

	embs := make([][]float32, len(resp.Data))
	for i, emb := range resp.Data {
		embs[i] = emb.Embedding
	}

	return embs, nil
}

func (e *Embedder) EmbedDocuments(ctx context.Context, docs []document.Document) ([]document.Document, error) {
	texts := make([]string, len(docs))
	for i, doc := range docs {
		texts[i] = doc.Content
	}
	embs, err := e.EmbedTexts(ctx, texts)
	if err != nil {
		return nil, err
	}
	newDocs := make([]document.Document, len(docs))
	for i, doc := range docs {
		doc.Embedding = embs[i]
		newDocs[i] = doc
	}
	return newDocs, nil
}

func (e *Embedder) EmbedText(ctx context.Context, text string) ([]float32, error) {
	embs, err := e.EmbedTexts(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	return embs[0], nil
}
