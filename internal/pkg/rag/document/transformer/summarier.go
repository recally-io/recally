package transformer

import (
	"context"
	"fmt"
	"vibrain/internal/pkg/llms"
	"vibrain/internal/pkg/rag/document"

	"github.com/sashabaranov/go-openai"
)

const summaryPrompt = "Summarize the following text in a concise way:"

func WithLLMSummarier(llm *llms.LLM, opts ...llms.Option) Transformer {
	return func(docs []document.Document) ([]document.Document, error) {
		return Batch(func(d []document.Document) ([]document.Document, error) {
			ctx := context.Background()
			newDocs := make([]document.Document, 0, len(d))
			for i, doc := range d {
				messages := []openai.ChatCompletionMessage{
					{
						Role:    "user",
						Content: fmt.Sprintf("%s\n####\n%s", summaryPrompt, doc.Content),
					},
				}
				resp, err := llm.GenerateContent(ctx, messages, opts...)
				if err != nil {
					return nil, err
				}

				doc.Summaries = resp.Message.Content
				newDocs[i] = doc
			}
			return newDocs, nil
		}, docs)
	}
}
