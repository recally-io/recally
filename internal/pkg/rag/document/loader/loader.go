package loader

import (
	"context"

	"recally/internal/pkg/rag/document"
	"recally/internal/pkg/rag/document/transformer"
	"recally/internal/pkg/rag/textsplitter"
)

type Options struct {
	Splitter textsplitter.TextSplitter
}

func DefaultOptions() Options {
	return Options{
		Splitter: textsplitter.NewRecursiveCharacter(),
	}
}

type Loader interface {
	Load(ctx context.Context, transformers ...transformer.Transformer) ([]document.Document, error)
}

func transformerPipeline(docs []document.Document, transformers ...transformer.Transformer) ([]document.Document, error) {
	var err error
	for _, t := range transformers {
		docs, err = t(docs)
		if err != nil {
			return nil, err
		}
	}

	return docs, nil
}
