package llms

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
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

	respChan := make(chan *openai.ChatCompletionChoice)
	errChan := make(chan error)

	go l.generateContent(ctx, req, respChan, errChan)

	select {
	case resp := <-respChan:
		return resp.Message.Content, nil
	case err := <-errChan:
		return "", err
	}
}

type StreamingMessage struct {
	Choice *openai.ChatCompletionChoice
	Usage  *openai.Usage
	Err    error
}

func (l *LLM) GenerateContent(ctx context.Context, messages []openai.ChatCompletionMessage, streamingFunc func(msg StreamingMessage), options ...Option) {
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
		}
		req.Tools = llmTools(mapping)
	}

	respChan := make(chan *openai.ChatCompletionChoice)
	errChan := make(chan error)
	go l.generateContentStream(ctx, req, respChan, errChan)
	var choice *openai.ChatCompletionChoice
	sb := strings.Builder{}
	toolCalls := make([]openai.ToolCall, 0)
out:
	for {
		select {
		case delta := <-respChan:
			deltaToolCalls := delta.Message.ToolCalls
			if len(deltaToolCalls) > 0 {
				for _, tc := range deltaToolCalls {
					if len(toolCalls) <= *tc.Index {
						choice = delta
						toolCalls = append(toolCalls, tc)
					} else {
						toolCalls[*tc.Index].Function.Arguments += tc.Function.Arguments
					}
				}
			} else {
				if choice == nil {
					choice = delta
				} else {
					sb.WriteString(delta.Message.Content)
					choice.Message.Content = delta.Message.Content
				}
				if req.Stream {
					streamingFunc(StreamingMessage{Choice: choice})
				}
			}
		case err := <-errChan:
			if errors.Is(err, io.EOF) {
				choice.Message.Content = sb.String()

				streamingFunc(StreamingMessage{Err: err})
				break out
			}
			streamingFunc(StreamingMessage{Err: err})
			return
		}
	}

	if len(toolCalls) > 0 {
		choice.Message.ToolCalls = toolCalls
		req.Messages = append(req.Messages, choice.Message)
		for {
			toolMessages, err := l.invokeTools(ctx, toolCalls)
			if err != nil {
				streamingFunc(StreamingMessage{Err: err})
				return
			}
			req.Messages = append(req.Messages, toolMessages...)
			toolCalls = make([]openai.ToolCall, 0)
			go l.generateContentStream(ctx, req, respChan, errChan)
		outloop:
			for {
				select {
				case delta := <-respChan:
					deltaToolCalls := delta.Message.ToolCalls
					if len(deltaToolCalls) > 0 {
						for _, tc := range deltaToolCalls {
							if len(toolCalls) <= *tc.Index {
								toolCalls = append(toolCalls, tc)
							} else {
								toolCalls[*tc.Index].Function.Arguments += tc.Function.Arguments
							}
						}
					} else {
						if choice == nil {
							choice = delta
						} else {
							sb.WriteString(delta.Message.Content)
							choice.Message.Content = delta.Message.Content
						}
						if req.Stream {
							streamingFunc(StreamingMessage{Choice: choice})
						}
					}
				case err := <-errChan:
					if errors.Is(err, io.EOF) {
						choice.Message.Content = sb.String()
						streamingFunc(StreamingMessage{Err: err})
						break outloop
					}
					streamingFunc(StreamingMessage{Err: err})
					return
				}
			}
			if len(toolCalls) == 0 {
				break
			}
		}
	}
	if !req.Stream {
		streamingFunc(StreamingMessage{Choice: choice})
	}
}

func (l *LLM) generateContent(ctx context.Context, req openai.ChatCompletionRequest, choiceChan chan *openai.ChatCompletionChoice, errChan chan error) {
	start := time.Now()
	resp, err := l.client.CreateChatCompletion(ctx, req)
	if err != nil {

		errChan <- err
		return
	}
	logger.FromContext(ctx).Info("time for generated content",
		"duration", time.Since(start),
		"model", req.Model,
		"prompt_tokens", resp.Usage.PromptTokens,
		"completion_tokens", resp.Usage.CompletionTokens,
	)
	if len(resp.Choices) == 0 {
		errChan <- fmt.Errorf("no choices returned")

		return
	}
	choiceChan <- &resp.Choices[0]
}

func (l *LLM) generateContentStream(ctx context.Context, req openai.ChatCompletionRequest, choiceChan chan *openai.ChatCompletionChoice, errChan chan error) {
	stream, err := l.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		errChan <- err
		return
	}
	defer stream.Close()
	start := time.Now()
	for {
		response, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.FromContext(ctx).Info("time for generated content stream",
					"duration", time.Since(start),
					"model", req.Model,
				)
			}

			errChan <- err
			return
		}

		delta := response.Choices[0]
		choice := &openai.ChatCompletionChoice{
			Index: delta.Index,
			Message: openai.ChatCompletionMessage{
				Role:         delta.Delta.Role,
				Content:      delta.Delta.Content,
				FunctionCall: delta.Delta.FunctionCall,
				ToolCalls:    delta.Delta.ToolCalls,
			},
			FinishReason: delta.FinishReason,
		}
		choiceChan <- choice
	}
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
