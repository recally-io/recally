package assistants

import (
	"context"
	"fmt"
	"vibrain/internal/pkg/llms"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

type Service struct {
	llm *llms.LLM
}

func NewService(llm *llms.LLM) (*Service, error) {
	s := &Service{
		llm: llm,
	}

	return s, nil
}

func (s *Service) CreateAssistant(ctx context.Context, assistant *Assistant) error {
	r, err := RepositoryFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get repository from context: %w", err)
	}
	return r.CreateAssistant(ctx, assistant)
}

func (s *Service) GetAssistant(ctx context.Context, id uuid.UUID) (*Assistant, error) {
	r, err := RepositoryFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository from context: %w", err)
	}
	return r.GetAssistant(ctx, id)
}

func (s *Service) CreateThread(ctx context.Context, thread *Thread) error {
	r, err := RepositoryFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get repository from context: %w", err)
	}
	return r.CreateThread(ctx, thread)
}

func (s *Service) GetThread(ctx context.Context, id uuid.UUID) (*Thread, error) {
	r, err := RepositoryFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository from context: %w", err)
	}
	return r.GetThread(ctx, id)
}

func (s *Service) AddThreadMessage(ctx context.Context, thread *Thread, role, text string) error {
	thread.AddMessage(role, text)
	message := ThreadMessage{
		UserID:   thread.UserId,
		ThreadID: thread.Id,
		Model:    thread.Model,
		Role:     role,
		Text:     text,
	}
	r, err := RepositoryFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get repository from context: %w", err)
	}
	return r.CreateThreadMessage(ctx, thread.Id, message)
}

func (s *Service) RunThread(ctx context.Context, thread *Thread) (*ThreadMessage, error) {
	r, err := RepositoryFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository from context: %w", err)
	}
	oaiMessages := make([]openai.ChatCompletionMessage, 0)
	oaiMessages = append(oaiMessages, openai.ChatCompletionMessage{
		Role:    "system",
		Content: thread.SystemPrompt,
	})
	for _, m := range thread.Messages {
		oaiMessages = append(oaiMessages, openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Text,
		})
	}
	resp, usage, err := s.llm.GenerateContent(ctx, oaiMessages)
	if err != nil {
		return nil, err
	}

	message := ThreadMessage{
		UserID:   thread.UserId,
		ThreadID: thread.Id,
		Model:    thread.Model,
		Role:     resp.Message.Role,
		Text:     resp.Message.Content,
		Token:    usage.TotalTokens,
	}

	if err := r.CreateThreadMessage(ctx, thread.Id, message); err != nil {
		return nil, fmt.Errorf("failed to save thread message: %w", err)
	}
	return &message, err
}
