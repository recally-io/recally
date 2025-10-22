package bookmarks

import (
	"context"
	"fmt"
	"slices"

	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"

	"github.com/google/uuid"
)

func (s *Service) linkContentTags(
	ctx context.Context,
	tx db.DBTX,
	newTags []string,
	contentID, userID uuid.UUID,
) error {
	// 1) Fetch all tags currently linked to this bookmark
	currentLinkedTags, err := s.dao.ListBookmarkTagsByBookmarkId(ctx, tx, contentID)
	if err != nil {
		return fmt.Errorf("failed to list existing linked tags for content: %w", err)
	}

	// 2) Figure out which tags need adding/removing
	toAdd := difference(newTags, currentLinkedTags)
	toRemove := difference(currentLinkedTags, newTags)

	// 3) Ensure any “toAdd” tags actually exist in the user’s DB
	if err := s.ensureTagsExist(ctx, tx, toAdd, userID); err != nil {
		return fmt.Errorf("failed to ensure tags exist: %w", err)
	}

	// 4) Link new tags
	if len(toAdd) > 0 {
		if err := s.dao.LinkBookmarkWithTags(ctx, tx, db.LinkBookmarkWithTagsParams{
			BookmarkID: contentID,
			Column2:    toAdd,
			UserID:     userID,
		}); err != nil {
			return fmt.Errorf("failed to link new tags: %w", err)
		}
	}

	// 5) Unlink removed tags
	if len(toRemove) > 0 {
		if err := s.dao.UnLinkBookmarkWithTags(ctx, tx, db.UnLinkBookmarkWithTagsParams{
			BookmarkID: contentID,
			Column2:    toRemove,
			UserID:     userID,
		}); err != nil {
			return fmt.Errorf("failed to unlink removed tags: %w", err)
		}
	}

	// 6) Log result
	logger.FromContext(ctx).Info("link content with tags",
		"content_id", contentID,
		"new_tags", newTags,
		"added_tags", toAdd,
		"removed_tags", toRemove,
	)

	return nil
}

// difference returns elements in sliceA that are not in sliceB.
func difference(sliceA, sliceB []string) []string {
	setB := make(map[string]struct{}, len(sliceB))
	for _, val := range sliceB {
		setB[val] = struct{}{}
	}

	var diff []string

	for _, val := range sliceA {
		if _, found := setB[val]; !found {
			diff = append(diff, val)
		}
	}

	return diff
}

// ensureTagsExist creates any tags in targetTags that do not currently exist in the DB.
func (s *Service) ensureTagsExist(ctx context.Context, tx db.DBTX, targetTags []string, userID uuid.UUID) error {
	if len(targetTags) == 0 {
		return nil
	}

	existing, err := s.dao.ListExistingBookmarkTagsByTags(ctx, tx, db.ListExistingBookmarkTagsByTagsParams{
		Column1: targetTags,
		UserID:  userID,
	})
	if err != nil {
		return fmt.Errorf("failed to list existing tags: %w", err)
	}

	for _, tag := range targetTags {
		if slices.Contains(existing, tag) {
			continue
		}

		if _, err := s.dao.CreateBookmarkTag(ctx, tx, db.CreateBookmarkTagParams{
			Name:   tag,
			UserID: userID,
		}); err != nil {
			return fmt.Errorf("failed to create tag '%s': %w", tag, err)
		}
	}

	return nil
}
