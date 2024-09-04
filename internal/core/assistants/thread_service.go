package assistants

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/llms"
	"vibrain/internal/pkg/logger"
	"vibrain/internal/pkg/rag/document"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
	"github.com/sashabaranov/go-openai"
)

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

func (s *Service) RunThread(ctx context.Context, tx db.DBTX, id uuid.UUID) (*MessageDTO, error) {
	thread, err := s.GetThread(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread: %w", err)
	}

	oaiMessages, model, toolNames, err := s.buildChatMessages(ctx, tx, thread)
	if err != nil {
		return nil, fmt.Errorf("failed to build chat messages: %w", err)
	}

	opts := []llms.Option{
		llms.WithModel(model),
		llms.WithToolNames(toolNames),
	}

	resp, usage, err := s.llm.GenerateContent(ctx, oaiMessages, opts...)
	if err != nil {
		return nil, err
	}

	message := &MessageDTO{
		UserID:          thread.UserId,
		AssistantID:     thread.AssistantId,
		ThreadID:        thread.Id,
		Model:           model,
		Role:            resp.Message.Role,
		Text:            resp.Message.Content,
		PromptToken:     int32(usage.PromptTokens),
		CompletionToken: int32(usage.CompletionTokens),
	}

	message, err = s.CreateThreadMessage(ctx, tx, thread.Id, message)
	if err != nil {
		return nil, fmt.Errorf("failed to save thread message: %w", err)
	}
	return message, err
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

func (s *Service) buildChatMessages(ctx context.Context, tx db.DBTX, thread *ThreadDTO) ([]openai.ChatCompletionMessage, string, []string, error) {
	oaiMessages := make([]openai.ChatCompletionMessage, 0)
	messages, err := s.ListThreadMessages(ctx, tx, thread.Id)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to get thread messages: %w", err)
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
	s.rewriteUserMessage(ctx, tx, &lastMessage)

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
	return oaiMessages, model, toolNames, nil
}

func (s *Service) rewriteUserMessage(ctx context.Context, tx db.DBTX, message *MessageDTO) {
	if message.Text == "" {
		logger.FromContext(ctx).Info("RAG for user message: message text is empty")
		return
	}
	// 1. search for similar documents
	// 2. search chat history for similar questions
	// 3. construct a new message
	var err error
	docs := make([]document.Document, 0)
	// chatHistory := make([]string, 0)

	embeddings, err := s.llm.CreateEmbeddings(ctx, message.Text)
	if err != nil {
		logger.FromContext(ctx).Error("failed to create embeddings", "err", err)
	}

	docsRes, err := s.dao.SimilaritySearchByThreadId(ctx, tx, db.SimilaritySearchByThreadIdParams{
		Uuid:       message.ThreadID,
		Embeddings: pgvector.NewVector(embeddings),
		Limit:      10,
	})
	if err != nil {
		logger.FromContext(ctx).Error("failed to search for similar documents", "err", err)
	}
	for _, d := range docsRes {
		var metadata map[string]any
		if err := json.Unmarshal(d.Metadata, &metadata); err != nil {
			logger.FromContext(ctx).Error("failed to unmarshal metadata", "err", err)
		}
		docs = append(docs, document.Document{
			Content:  d.Text,
			Metadata: metadata,
		})
	}

	// chatHistoryRes, err := s.dao.SimilaritySearchMessages(ctx, tx, db.SimilaritySearchMessagesParams{
	// 	ThreadID:   pgtype.UUID{Bytes: message.ThreadID, Valid: true},
	// 	Embeddings: pgvector.NewVector(embeddings),
	// 	Limit:      10,
	// })
	// if err != nil {
	// 	logger.FromContext(ctx).Error("failed to search for similar messages", "err", err)
	// }
	// for _, m := range chatHistoryRes {
	// 	chatHistory = append(chatHistory, m.Text.String)
	// }

	// chatHistoryStr := strings.Join(chatHistory, "\n")
	docStrBuilder := strings.Builder{}
	for _, d := range docs {
		docStrBuilder.WriteString(fmt.Sprintf("Content: %s\nMetadata: %v\nScore: %f\n\n", d.Content, d.Metadata, d.Score))
	}
	docsStr := docStrBuilder.String()
	prompt, err := getChatMessageWithRagPrompt(docsStr, "", message.Text)
	if err != nil {
		logger.FromContext(ctx).Error("failed to get chat message with RAG prompt", "err", err)
	}

	message.Text = prompt
}
