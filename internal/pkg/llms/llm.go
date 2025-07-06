package llms

import (
	"context"
	"errors"
	"fmt"
	"io"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/config"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/tools"
	"strings"
	"sync"
	"time"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/sync/errgroup"
)

var DefaultLLM *LLM

func init() {
	DefaultLLM = New(config.Settings.OpenAI.BaseURL, config.Settings.OpenAI.ApiKey)
}

const (
	IntermediateStepTool = "tool"
	IntermediateStepRag  = "rag"
)

type IntermediateStep struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Input  any    `json:"input"`
	Output any    `json:"output"`
}

type StreamingMessage struct {
	Choice            *openai.ChatCompletionChoice
	Usage             *openai.Usage
	Err               error
	IntermediateSteps []IntermediateStep `json:"intermediate_steps"`
}

func (m *StreamingMessage) ToStreamingString() StreamingString {
	if m.Choice != nil {
		return StreamingString{Content: m.Choice.Message.Content, Err: m.Err}
	} else {
		return StreamingString{Err: m.Err}
	}
}

type StreamingString struct {
	Content string `json:"content"`
	Err     error  `json:"err"`
}

type LLM struct {
	client       *openai.Client
	toolMappings map[string]tools.Tool
}

type Model struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

func (l *LLM) ListModels(ctx context.Context) ([]Model, error) {
	models, err := cache.RunInCache[[]Model](ctx, cache.MemCache, cache.NewCacheKey("llm", "list-models"), time.Hour, func() (*[]Model, error) {
		models, err := l.client.ListModels(ctx)
		if err != nil {
			return nil, err
		}

		data := make([]Model, 0, len(models.Models))
		for _, m := range models.Models {
			data = append(data, Model{
				ID:   m.ID,
				Name: m.ID,
			})
		}

		return &data, nil
	})

	return *models, err
}

func (l *LLM) ListTools(ctx context.Context) ([]tools.BaseTool, error) {
	availableTools := make([]tools.BaseTool, 0, len(AllToolMappings))
	for _, tool := range AllToolMappings {
		availableTools = append(availableTools, tools.BaseTool{
			Name:        tool.LLMName(),
			Description: tool.LLMDescription(),
		})
	}

	return availableTools, nil
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

func (l *LLM) StreamingTextCompletion(ctx context.Context, prompt string, streamingFunc func(content StreamingString), options ...Option) {
	opts := &Options{}
	for _, o := range options {
		o(opts)
	}

	options = append(options, WithStream(true))

	req := opts.ToChatCompletionRequest()
	req.Messages = []openai.ChatCompletionMessage{{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	}}

	sendToUser := func(m StreamingMessage) {
		streamingFunc(m.ToStreamingString())
	}

	l.GenerateContent(ctx, req.Messages, sendToUser, options...)
}

func (l *LLM) GenerateContent(ctx context.Context, messages []openai.ChatCompletionMessage, streamingFunc func(msg StreamingMessage), options ...Option) {
	opts := &Options{}
	for _, o := range options {
		o(opts)
	}

	req := opts.ToChatCompletionRequest()
	req.Messages = messages

	mu := new(sync.Mutex)

	syncStreamFunc := func(msg StreamingMessage) {
		mu.Lock()
		defer mu.Unlock()
		streamingFunc(msg)
	}

	// dynamically add tools to the request
	logger.FromContext(ctx).Info("generating content", "tool_names", opts.ToolNames, "model", req.Model)

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

	// o1- models don't support system messages or tools
	// TODO: refactor this to be more elegant
	if strings.HasPrefix(req.Model, "o1-") {
		if req.Messages[0].Role == openai.ChatMessageRoleSystem {
			req.Messages[0].Role = openai.ChatMessageRoleUser
		}

		if len(req.Tools) > 0 {
			req.Tools = nil
		}

		req.Stream = false
	}

	respChan := make(chan *openai.ChatCompletionChoice)
	errChan := make(chan error)

	go l.generateContent(ctx, req, respChan, errChan)

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
					syncStreamFunc(StreamingMessage{Choice: choice})
				}
			}
		case err := <-errChan:
			if errors.Is(err, io.EOF) {
				if req.Stream {
					choice.Message.Content = sb.String()
				}

				break out
			}
			syncStreamFunc(StreamingMessage{Err: err})

			return
		}
	}

	if len(toolCalls) > 0 {
		for {
			req.Messages = append(req.Messages, openai.ChatCompletionMessage{
				Role:      openai.ChatMessageRoleAssistant,
				ToolCalls: toolCalls,
			})

			toolMessages, steps, err := l.invokeTools(ctx, toolCalls)
			if err != nil {
				syncStreamFunc(StreamingMessage{Err: err})

				return
			}

			syncStreamFunc(StreamingMessage{IntermediateSteps: steps})

			req.Messages = append(req.Messages, toolMessages...)
			toolCalls = make([]openai.ToolCall, 0)

			go l.generateContent(ctx, req, respChan, errChan)
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
							syncStreamFunc(StreamingMessage{Choice: choice})
						}
					}
				case err := <-errChan:
					if errors.Is(err, io.EOF) {
						if req.Stream {
							choice.Message.Content = sb.String()
						}

						break outloop
					}
					syncStreamFunc(StreamingMessage{Err: err})

					return
				}
			}

			if len(toolCalls) == 0 {
				break
			}
		}
	}

	if !req.Stream {
		syncStreamFunc(StreamingMessage{Choice: choice})
	}

	syncStreamFunc(StreamingMessage{Err: io.EOF})
}

