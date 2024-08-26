package assistants

import (
	"context"
	"fmt"
	"strings"
	"time"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/llms"
	"vibrain/internal/pkg/tools"

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
	ass, err := s.GetAssistant(ctx, tx, thread.AssistantId)
	if err != nil {
		return nil, fmt.Errorf("failed to get assistant: %w", err)
	}
	if thread.Model == "" {
		thread.Model = ass.Model
	}
	if thread.Metadata.Tools == nil {
		thread.Metadata.Tools = ass.Metadata.Tools
	}
	if thread.SystemPrompt == "" {
		thread.SystemPrompt = ass.SystemPrompt
	}

	model := thread.Dump()
	th, err := s.dao.CreateAssistantThread(ctx, tx, db.CreateAssistantThreadParams{
		Uuid:         model.Uuid,
		UserID:       model.UserID,
		AssistantID:  model.AssistantID,
		Name:         model.Name,
		Description:  model.Description,
		SystemPrompt: model.SystemPrompt,
		Model:        model.Model,
		Metadata:     model.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create thread: %w", err)
	}
	thread.Load(&th)
	return thread, nil
}

func (s *Service) UpdateThread(ctx context.Context, tx db.DBTX, thread *ThreadDTO) (*ThreadDTO, error) {
	model := thread.Dump()
	th, err := s.dao.UpdateAssistantThread(ctx, tx, db.UpdateAssistantThreadParams{
		Uuid:         thread.Id,
		Name:         model.Name,
		Description:  model.Description,
		Model:        model.Model,
		Metadata:     model.Metadata,
		SystemPrompt: model.SystemPrompt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update thread: %w", err)
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

	// messages, err := s.ListThreadMessages(ctx, tx, th.Uuid)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get thread messages: %w", err)
	// }
	// t.Messages = messages
	return &t, nil
}

func (s *Service) DeleteThread(ctx context.Context, tx db.DBTX, id uuid.UUID) error {
	if err := s.dao.DeleteThreadMessagesByThread(ctx, tx, pgtype.UUID{Bytes: id, Valid: true}); err != nil {
		return fmt.Errorf("failed to delete thread messages: %w", err)
	}
	if err := s.dao.DeleteAssistantThread(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete thread: %w", err)
	}
	return nil
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

func (s *Service) GetThreadMessage(ctx context.Context, tx db.DBTX, id uuid.UUID) (*ThreadMessageDTO, error) {
	msg, err := s.dao.GetThreadMessage(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread message: %w", err)
	}
	var m ThreadMessageDTO
	m.Load(&msg)
	return &m, nil
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

func (s *Service) DeleteThreadMessage(ctx context.Context, tx db.DBTX, id uuid.UUID) error {
	msg, err := s.dao.GetThreadMessage(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("failed to get thread message: %w", err)
	}
	if err := s.dao.DeleteThreadMessageByThreadAndCreatedAt(ctx, tx, db.DeleteThreadMessageByThreadAndCreatedAtParams{
		ThreadID:  msg.ThreadID,
		CreatedAt: msg.CreatedAt,
	}); err != nil {
		return fmt.Errorf("failed to delete thread message: %w", err)
	}
	return nil
}

func (s *Service) RunThread(ctx context.Context, tx db.DBTX, id uuid.UUID) (*ThreadMessageDTO, error) {
	thread, err := s.GetThread(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread: %w", err)
	}
	oaiMessages := make([]openai.ChatCompletionMessage, 0)
	messages, err := s.ListThreadMessages(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread messages: %w", err)
	}
	oaiMessages = append(oaiMessages, openai.ChatCompletionMessage{
		Role:    "system",
		Content: thread.SystemPrompt,
	})
	model := thread.Model
	toolNames := thread.Metadata.Tools
	for _, m := range messages {
		oaiMessages = append(oaiMessages, openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Text,
		})
	}
	lastMessage := messages[len(messages)-1]
	// Use the model from the last message
	if lastMessage.Model != "" {
		model = lastMessage.Model
	}
	if lastMessage.Metadata.Tools != nil {
		toolNames = lastMessage.Metadata.Tools
	}

	if len(lastMessage.Metadata.Images) > 0 {
		multiContent := make([]openai.ChatMessagePart, 0)

		if lastMessage.Text != "" {
			multiContent = append(multiContent, openai.ChatMessagePart{
				Type: openai.ChatMessagePartTypeText,
				Text: lastMessage.Text,
			})
		}
		for _, img := range lastMessage.Metadata.Images {
			multiContent = append(multiContent, openai.ChatMessagePart{
				Type: openai.ChatMessagePartTypeImageURL,
				ImageURL: &openai.ChatMessageImageURL{
					URL: img,
				},
			})
		}
		oaiMessages[len(oaiMessages)-1] = openai.ChatCompletionMessage{
			Role:         "user",
			MultiContent: multiContent,
		}
	}

	fmt.Println(oaiMessages)

	opts := []llms.Option{
		llms.WithModel(model),
		llms.WithToolNames(toolNames),
	}

	resp, usage, err := s.llm.GenerateContent(ctx, oaiMessages, opts...)
	if err != nil {
		return nil, err
	}

	message := &ThreadMessageDTO{
		UserID:   thread.UserId,
		ThreadID: thread.Id,
		Model:    model,
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

// GenerateThreadTitle generates a title for the thread based on the conversation.
// It uses the last 4 messages in the thread to generate a title.
// It uses the LLM to generate the title.
// It updates the thread with the generated title.
func (s *Service) GenerateThreadTitle(ctx context.Context, tx db.DBTX, id uuid.UUID) (string, error) {
	thread, err := s.GetThread(ctx, tx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get thread: %w", err)
	}
	messages, err := s.ListThreadMessages(ctx, tx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get thread messages: %w", err)
	}
	if len(messages) < 4 {
		return "", fmt.Errorf("not enough messages to generate title")
	}

	conversationStr := strings.Builder{}
	for _, m := range messages[:4] {
		conversationStr.WriteString(fmt.Sprintf("%s: %s\n", m.Role, m.Text))
		conversationStr.WriteString("\n")
	}

	prompt, err := getTitleGenerationPrompt(conversationStr.String())
	if err != nil {
		return "", fmt.Errorf("failed to get title generation prompt: %w", err)
	}

	title, err := s.llm.TextCompletion(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("GenerateThreadTitle: %w", err)
	}
	thread.Name = title
	thread.Metadata.IsGeneratedTitle = true

	_, err = s.UpdateThread(ctx, tx, thread)
	if err != nil {
		return "", fmt.Errorf("failed to update thread: %w", err)
	}

	return title, nil
}

func (s *Service) DeleteAssistant(ctx context.Context, tx db.DBTX, assistantId uuid.UUID) error {
	// Delete associated threads and messages
	if err := s.dao.DeleteThreadMessagesByAssistant(ctx, tx, pgtype.UUID{Bytes: assistantId, Valid: true}); err != nil {
		return fmt.Errorf("failed to delete thread messages by assistant: %w", err)
	}

	if err := s.dao.DeleteAssistantThreadsByAssistant(ctx, tx, pgtype.UUID{Bytes: assistantId, Valid: true}); err != nil {
		return fmt.Errorf("failed to delete assistant threads: %w", err)
	}

	// Delete the assistant
	if err := s.dao.DeleteAssistant(ctx, tx, assistantId); err != nil {
		return fmt.Errorf("failed to delete assistant: %w", err)
	}

	return nil
}

func (s *Service) ListTools(ctx context.Context) ([]tools.BaseTool, error) {
	toolMappings := llms.AllToolMappings
	availableTools := make([]tools.BaseTool, 0, len(toolMappings))
	for _, tool := range toolMappings {
		availableTools = append(availableTools, tools.BaseTool{
			Name:        tool.LLMName(),
			Description: tool.LLMDescription(),
		})
	}
	return availableTools, nil
}
