package llms

import (
	"context"
	"time"
	"vibrain/internal/pkg/logger"

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
	start := time.Now()
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
	logger.FromContext(ctx).Info("time for generated content",
		"duration", time.Since(start),
		"model", req.Model,
		"prompt_tokens", resp.Usage.PromptTokens,
		"completion_tokens", resp.Usage.CompletionTokens,
	)
	return resp.Choices[0], resp.Usage, nil
}
