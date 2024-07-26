package llms

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type LLM struct {
	client *openai.Client
}

func New(baseUrl, apiKey string) *LLM {
	cfg := openai.DefaultConfig(apiKey)
	if baseUrl != "" {
		cfg.BaseURL = baseUrl
	}
	return &LLM{
		client: openai.NewClientWithConfig(cfg),
	}
}

func (l *LLM) GenerateContent(ctx context.Context, messages []openai.ChatCompletionMessage, options ...Option) (openai.ChatCompletionChoice, openai.Usage, error) {
	opts := &Options{}
	for _, o := range options {
		o(opts)
	}
	req := opts.ToChatCompletionRequest()
	req.Messages = messages

	resp, err := l.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return openai.ChatCompletionChoice{}, openai.Usage{}, err
	}

	return resp.Choices[0], resp.Usage, nil
}
