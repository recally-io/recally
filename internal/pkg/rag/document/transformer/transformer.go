package transformer

import (
	"vibrain/internal/pkg/rag/document"
)

type Transformer func([]document.Document) ([]document.Document, error)

