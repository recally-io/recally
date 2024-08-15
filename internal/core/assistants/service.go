package assistants

import (
	"context"
	"fmt"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/llms"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

type Service struct {
	llm *llms.LLM
	r   Repository
}

func NewService(llm *llms.LLM) *Service {
	return &Service{
		llm: llm,
		r:   NewRepository(),
	}
}

func (s *Service) ListAssistants(ctx context.Context, tx db.DBTX, userId uuid.UUID) ([]AssistantDTO, error) {
	return s.r.ListAssistants(ctx, tx, userId)
}

func (s *Service) CreateAssistant(ctx context.Context, tx db.DBTX, assistant *AssistantDTO) error {
	return s.r.CreateAssistant(ctx, tx, assistant)
}

func (s *Service) GetAssistant(ctx context.Context, tx db.DBTX, id uuid.UUID) (*AssistantDTO, error) {
	return s.r.GetAssistant(ctx, tx, id)
}

func (s *Service) ListThreads(ctx context.Context, tx db.DBTX, assistantID uuid.UUID) ([]ThreadDTO, error) {
	return s.r.ListThreads(ctx, tx, assistantID)
}

func (s *Service) CreateThread(ctx context.Context, tx db.DBTX, thread *ThreadDTO) error {
	return s.r.CreateThread(ctx, tx, thread)
}

func (s *Service) GetThread(ctx context.Context, tx db.DBTX, id uuid.UUID) (*ThreadDTO, error) {
	return s.r.GetThread(ctx, tx, id)
}

func (s *Service) ListThreadMessages(ctx context.Context, tx db.DBTX, threadID uuid.UUID) ([]ThreadMessageDTO, error) {
	return s.r.ListThreadMessages(ctx, tx, threadID)
}

func (s *Service) AddThreadMessage(ctx context.Context, tx db.DBTX, thread *ThreadDTO, role, text string) error {
	thread.AddMessage(role, text)
	message := ThreadMessageDTO{
		UserID:   thread.UserId,
		ThreadID: thread.Id,
		Model:    thread.Model,
		Role:     role,
		Text:     text,
	}
	return s.r.CreateThreadMessage(ctx, tx, thread.Id, message)
}

func (s *Service) RunThread(ctx context.Context, tx db.DBTX, thread *ThreadDTO) (*ThreadMessageDTO, error) {
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

	message := ThreadMessageDTO{
		UserID:   thread.UserId,
		ThreadID: thread.Id,
		Model:    thread.Model,
		Role:     resp.Message.Role,
		Text:     resp.Message.Content,
		Token:    usage.TotalTokens,
	}

	if err := s.r.CreateThreadMessage(ctx, tx, thread.Id, message); err != nil {
		return nil, fmt.Errorf("failed to save thread message: %w", err)
	}
	return &message, err
}

func (s *Service) GetTelegramUser(ctx context.Context, tx db.DBTX, userID string) (*User, error) {
	return s.r.GetTelegramUser(ctx, tx, userID)
}

func (s *Service) CreateTelegramUser(ctx context.Context, tx db.DBTX, userName string, userID string) (*User, error) {
	return s.r.CreateTelegramUser(ctx, tx, userName, userID)
}

func (s *Service) UpdateTelegramUser(ctx context.Context, tx db.DBTX, user User) (*User, error) {
	return s.r.UpdateTelegramUser(ctx, tx, user)
}
