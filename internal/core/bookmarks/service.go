package bookmarks

import (
	"context"
	"fmt"
	"net/url"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader"
	"recally/internal/pkg/webreader/fetcher"
	"recally/internal/pkg/webreader/processor"
	"recally/internal/pkg/webreader/reader"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	dao DAO
	llm *llms.LLM
}

func NewService(llm *llms.LLM) *Service {
	return &Service{
		dao: db.New(),
		llm: llm,
	}
}

// CreateBookmark creates a new bookmark with content fetching and embedding generation
func (s *Service) Create(ctx context.Context, tx db.DBTX, dto *ContentDTO) (*ContentDTO, error) {
	// Validate URL
	if _, err := url.ParseRequestURI(dto.URL); err != nil {
		return nil, fmt.Errorf("%w: invalid URL", ErrInvalidInput)
	}

	// Check for existing bookmark
	isExisting, err := s.dao.IsContentExistWithURL(ctx, tx, db.IsContentExistWithURLParams{
		Url:    pgtype.Text{String: dto.URL, Valid: true},
		UserID: dto.UserID,
	})
	if err != nil && !db.IsNotFoundError(err) {
		return nil, fmt.Errorf("failed to check existing bookmark for url '%s': %w", dto.URL, err)
	}
	if isExisting {
		return nil, fmt.Errorf("%w, id: %s", ErrDuplicate, dto.URL)
	}

	// create content
	c, err := s.dao.CreateContent(ctx, tx, dto.Dump())
	if err != nil {
		return nil, fmt.Errorf("failed to create bookmark for url '%s': %w", dto.URL, err)
	}
	dto.Load(&c)

	if len(dto.Tags) > 0 {
		if err := s.linkContentTags(ctx, tx, []string{}, dto.Tags, c.ID, dto.UserID); err != nil {
			return nil, err
		}
	}
	return dto, nil
}

// GetBookmark retrieves a bookmark by ID
func (s *Service) Get(ctx context.Context, tx db.DBTX, id, userID uuid.UUID) (*ContentDTO, error) {
	c, err := s.dao.GetContent(ctx, tx, db.GetContentParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get content by id '%s': %w", id.String(), err)
	}

	var dto ContentDTO
	dto.LoadWithTags(&c)
	// Clear content and HTML
	dto.HTML = ""

	// add share content info
	if sc, err := s.GetShareContent(ctx, tx, id); err == nil {
		dto.Metadata.Share = sc
	}

	return &dto, nil
}

func parseListFilter(filters []string) (domains, contentTypes, tags []string) {
	if len(filters) == 0 {
		return
	}

	domains = make([]string, 0)
	contentTypes = make([]string, 0)
	tags = make([]string, 0)

	// Parse filter=category:article;type:rss
	for _, part := range filters {
		kv := strings.Split(part, ":")
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "domain":
			domains = append(domains, kv[1])
		case "type":
			contentTypes = append(contentTypes, kv[1])
		case "tag":
			tags = append(tags, kv[1])
		}
	}
	if len(domains) == 0 {
		domains = nil
	}
	if len(contentTypes) == 0 {
		contentTypes = nil
	}
	if len(tags) == 0 {
		tags = nil
	}
	return
}

// ListBookmarks retrieves a paginated list of bookmarks for a user
func (s *Service) List(ctx context.Context, tx db.DBTX, userID uuid.UUID, filters []string, query string, limit, offset int32) ([]*ContentDTO, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	// Use List instead of Search if no query provided since Search has worse performance
	if query != "" {
		return s.Search(ctx, tx, userID, filters, query, limit, offset)
	}

	domains, contentTypes, tags := parseListFilter(filters)
	totalCount := int64(0)
	cs, err := s.dao.ListContents(ctx, tx, db.ListContentsParams{
		UserID:  userID,
		Limit:   limit,
		Offset:  offset,
		Domains: domains,
		Types:   contentTypes,
		Tags:    tags,
	})

	dtos := make([]*ContentDTO, 0, len(cs))
	for _, c := range cs {
		var dto ContentDTO
		dto.LoadWithTagsAndTotalCount(&c)
		dto.HTML = ""
		dto.Content = ""
		dto.Summary = ""
		dtos = append(dtos, &dto)
		totalCount = c.TotalCount
	}
	return dtos, totalCount, err
}

