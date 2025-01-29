package bookmarks

import (
	"context"
	"fmt"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"slices"

	"github.com/google/uuid"
)

func (s *Service) linkContentTags(ctx context.Context, tx db.DBTX, originTags, newTags []string, contentID, userID uuid.UUID) error {
	// create tags if not exist
	allExistingTags, err := s.dao.ListExistingBookmarkTagsByTags(ctx, tx, db.ListExistingBookmarkTagsByTagsParams{
		Column1: newTags,
		UserID:  userID,
	})
	if err != nil {
		return fmt.Errorf("failed to list existing tags: %w", err)
	}
	for _, tag := range newTags {
		if slices.Contains(allExistingTags, tag) {
			continue
		}
		if _, err := s.dao.CreateBookmarkTag(ctx, tx, db.CreateBookmarkTagParams{
			Name:   tag,
			UserID: userID,
		}); err != nil {
			return fmt.Errorf("failed to create tag '%s': %w", tag, err)
		}
	}

	// link content with tags that not linked before
	contentExistingTags, err := s.dao.ListBookmarkTagsByBookmarkId(ctx, tx, contentID)
	if err != nil {
		return fmt.Errorf("failed to list content tags: %w", err)
	}
	newLinkedTags := make([]string, 0)
	for _, tag := range newTags {
		if !slices.Contains(contentExistingTags, tag) {
			newLinkedTags = append(newLinkedTags, tag)
		}
	}
	if err := s.dao.LinkBookmarkWithTags(ctx, tx, db.LinkBookmarkWithTagsParams{
		BookmarkID: contentID,
		Column2:    newLinkedTags,
		UserID:     userID,
	}); err != nil {
		return fmt.Errorf("failed to link tags with content: %w", err)
	}

	// unlink content with tags in original but not in new
	removedTags := make([]string, 0)
	for _, tag := range originTags {
		if !slices.Contains(newTags, tag) {
			removedTags = append(removedTags, tag)
		}
	}

	if err := s.dao.UnLinkBookmarkWithTags(ctx, tx, db.UnLinkBookmarkWithTagsParams{
		BookmarkID: contentID,
		Column2:    removedTags,
		UserID:     userID,
	}); err != nil {
		return fmt.Errorf("failed to unlink tags with content: %w", err)
	}

	logger.FromContext(ctx).Info("link content with tags", "content_id", contentID, "new_tags", newTags, "origin_tags", originTags, "removed_tags", removedTags, "new_linked_tags", newLinkedTags)
	return nil
}
