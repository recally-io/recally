package assistants

import (
	"context"
	"fmt"
	"time"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/llms"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sashabaranov/go-openai"
)

type Service struct {
	llm *llms.LLM
	dao dao
}

func NewService(llm *llms.LLM) *Service {
	return &Service{
		llm: llm,
		dao: db.New(),
	}
}

func (s *Service) ListAssistants(ctx context.Context, tx db.DBTX, userId uuid.UUID) ([]AssistantDTO, error) {
	asts, err := s.dao.ListAssistantsByUser(ctx, tx, pgtype.UUID{Bytes: userId, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get assistants: %w", err)
	}
	asstants := make([]AssistantDTO, 0, len(asts))
	for _, ast := range asts {
		var a AssistantDTO
		a.Load(&ast)
		asstants = append(asstants, a)
	}
	return asstants, nil
}

func (s *Service) CreateAssistant(ctx context.Context, tx db.DBTX, assistant *AssistantDTO) (*AssistantDTO, error) {
	model := assistant.Dump()
	ast, err := s.dao.CreateAssistant(ctx, tx, db.CreateAssistantParams{
		UserID:       model.UserID,
		Name:         model.Name,
		Description:  model.Description,
		SystemPrompt: model.SystemPrompt,
		Model:        model.Model,
		Metadata:     model.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create assistant: %w", err)
	}
	assistant.Load(&ast)
	return assistant, nil
}

func (s *Service) UpdateAssistant(ctx context.Context, tx db.DBTX, assistant *AssistantDTO) (*AssistantDTO, error) {
	model := assistant.Dump()

	ast, err := s.dao.UpdateAssistant(ctx, tx, db.UpdateAssistantParams{
		Uuid:         assistant.Id,
		Name:         model.Name,
		Description:  model.Description,
		SystemPrompt: model.SystemPrompt,
		Model:        model.Model,
		Metadata:     model.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update assistant: %w", err)
	}
	assistant.Load(&ast)
	return assistant, nil
}

func (s *Service) GetAssistant(ctx context.Context, tx db.DBTX, id uuid.UUID) (*AssistantDTO, error) {
	ast, err := s.dao.GetAssistant(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get assistant: %w", err)
	}
	var assistant AssistantDTO
	assistant.Load(&ast)
	return &assistant, nil
}

func (s *Service) ListThreads(ctx context.Context, tx db.DBTX, assistantID uuid.UUID) ([]ThreadDTO, error) {
	threads, err := s.dao.ListAssistantThreads(ctx, tx, pgtype.UUID{Bytes: assistantID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get threads: %w", err)
	}

	var result []ThreadDTO
	for _, th := range threads {
		var t ThreadDTO
		t.Load(&th)
		result = append(result, t)
	}

	return result, nil
}

func (s *Service) CreateThread(ctx context.Context, tx db.DBTX, thread *ThreadDTO) (*ThreadDTO, error) {
	model := thread.Dump()
	th, err := s.dao.CreateAssistantThread(ctx, tx, db.CreateAssistantThreadParams{
		UserID:      model.UserID,
		AssistantID: model.AssistantID,
		Name:        model.Name,
		Description: model.Description,
		Model:       model.Model,
		Metadata:    model.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create thread: %w", err)
	}
	thread.Load(&th)
	return thread, nil
}

func (s *Service) GetThread(ctx context.Context, tx db.DBTX, id uuid.UUID) (*ThreadDTO, error) {
	th, err := s.dao.GetAssistantThread(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread: %w", err)
	}
	var t ThreadDTO
	t.Load(&th)

	messages, err := s.ListThreadMessages(ctx, tx, th.Uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread messages: %w", err)
	}
	t.Messages = messages
	return &t, nil
}

func (s *Service) ListThreadMessages(ctx context.Context, tx db.DBTX, threadID uuid.UUID) ([]ThreadMessageDTO, error) {
	messages, err := s.dao.ListThreadMessages(ctx, tx, pgtype.UUID{Bytes: threadID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get thread messages: %w", err)
	}

	var result []ThreadMessageDTO
	for _, msg := range messages {
		var m ThreadMessageDTO
		m.Load(&msg)

		result = append(result, m)
	}

	return result, nil
}

func (s *Service) CreateThreadMessage(ctx context.Context, tx db.DBTX, threadId uuid.UUID, message *ThreadMessageDTO) (*ThreadMessageDTO, error) {
	model := message.Dump()
	tm, err := s.dao.CreateThreadMessage(ctx, tx, db.CreateThreadMessageParams{
		UserID:      model.UserID,
		ThreadID:    model.ThreadID,
		Model:       model.Model,
		Role:        model.Role,
		Text:        model.Text,
		Attachments: model.Attachments,
		Metadata:    model.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save thread message: %w", err)
	}
	message.Load(&tm)
	return message, nil
}

func (s *Service) AddThreadMessage(ctx context.Context, tx db.DBTX, thread *ThreadDTO, role, text string) (*ThreadMessageDTO, error) {
	thread.AddMessage(role, text)
	message := &ThreadMessageDTO{
		UserID:   thread.UserId,
		ThreadID: thread.Id,
		Model:    thread.Model,
		Role:     role,
		Text:     text,
	}
	return s.CreateThreadMessage(ctx, tx, thread.Id, message)
}

func (s *Service) RunThread(ctx context.Context, tx db.DBTX, thread *ThreadDTO) (*ThreadMessageDTO, error) {
	oaiMessages := make([]openai.ChatCompletionMessage, 0)
	messages, err := s.ListThreadMessages(ctx, tx, thread.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread messages: %w", err)
	}
	for _, m := range messages {
		oaiMessages = append(oaiMessages, openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Text,
		})
	}
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

	message := &ThreadMessageDTO{
		UserID:   thread.UserId,
		ThreadID: thread.Id,
		Model:    thread.Model,
		Role:     resp.Message.Role,
		Text:     resp.Message.Content,
		Token:    usage.TotalTokens,
	}

	message, err = s.CreateThreadMessage(ctx, tx, thread.Id, message)
	if err != nil {
		return nil, fmt.Errorf("failed to save thread message: %w", err)
	}
	return message, err
}

func (s *Service) ListModels(ctx context.Context) ([]string, error) {
	cacheKey := cache.NewCacheKey("list-models", "")
	if models, ok := cache.Get[[]string](ctx, cache.MemCache, cacheKey); ok {
		return *models, nil
	}
	models, err := s.llm.ListModels(ctx)
	if err != nil {
		return nil, fmt.Errorf("list models error: %w", err)
	}
	cache.MemCache.Set(cacheKey, &models, time.Hour)
	return models, nil
}