func (s *Service) Search(ctx context.Context, tx db.DBTX, userID uuid.UUID, filters []string, query string, limit, offset int32) ([]*ContentDTO, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	domains, contentTypes, tags := parseListFilter(filters)
	totalCount := int64(0)
	cs, err := s.dao.SearchContentsWithFilter(ctx, tx, db.SearchContentsWithFilterParams{
		UserID:  userID,
		Limit:   limit,
		Offset:  offset,
		Domains: domains,
		Types:   contentTypes,
		Tags:    tags,
		Query:   pgtype.Text{String: query, Valid: query != ""},
	})

	dtos := make([]*ContentDTO, 0, len(cs))
	for _, c := range cs {
		var dto ContentDTO
		dto.LoadWithTagsAndTotalCountFromSearch(&c)
		dto.HTML = ""
		dto.Content = ""
		dto.Summary = ""
		dtos = append(dtos, &dto)
		totalCount = c.TotalCount
	}
	return dtos, totalCount, err
}

type TagDTO struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

func (s *Service) ListTags(ctx context.Context, tx db.DBTX, userID uuid.UUID) ([]TagDTO, error) {
	tags, err := s.dao.ListTagsByUser(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags for user '%s': %w", userID.String(), err)
	}

	tagsList := make([]TagDTO, 0, len(tags))
	for _, tag := range tags {
		tagsList = append(tagsList, TagDTO{
			Name:  tag.Name,
			Count: tag.Count,
		})
	}
	return tagsList, nil
}

type DomainDTO struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

func (s *Service) ListDomains(ctx context.Context, tx db.DBTX, userID uuid.UUID) ([]DomainDTO, error) {
	domains, err := s.dao.ListContentDomains(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains for user '%s': %w", userID.String(), err)
	}
	tags := make([]DomainDTO, 0, len(domains))
	for _, domain := range domains {
		tags = append(tags, DomainDTO{
			Name:  domain.Domain.String,
			Count: domain.Count,
		})
	}
	return tags, nil
}

// UpdateBookmark updates an existing bookmark, PUT full update
func (s *Service) Update(ctx context.Context, tx db.DBTX, id, userID uuid.UUID, dto *ContentDTO) (*ContentDTO, error) {
	_, err := s.Get(ctx, tx, id, userID)
	if err != nil {
		return nil, err
	}
	updateParams := dto.DumpToUpdateParams()
	c, err := s.dao.UpdateContent(ctx, tx, updateParams)
	if err != nil {
		return nil, err
	}

	dto.Load(&c)
	return dto, nil
}

// DeleteBookmark removes a bookmark
func (s *Service) Delete(ctx context.Context, tx db.DBTX, id, userID uuid.UUID) error {
	_, err := s.Get(ctx, tx, id, userID)
	if err != nil {
		return err
	}

	return s.dao.DeleteContent(ctx, tx, db.DeleteContentParams{
		ID:     id,
		UserID: userID,
	})
}

// DeleteUserBookmarks removes all bookmarks for a user
func (s *Service) DeleteUserBookmarks(ctx context.Context, tx db.DBTX, userID uuid.UUID) error {
	return s.dao.DeleteContentsByUser(ctx, tx, userID)
}

func (s *Service) FetchContent(ctx context.Context, tx db.DBTX, id, userID uuid.UUID, opts fetcher.FetchOptions) (*ContentDTO, error) {
	dto, err := s.Get(ctx, tx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmark by id '%s': %w", id.String(), err)
	}

	content, err := s.FetchContentWithCache(ctx, dto.URL, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch content: %w", err)
	}

	dto.Content = content.Markwdown
	dto.Title = content.Title
	dto.HTML = content.Html

	// Update metadata
	dto.Metadata.Author = content.Author
	dto.Metadata.SiteName = content.SiteName
	dto.Metadata.Description = content.Description

	dto.Metadata.Cover = content.Cover
	dto.Metadata.Favicon = content.Favicon
	if content.Cover != "" {
		dto.Metadata.Image = content.Cover
	} else {
		dto.Metadata.Image = content.Favicon
	}

	if content.PublishedTime != nil {
		dto.Metadata.PublishedAt = *content.PublishedTime
	}
	return s.Update(ctx, tx, id, userID, dto)
}

