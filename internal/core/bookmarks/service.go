package bookmarks

import (
	"context"
	"fmt"
	"net/url"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader"
	"recally/internal/pkg/webreader/fetcher"
	"recally/internal/pkg/webreader/processor"
	"recally/internal/pkg/webreader/reader"
	"strings"

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
	if isExisting {
		return nil, fmt.Errorf("%w, id: %s", ErrDuplicate, dto.URL)
	}

	if !db.IsNotFoundError(err) {
		return nil, fmt.Errorf("failed to check existing bookmark for url '%s': %w", dto.URL, err)
	}

	// create content
	c, err := s.dao.CreateContent(ctx, tx, dto.Dump())
	if err != nil {
		return nil, fmt.Errorf("failed to create bookmark for url '%s': %w", dto.URL, err)
	}
	dto.Load(&c)

	if len(dto.Tags) > 0 {
		// create tags
		for _, tag := range dto.Tags {
			if _, err := s.dao.CreateContentTag(ctx, tx, db.CreateContentTagParams{
				Name:   tag,
				UserID: dto.UserID,
			}); err != nil {
				logger.FromContext(ctx).Error("failed to create tag", "err", err, "content_id", c.ID, "tag", tag)
			}
		}
		// link content with tags
		if err := s.dao.LinkContentWithTags(ctx, tx, db.LinkContentWithTagsParams{
			ContentID: c.ID,
			Column2:   dto.Tags,
			UserID:    dto.UserID,
		}); err != nil {
			logger.FromContext(ctx).Error("failed to link tags with content", "err", err, "content_id", c.ID, "tags", dto.Tags)
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
	return &dto, nil
}

func parseListFilter(filter string) (domains, contentTypes, tags []string) {
	if filter == "" {
		return
	}

	domains = make([]string, 0)
	contentTypes = make([]string, 0)
	tags = make([]string, 0)

	// Parse filter=category:article;type:rss
	parts := strings.Split(filter, ";")
	for _, part := range parts {
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
	return
}

// ListBookmarks retrieves a paginated list of bookmarks for a user
func (s *Service) List(ctx context.Context, tx db.DBTX, userID uuid.UUID, filter, query string, limit, offset int32) ([]*ContentDTO, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	domains, contentTypes, tags := parseListFilter(filter)
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
			Count: int64(tag.UsageCount.Int32),
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

func (s *Service) Refresh(ctx context.Context, tx db.DBTX, id, userID uuid.UUID, fetcherType fetcher.FecherType, regenerateSummary bool) (*ContentDTO, error) {
	var dto *ContentDTO
	var err error

	if fetcherType != fetcher.TypeNil {
		dto, err = s.FetchContent(ctx, tx, id, userID, fetcherType)
		if err != nil {
			return nil, err
		}
	}

	if regenerateSummary {
		dto, err = s.SummarierContent(ctx, tx, id, userID)
		if err != nil {
			return nil, err
		}
	}

	return dto, nil
}

func (s *Service) FetchContent(ctx context.Context, tx db.DBTX, id, userID uuid.UUID, fetcherType fetcher.FecherType) (*ContentDTO, error) {
	dto, err := s.Get(ctx, tx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmark by id '%s': %w", id.String(), err)
	}

	reader, err := reader.New(fetcherType, dto.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %w", err)
	}
	content, err := reader.Read(ctx, dto.URL)
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

func (s *Service) SummarierContent(ctx context.Context, tx db.DBTX, id, userID uuid.UUID) (*ContentDTO, error) {
	user, err := auth.LoadUserFromContext(ctx)
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

	summarier := newSummarier(s.llm, user)

	if err := summarier.Process(ctx, content); err != nil {
		logger.Default.Error("failed to generate summary", "err", err)
	} else {
		dto.Summary = content.Summary
	}

	return s.Update(ctx, tx, id, userID, dto)
}

func newSummarier(llm *llms.LLM, user *auth.UserDTO) *processor.SummaryProcessor {
	summaryOptions := make([]processor.SummaryOption, 0)
	if user.Settings.SummaryOptions.Prompt != "" {
		summaryOptions = append(summaryOptions, processor.WithSummaryOptionPrompt(user.Settings.SummaryOptions.Prompt))
	}
	if user.Settings.SummaryOptions.Model != "" {
		summaryOptions = append(summaryOptions, processor.WithSummaryOptionModel(user.Settings.SummaryOptions.Model))
	}
	if user.Settings.SummaryOptions.Language != "" {
		summaryOptions = append(summaryOptions, processor.WithSummaryOptionLanguage(user.Settings.SummaryOptions.Language))
	}

	return processor.NewSummaryProcessor(llm, summaryOptions...)
}
