package document

import "github.com/google/uuid"

// Document is the struct for interacting with a document.
type Document struct {
	ID        uuid.UUID      `json:"id"`
	Content   string         `json:"content"`
	Metadata  map[string]any `json:"metadata"`
	Embedding []float32      `json:"embedding"`
	Summaries string         `json:"summaries"`
	Score     float32        `json:"score"`
}