func (l *LLM) generateContent(ctx context.Context, req openai.ChatCompletionRequest, choiceChan chan *openai.ChatCompletionChoice, errChan chan error) {
	if req.Stream {
		l.generateContentStream(ctx, req, choiceChan, errChan)

		return
	}

	start := time.Now()

	resp, err := l.client.CreateChatCompletion(ctx, req)
	if err != nil {
		errChan <- err

		return
	}

	logger.FromContext(ctx).Info("time for generated content",
		"duration", time.Since(start).Milliseconds(),
		"model", req.Model,
		"prompt_tokens", resp.Usage.PromptTokens,
		"completion_tokens", resp.Usage.CompletionTokens,
	)

	if len(resp.Choices) == 0 {
		errChan <- fmt.Errorf("no choices returned")

		return
	}
	choiceChan <- &resp.Choices[0]
	errChan <- io.EOF
}

func (l *LLM) generateContentStream(ctx context.Context, req openai.ChatCompletionRequest, choiceChan chan *openai.ChatCompletionChoice, errChan chan error) {
	// get token usage data for streamed chat completion response
	req.StreamOptions = &openai.StreamOptions{
		IncludeUsage: true,
	}

	stream, err := l.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		errChan <- err

		return
	}

	defer stream.Close()

	start := time.Now()
	usage := &openai.Usage{}

	for {
		response, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.FromContext(ctx).Info("time for generated content stream",
					"duration", time.Since(start).Milliseconds(),
					"model", req.Model,
					"prompt_tokens", usage.PromptTokens,
					"completion_tokens", usage.CompletionTokens,
				)
			}

			errChan <- err

			return
		}

		if response.Usage != nil {
			usage = response.Usage
		}

		if len(response.Choices) > 0 {
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
}

func (l *LLM) invokeTools(ctx context.Context, toolCalls []openai.ToolCall) ([]openai.ChatCompletionMessage, []IntermediateStep, error) {
	var messages []openai.ChatCompletionMessage

	eg, ctx := errgroup.WithContext(ctx)
	steps := make([]IntermediateStep, 0)

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

			logger.FromContext(ctx).Info("tool response", "tool", toolName, "response", toolResp[:min(200, len(toolResp))], "duration", time.Since(start).Milliseconds())
			messages = append(messages, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				ToolCallID: tc.ID,
				Content:    toolResp,
			})

			steps = append(steps, IntermediateStep{
				Type:   IntermediateStepTool,
				Name:   toolName,
				Input:  toolArgs,
				Output: toolResp,
			})

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, steps, err
	}

	return messages, steps, nil
}
