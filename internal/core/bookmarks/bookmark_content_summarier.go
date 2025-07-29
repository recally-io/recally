package bookmarks

import (
	"context"
	"fmt"
	"io"
	"recally/internal/core/files"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader"
	"recally/internal/pkg/webreader/processor"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (s *Service) SaveContentTags(bookmarkID, userID uuid.UUID, newTags []string) {
	if len(newTags) > 0 {
		// link tags in background
		newUserCtx := auth.SetUserToContextByUserID(context.Background(), userID)
		go func() {
			if err := db.RunInTransaction(newUserCtx, db.DefaultPool.Pool, func(ctx context.Context, tx pgx.Tx) error {
				return s.linkContentTags(ctx, tx, newTags, bookmarkID, userID)
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
		logger.FromContext(ctx).Info("content is too short to summarize")

		return nil
	}

	if err := summarier.Process(ctx, content); err != nil {
		logger.Default.Error("failed to generate summary", "err", err)
	} else {
		summary, tags := summarier.ParseSummaryInfo(content.Summary)
		bookmarkContent.Summary = summary

		if len(tags) > 0 {
			bookmarkContent.Tags = tags
			s.SaveContentTags(bookmarkID, user.ID, tags)
		}
	}

	return nil
}

func (s *Service) summarierImageContent(ctx context.Context, bookmarkID uuid.UUID, user *auth.UserDTO, bookmarkContent *BookmarkContentDTO) error {
	if bookmarkContent.S3Key == "" {
		return fmt.Errorf("s3 key is empty")
	}

	imgReader, err := files.DefaultService.LoadFileContentByS3Key(ctx, bookmarkContent.S3Key)
	if err != nil {
		return fmt.Errorf("failed to load image content: %w", err)
	}

	// summarize image
	summarier := processor.NewSummaryImageProcessor(s.llm, processor.WithSummaryImageOptionUser(user))

	imgDataUrl, err := summarier.EncodeImage(imgReader, bookmarkContent.Content)
	if err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	streamingFunc := func(content llms.StreamingMessage) {
		if content.Err != nil && content.Err != io.EOF {
			logger.FromContext(ctx).Error("failed to generate describe image content", "err", content.Err)

			return
		}

		if content.Choice != nil {
			text := content.Choice.Message.Content
			title, description, tags := summarier.ParseSummaryInfo(text)
			bookmarkContent.Title = title
			bookmarkContent.Summary = description
			bookmarkContent.Tags = tags
		}
	}

	summarier.Summary(ctx, imgDataUrl, streamingFunc)

	if len(bookmarkContent.Tags) > 0 {
		s.SaveContentTags(bookmarkID, user.ID, bookmarkContent.Tags)
	}

	return nil
}
