package transformer

import (
	"errors"
	"recally/internal/pkg/rag/document"
	"recally/internal/pkg/rag/textsplitter"
)

// ErrMismatchMetadatasAndText is returned when the number of texts and metadatas
// given to CreateDocuments does not match. The function will not error if the
// length of the metadatas slice is zero.
var ErrMismatchMetadatasAndText = errors.New("number of texts and metadatas does not match")

func WithTextSplitter(textSplitter textsplitter.TextSplitter) Transformer {
	return func(docs []document.Document) ([]document.Document, error) {
		return Batch(func(d []document.Document) ([]document.Document, error) {
			texts := make([]string, 0)
			metadatas := make([]map[string]any, 0)
			for _, document := range d {
				texts = append(texts, document.Content)
				metadatas = append(metadatas, document.Metadata)
			}

			return createDocuments(textSplitter, texts, metadatas)
		}, docs)
	}
}

func createDocuments(textSplitter textsplitter.TextSplitter, texts []string, metadatas []map[string]any) ([]document.Document, error) {
	if len(metadatas) == 0 {
		metadatas = make([]map[string]any, len(texts))
	}

	if len(texts) != len(metadatas) {
		return nil, ErrMismatchMetadatasAndText
	}

	documents := make([]document.Document, 0)

	for i := 0; i < len(texts); i++ {
		chunks, err := textSplitter.Split(texts[i])
		if err != nil {
			return nil, err
		}

		for _, chunk := range chunks {
			// Copy the document metadata
			curMetadata := make(map[string]any, len(metadatas[i]))
			for key, value := range metadatas[i] {
				curMetadata[key] = value
			}

			documents = append(documents, document.Document{
				Content:  chunk,
				Metadata: curMetadata,
			})
		}
	}

	return documents, nil
}
