// Package bookmarks provides functionality for managing user bookmarks and their content.
package bookmarks

import (
	"context"
	"fmt"
	"net/url"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateBookmark creates a new bookmark for a user with the given content.
// It first validates the URL and checks if content already exists for the URL.
// If content exists and belongs to the user, returns ErrDuplicate.
// If content doesn't exist, creates new content.
// Finally creates the bookmark linking the user and content.
//
// Parameters:
//   - ctx: Context for the operation
//   - tx: Database transaction
//   - userId: UUID of the user creating the bookmark
//   - dto: BookmarkContentDTO containing bookmark details
//
// Returns:
//   - *BookmarkDTO: Created bookmark data
//   - error: ErrInvalidInput for invalid URL
//     ErrDuplicate if bookmark already exists
//     Other errors for database operations
func (s *Service) CreateBookmark(ctx context.Context, tx db.DBTX, userId uuid.UUID, dto *BookmarkContentDTO) (*BookmarkDTO, error) {
	// Validate URL format before proceeding
	u, err := url.Parse(dto.URL)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid URL", ErrInvalidInput)
	}

	dto.Domain = u.Host
	contentDTO := &BookmarkContentDTO{}

	// Create bookmark content for PDF and EPUB types
	if dto.Type == ContentTypePDF || dto.Type == ContentTypeEPUB || dto.Type == ContentTypeImage {
		contentDTO, err = s.CreateBookmarkContent(ctx, tx, dto)
		if err != nil {
			return nil, fmt.Errorf("failed to create bookmark content: %w", err)
		}
	} else {
		// Create bookmark content for other types
		// Check if bookmark content already exists for this URL and user
		content, err := s.dao.GetBookmarkContentByURL(ctx, tx, db.GetBookmarkContentByURLParams{
			Url:    dto.URL,
			UserID: pgtype.UUID{Bytes: userId, Valid: true},
		})
		if err == nil {
			contentDTO.Load(&content)
		} else {
			if db.IsNotFoundError(err) {
				contentDTO, err = s.CreateBookmarkContent(ctx, tx, dto)
				if err != nil {
					return nil, fmt.Errorf("failed to create bookmark content: %w", err)
				}
			} else {
				// return other errors
				return nil, fmt.Errorf("failed to check existing bookmark for url '%s': %w", dto.URL, err)
			}
		}
	}
	// Create the bookmark entry linking user to content
	bookmark, err := s.dao.CreateBookmark(ctx, tx, db.CreateBookmarkParams{
		UserID:    pgtype.UUID{Bytes: userId, Valid: true},
		ContentID: pgtype.UUID{Bytes: contentDTO.ID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create bookmark: %w", err)
	}

	// Convert database model to DTO
	var bookmarkDTO BookmarkDTO

	bookmarkDTO.Load(&bookmark)

	return &bookmarkDTO, nil
}

func (s *Service) GetBookmarkWithContent(ctx context.Context, tx db.DBTX, userId, id uuid.UUID) (*BookmarkDTO, error) {
	bookmark, err := s.dao.GetBookmarkWithContent(ctx, tx, db.GetBookmarkWithContentParams{
		ID:     id,
		UserID: pgtype.UUID{Bytes: userId, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	var result BookmarkDTO

	result.LoadWithContent(&bookmark)

	return &result, nil
}

func (s *Service) ListBookmarks(ctx context.Context, tx db.DBTX, userID uuid.UUID, filters []string, query string, limit, offset int32) ([]BookmarkDTO, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}

	if offset < 0 {
		offset = 0
	}

	// Use List instead of Search if no query provided since Search has worse performance
	if query != "" {
		return s.SearchBookmarks(ctx, tx, userID, filters, query, limit, offset)
	}

	domains, contentTypes, tags := parseListFilter(filters)
	totalCount := int64(0)
	bs, err := s.dao.ListBookmarks(ctx, tx, db.ListBookmarksParams{
		UserID:  pgtype.UUID{Bytes: userID, Valid: true},
		Limit:   limit,
		Offset:  offset,
		Domains: domains,
		Types:   contentTypes,
		Tags:    tags,
	})
	if err != nil {
		return nil, totalCount, fmt.Errorf("failed to list bookmarks: %w", err)
	}

	dtos := loadListBookmarks(bs)

	if len(bs) > 0 {
		totalCount = bs[0].TotalCount
	}

	return dtos, totalCount, nil
}

func (s *Service) SearchBookmarks(ctx context.Context, tx db.DBTX, userID uuid.UUID, filters []string, query string, limit, offset int32) ([]BookmarkDTO, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}

	if offset < 0 {
		offset = 0
	}

	domains, contentTypes, tags := parseListFilter(filters)
	totalCount := int64(0)

	bs, err := s.dao.SearchBookmarks(ctx, tx, db.SearchBookmarksParams{
		UserID:  pgtype.UUID{Bytes: userID, Valid: true},
		Limit:   limit,
		Offset:  offset,
		Domains: domains,
		Types:   contentTypes,
		Tags:    tags,
		Query: pgtype.Text{
			String: query,
			Valid:  query != "",
		},
	})
	if err != nil {
		return nil, totalCount, fmt.Errorf("failed to search bookmarks: %w", err)
	}

	dtos := loadSearchBookmarks(bs)

	if len(bs) > 0 {
		totalCount = bs[0].TotalCount
	}

	return dtos, totalCount, nil
}

func (s *Service) DeleteBookmark(ctx context.Context, tx db.DBTX, userId, id uuid.UUID) error {
	logger.FromContext(ctx).Info("deleting bookmark", "id", id.String(), "user_id", userId.String())

	return s.dao.DeleteBookmark(ctx, tx, db.DeleteBookmarkParams{
		ID:     id,
		UserID: pgtype.UUID{Bytes: userId, Valid: true},
	})
}

func (s *Service) DeleteBookmarksByUser(ctx context.Context, tx db.DBTX, userId uuid.UUID) error {
	return s.dao.DeleteBookmarksByUser(ctx, tx, pgtype.UUID{Bytes: userId, Valid: true})
}

func (s *Service) UpdateBookmark(ctx context.Context, tx db.DBTX, userId, id uuid.UUID, dto *BookmarkDTO) (*BookmarkDTO, error) {
	bookmark, err := s.GetBookmarkWithContent(ctx, tx, userId, id)
	if err != nil {
		return nil, err
	}

	if dto.Content != nil {
		new := dto.Content
		old := bookmark.Content
		// Update content if it's changed
		if new.Content != "" {
			old.Content = new.Content
		}

		if new.Description != "" {
			old.Description = new.Description
		}

		if new.Html != "" {
			old.Html = new.Html
		}

		if new.Summary != "" {
			old.Summary = new.Summary
		}

		if _, err = s.UpdateBookmarkContent(ctx, tx, old); err != nil {
			return nil, fmt.Errorf("failed to update bookmark content: %w", err)
		}
	}

	dbo, err := s.dao.UpdateBookmark(ctx, tx, dto.DumpToUpdateParams())
	if err != nil {
		return nil, fmt.Errorf("failed to update bookmark: %w", err)
	}

	bookmark.Load(&dbo)

	return bookmark, nil
}
