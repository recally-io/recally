package main

import (
	"context"
	"fmt"
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/llms"
	"vibrain/internal/pkg/logger"

	"github.com/sashabaranov/go-openai"
)

func main() {
	ctx := context.Background()
	llm := llms.New(config.Settings.OpenAI.BaseURL, config.Settings.OpenAI.ApiKey)
	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "You are a helpful assistant.",
		},
		{
			Role:    "user",
			Content: "What's news in hackernews?",
		},
	}

	resp, usage, err := llm.GenerateContent(ctx, messages)
	if err != nil {
		logger.Default.Error("failed to generate content", "error", err)
		return
	}
	logger.Default.Info("generated content", "response", resp, "usage", usage)
	fmt.Println(resp.Message.Content)
}
