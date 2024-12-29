package transformer

import (
	"context"
	"fmt"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/rag/document"
)

const summaryPrompt = "Summarize the following text in a concise way:"

func WithLLMSummarier(llm *llms.LLM, opts ...llms.Option) Transformer {
	return func(docs []document.Document) ([]document.Document, error) {
		return Batch(func(d []document.Document) ([]document.Document, error) {
			ctx := context.Background()
			newDocs := make([]document.Document, 0, len(d))
			for i, doc := range d {
				prompt := fmt.Sprintf("%s\n####\n%s", summaryPrompt, doc.Content)
				resp, err := llm.TextCompletion(ctx, prompt, opts...)
				if err != nil {
					return nil, err
				}

				doc.Summaries = resp
				newDocs[i] = doc
			}
			return newDocs, nil
		}, docs)
	}
}
