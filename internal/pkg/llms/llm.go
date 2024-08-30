package llms

import (
	"context"
	"fmt"
	"time"
	"vibrain/internal/pkg/logger"
	"vibrain/internal/pkg/tools"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/sync/errgroup"
)

type LLM struct {
	client       *openai.Client
	toolMappings map[string]tools.Tool
}

func New(baseUrl, apiKey string) *LLM {
	cfg := openai.DefaultConfig(apiKey)
	if baseUrl != "" {
		cfg.BaseURL = baseUrl
	}
	return &LLM{
		client:       openai.NewClientWithConfig(cfg),
		toolMappings: AllToolMappings,
	}
}

func (l *LLM) ListModels(ctx context.Context) ([]string, error) {
	models, err := l.client.ListModels(ctx)
	if err != nil {
		return nil, err
	}
	data := make([]string, 0, len(models.Models))
	for _, m := range models.Models {
		data = append(data, m.ID)
	}
	return data, nil
}

func (l *LLM) CreateEmbeddings(ctx context.Context, text string) ([]float32, error) {
	embeddings, err := l.client.CreateEmbeddings(ctx, openai.EmbeddingRequestStrings{
		Input: []string{text},
		Model: openai.SmallEmbedding3,
	})
	if err != nil {
		return nil, err
	}
	return embeddings.Data[0].Embedding, nil
}

func (l *LLM) TextCompletion(ctx context.Context, prompt string, options ...Option) (string, error) {
	opts := &Options{}
	for _, o := range options {
		o(opts)
	}
	req := opts.ToChatCompletionRequest()
	req.Messages = []openai.ChatCompletionMessage{{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	}}
	choice, _, err := l.generateContent(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to generate text: %w", err)
	}
	return choice.Message.Content, nil
}

func (l *LLM) GenerateContent(ctx context.Context, messages []openai.ChatCompletionMessage, options ...Option) (openai.ChatCompletionChoice, openai.Usage, error) {
	opts := &Options{}
	for _, o := range options {
		o(opts)
	}
	req := opts.ToChatCompletionRequest()
	req.Messages = messages

	// dynamically add tools to the request
	logger.FromContext(ctx).Info("generating content", "tool_names", opts.ToolNames)
	if len(opts.ToolNames) > 0 {
		mapping := make(map[string]tools.Tool)
		for _, name := range opts.ToolNames {
			tool, ok := l.toolMappings[name]
			if ok {
				mapping[name] = tool
			}
			req.Tools = llmTools(mapping)
		}
	}

	choice, usage, err := l.generateContent(ctx, req)
	if err != nil {
		return choice, usage, err
	}
	if len(choice.Message.ToolCalls) > 0 {
		req.Messages = append(req.Messages, choice.Message)
		for {
			toolMessages, err := l.invokeTools(ctx, choice.Message.ToolCalls)
			if err != nil {
				return choice, usage, err
			}
			req.Messages = append(req.Messages, toolMessages...)
			choice, usage, err = l.generateContent(ctx, req)
			if err != nil {
				return choice, usage, err
			}
			if len(choice.Message.ToolCalls) == 0 {
				break
			}
		}
	}

	return choice, usage, nil
}

func (l *LLM) generateContent(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionChoice, openai.Usage, error) {
	start := time.Now()
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
	if len(resp.Choices) == 0 {
		return openai.ChatCompletionChoice{}, openai.Usage{}, fmt.Errorf("no choices returned")
	}
	return resp.Choices[0], resp.Usage, nil
}

func (l *LLM) invokeTools(ctx context.Context, toolCalls []openai.ToolCall) ([]openai.ChatCompletionMessage, error) {
	var messages []openai.ChatCompletionMessage
	eg, ctx := errgroup.WithContext(ctx)
	for _, tc := range toolCalls {
		eg.Go(func() error {
			toolName := tc.Function.Name
			toolArgs := tc.Function.Arguments
			tool, ok := l.toolMappings[toolName]
			if !ok {
				return fmt.Errorf("tool %s not found", toolName)
			}
			logger.FromContext(ctx).Info("invoking tool", "tool", toolName, "args", toolArgs)
			start := time.Now()
			toolResp, err := tool.Invoke(ctx, toolArgs)
			if err != nil {
				return fmt.Errorf("failed to invoke tool %s: %w", toolName, err)
			}
			logger.FromContext(ctx).Info("tool response", "tool", toolName, "response", toolResp[:200], "duration", time.Since(start))
			messages = append(messages, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				ToolCallID: tc.ID,
				Content:    toolResp,
			})
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return messages, nil
}
