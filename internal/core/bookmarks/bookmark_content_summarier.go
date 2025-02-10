package bookmarks

import (
	"context"
	"fmt"
	"io"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/config"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader"
	"recally/internal/pkg/webreader/processor"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sashabaranov/go-openai"
)

func (s *Service) SummarierContent(ctx context.Context, tx db.DBTX, bookmarkID, userID uuid.UUID) (*BookmarkContentDTO, error) {
	user, err := auth.LoadUser(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user: %w", err)
	}

	bookmarkContent, err := s.GetBookmarkContentByBookmarkID(ctx, tx, bookmarkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmark content by id '%s': %w", bookmarkID.String(), err)
	}

	switch bookmarkContent.Type {
	case ContentTypeBookmark:
		err = s.summarierArticleContent(ctx, bookmarkID, user, bookmarkContent)
	case ContentTypeImage:
		err = s.summarierImageContent(ctx, bookmarkID, user, bookmarkContent)
	default:
		return bookmarkContent, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to summarize content: %w", err)
	}

	return s.UpdateBookmarkContent(ctx, tx, bookmarkContent)
}

func (s *Service) ProcessSummaryTags(bookmarkID, userID uuid.UUID, summary string) (string, []string) {
	tags, summary := parseTagsFromSummary(summary)
	if len(tags) > 0 {
		// link tags in background
		s.saveContentTags(bookmarkID, userID, tags, []string{})
	}
	return summary, tags
}

func (s *Service) saveContentTags(bookmarkID, userID uuid.UUID, newTags, oldTags []string) {
	if len(newTags) > 0 {
		// link tags in background
		newUserCtx := auth.SetUserToContextByUserID(context.Background(), userID)
		go func() {
			if err := db.RunInTransaction(newUserCtx, db.DefaultPool.Pool, func(ctx context.Context, tx pgx.Tx) error {
				return s.linkContentTags(ctx, tx, oldTags, newTags, bookmarkID, userID)
			}); err != nil {
				logger.Default.Error("failed to link content tags", "err", err, "bookmark_id", bookmarkID)
			}
		}()
	}
}

func (s *Service) summarierArticleContent(ctx context.Context, bookmarkID uuid.UUID, user *auth.UserDTO, bookmarkContent *BookmarkContentDTO) error {
	content := &webreader.Content{
		Markwdown: bookmarkContent.Content,
	}

	summarier := processor.NewSummaryProcessor(s.llm, processor.WithSummaryOptionUser(user))

	if len(content.Markwdown) < 1000 {
		logger.FromContext(ctx).Info("content is too short to summarise")
		return nil
	}

	if err := summarier.Process(ctx, content); err != nil {
		logger.Default.Error("failed to generate summary", "err", err)
	} else {
		summary, tags := s.ProcessSummaryTags(bookmarkID, user.ID, content.Summary)
		bookmarkContent.Summary = summary
		if len(tags) > 0 {
			bookmarkContent.Tags = tags
		}
	}
	return nil
}

func (s *Service) summarierImageContent(ctx context.Context, bookmarkID uuid.UUID, user *auth.UserDTO, bookmarkContent *BookmarkContentDTO) error {
	// get image public url
	imgUrl, err := bookmarkContent.GetFilePublicURL(ctx)
	if err != nil {
		return fmt.Errorf("failed to get image public url: %w", err)
	}

	prompt := defaultDescribeImagePrompt
	language := "English"
	model := config.Settings.OpenAI.VisionModel
	if user.Settings.DescribeImageOptions.Model != "" {
		model = user.Settings.DescribeImageOptions.Model
	}
	if user.Settings.DescribeImageOptions.Language != "" {
		language = user.Settings.DescribeImageOptions.Language
	}
	if user.Settings.DescribeImageOptions.Prompt != "" {
		prompt = user.Settings.DescribeImageOptions.Prompt
	}

	logger.FromContext(ctx).Info("start describe image", "model", model, "language", language)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: prompt,
		},
		{
			Role: openai.ChatMessageRoleUser,
			MultiContent: []openai.ChatMessagePart{
				{
					Type: openai.ChatMessagePartTypeText,
					Text: fmt.Sprintf("Describe the image in %s", language),
				},
				{
					Type: openai.ChatMessagePartTypeImageURL,
					ImageURL: &openai.ChatMessageImageURL{
						URL: imgUrl,
					},
				},
			},
		},
	}

	streamingFunc := func(content llms.StreamingMessage) {
		if content.Err != nil && content.Err != io.EOF {
			logger.FromContext(ctx).Error("failed to generate describe image content", "err", content.Err)
			return
		}

		if content.Choice != nil {
			text := content.Choice.Message.Content
			bookmarkContent.Title = parseXmlContent(text, "title")
			bookmarkContent.Summary = parseXmlContent(text, "description")
			tagString := parseXmlContent(text, "tags")
			tags := []string{}
			for _, tag := range strings.Split(tagString, ", ") {
				if tag != "" {
					tags = append(tags, strings.TrimSpace(tag))
				}
			}
			bookmarkContent.Tags = tags
		}
	}

	s.llm.GenerateContent(ctx, messages, streamingFunc, llms.WithModel(model), llms.WithStream(false))

	if len(bookmarkContent.Tags) > 0 {
		s.saveContentTags(bookmarkID, user.ID, bookmarkContent.Tags, []string{})
	}

	return nil
}
