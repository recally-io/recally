package llms

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type LLM struct {
	client *openai.Client
}

func New(client *openai.Client) *LLM {
	return &LLM{
		client: client,
	}
}

func (l *LLM) GenerateContent(ctx context.Context, messages []openai.ChatCompletionMessage, options ...Option) (openai.ChatCompletionChoice, error) {
	opts := &Options{}
	for _, o := range options {
		o(opts)
	}
	req := opts.ToChatCompletionRequest()
	req.Messages = messages

	resp, err := l.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return openai.ChatCompletionChoice{}, err
	}

	return resp.Choices[0], nil
}
