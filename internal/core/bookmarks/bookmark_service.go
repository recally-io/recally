// Package bookmarks provides functionality for managing user bookmarks and their content
package bookmarks

import (
	"context"
	"fmt"
	"net/url"
	"recally/internal/pkg/db"

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
	if _, err := url.ParseRequestURI(dto.URL); err != nil {
		return nil, fmt.Errorf("%w: invalid URL", ErrInvalidInput)
	}

	// Track if content already exists in the database
	isContentExist := false

	// Check if bookmark content already exists for this URL and user
	content, err := s.dao.GetBookmarkContentByURL(ctx, tx, db.GetBookmarkContentByURLParams{
		Url:    dto.URL,
		UserID: pgtype.UUID{Bytes: userId, Valid: true},
	})

	// Handle the response from content lookup
	if err == nil {
		isContentExist = true
	} else if !db.IsNotFoundError(err) {
		// Return error if it's not a "not found" error
		return nil, fmt.Errorf("failed to check existing bookmark for url '%s': %w", dto.URL, err)
	}

	if isContentExist {
		if content.UserID.Valid {
			// If content exists and belongs to a user, return duplicate error
			return nil, fmt.Errorf("%w, id: %s", ErrDuplicate, content.ID)
		}
	} else {
		// Create new content if it doesn't exist
		createBookmarkContentParams := dto.Dump()
		// Set UserID to Nil as this content can be shared
		createBookmarkContentParams.UserID = pgtype.UUID{Bytes: uuid.Nil, Valid: false}
		content, err = s.dao.CreateBookmarkContent(ctx, tx, createBookmarkContentParams)
		if err != nil {
			return nil, fmt.Errorf("failed to create new bookmark content: %w", err)
		}
	}

	// Create the bookmark entry linking user to content
	bookmark, err := s.dao.CreateBookmark(ctx, tx, db.CreateBookmarkParams{
		UserID:    pgtype.UUID{Bytes: userId, Valid: true},
		ContentID: pgtype.UUID{Bytes: content.ID, Valid: true},
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
	return s.dao.DeleteBookmark(ctx, tx, db.DeleteBookmarkParams{
		ID:     id,
		UserID: pgtype.UUID{Bytes: userId, Valid: true},
	})
}

func (s *Service) DeleteBookmarksByUser(ctx context.Context, tx db.DBTX, userId uuid.UUID) error {
	return s.dao.DeleteBookmarksByUser(ctx, tx, pgtype.UUID{Bytes: userId, Valid: true})
}

func (s *Service) UpdateBookmark(ctx context.Context, tx db.DBTX, userId uuid.UUID, id uuid.UUID, content *BookmarkContentDTO) (*BookmarkDTO, error) {
	bookmark, err := s.GetBookmarkWithContent(ctx, tx, userId, id)
	if err != nil {
		return nil, err
	}

	updateContent := bookmark.Content
	// Update content if it's changed
	if content.Content != "" {
		updateContent.Content = content.Content
	}
	if content.Description != "" {
		updateContent.Description = content.Description
	}
	if content.Html != "" {
		updateContent.Html = content.Html
	}
	if content.Summary != "" {
		updateContent.Summary = content.Summary
	}
	if _, err = s.UpdateBookmarkContent(ctx, tx, &updateContent); err != nil {
		return nil, err
	}
	return bookmark, nil
}
