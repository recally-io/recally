package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"recally/internal/pkg/config"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"

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
			Content: "Please tell me what in the page: https://labs.watchtowr.com/we-spent-20-to-achieve-rce-and-accidentally-became-the-admins-of-mobi/",
		},
	}

	sendToUser := func(m llms.StreamingMessage) {
		choice := m.Choice
		err := m.Err
		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.Default.Info("done")
				return
			} else {
				logger.Default.Error("failed to send message to user", "error", err)
				return
			}
		}
		fmt.Print(choice.Message.Content)
	}

	llm.GenerateContent(ctx, messages, sendToUser,
		llms.WithStream(true),
		// llms.WithToolNames([]string{"googlesearch", "jinareader"}),
	)
}