func (s *Service) FetchContentWithCache(ctx context.Context, uri string, opts fetcher.FetchOptions) (*webreader.Content, error) {
	return cache.RunInCache(ctx, cache.DefaultDBCache,
		cache.NewCacheKey(fmt.Sprintf("WebReader-%s", opts.String()), uri), 24*time.Hour, func() (*webreader.Content, error) {
			u, err := url.Parse(uri)
			if err != nil {
				return nil, fmt.Errorf("invalid url '%s': %w", uri, err)
			}

			// read the content using jina reader
			reader, err := reader.New(u.Host, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to create reader: %w", err)
			}
			content, err := reader.Read(ctx, uri)
			if err != nil {
				return nil, fmt.Errorf("failed to read content: %w", err)
			}
			return content, nil
		})
}

func (s *Service) SummarierContent(ctx context.Context, tx db.DBTX, id, userID uuid.UUID) (*ContentDTO, error) {
	user, err := auth.LoadUser(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	dto, err := s.Get(ctx, tx, id, user.ID)
	if err != nil {
		return nil, err
	}

	content := &webreader.Content{
		Markwdown: dto.Content,
	}

	summarier := processor.NewSummaryProcessor(s.llm, processor.WithSummaryOptionUser(user))

	if len(content.Markwdown) < 1000 {
		logger.FromContext(ctx).Info("content is too short to summarise")
		return dto, nil
	}

	if err := summarier.Process(ctx, content); err != nil {
		logger.Default.Error("failed to generate summary", "err", err)
	} else {
		tags, summary := parseTagsFromSummary(content.Summary)
		if len(tags) > 0 {
			if err := s.linkContentTags(ctx, tx, dto.Tags, tags, id, userID); err != nil {
				return nil, err
			}
		}
		dto.Summary = summary
	}
	return s.Update(ctx, tx, id, userID, dto)
}

// parseTagsFromSummary extracts tags from a string and returns the tags array and the string without tags
func parseTagsFromSummary(input string) ([]string, string) {
	// Regular expression to match the tags section
	tagsRegex := regexp.MustCompile(`(?s)<tags>.*?</tags>`)

	// Find tags section
	tagsSection := tagsRegex.FindString(input)

	// If no tags section found, return empty array and original string
	if tagsSection == "" {
		return []string{}, input
	}

	// Extract content between tags
	content := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(tagsSection, "<tags>"), "</tags>"))

	// Split content by whitespace
	words := strings.Fields(content)

	// Process valid tags
	tagMap := make(map[string]bool) // Use map to ensure uniqueness
	var tags []string

	for _, word := range words {
		if strings.HasPrefix(word, "#") {
			tag := strings.TrimPrefix(word, "#")
			if tag != "" && !tagMap[tag] {
				tagMap[tag] = true
				tags = append(tags, tag)
			}
		}
	}

	// Remove tags section from original string
	cleanedString := strings.TrimSpace(tagsRegex.ReplaceAllString(input, "\n"))

	return tags, cleanedString
}

func (s *Service) linkContentTags(ctx context.Context, tx db.DBTX, originTags, newTags []string, contentID, userID uuid.UUID) error {
	// create tags if not exist
	allExistingTags, err := s.dao.ListExistingTagsByTags(ctx, tx, db.ListExistingTagsByTagsParams{
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
		if _, err := s.dao.CreateContentTag(ctx, tx, db.CreateContentTagParams{
			Name:   tag,
			UserID: userID,
		}); err != nil {
			return fmt.Errorf("failed to create tag '%s': %w", tag, err)
		}
	}

	// link content with tags that not linked before
	contentExistingTags, err := s.dao.ListContentTags(ctx, tx, db.ListContentTagsParams{
		ContentID: contentID,
		UserID:    userID,
	})
	if err != nil {
		return fmt.Errorf("failed to list content tags: %w", err)
	}
	newLinkedTags := make([]string, 0)
	for _, tag := range newTags {
		if !slices.Contains(contentExistingTags, tag) {
			newLinkedTags = append(newLinkedTags, tag)
		}
	}
	if err := s.dao.LinkContentWithTags(ctx, tx, db.LinkContentWithTagsParams{
		ContentID: contentID,
		Column2:   newLinkedTags,
		UserID:    userID,
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

	if err := s.dao.UnLinkContentWithTags(ctx, tx, db.UnLinkContentWithTagsParams{
		ContentID: contentID,
		Column2:   removedTags,
		UserID:    userID,
	}); err != nil {
		return fmt.Errorf("failed to unlink tags with content: %w", err)
	}

	logger.FromContext(ctx).Info("link content with tags", "content_id", contentID, "new_tags", newTags, "origin_tags", originTags, "removed_tags", removedTags, "new_linked_tags", newLinkedTags)
	return nil
}
