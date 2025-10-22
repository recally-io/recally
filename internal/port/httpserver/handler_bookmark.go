package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"recally/internal/core/bookmarks"
	"recally/internal/core/queue"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader/fetcher"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

// BookmarkService defines operations for managing bookmarks.
type BookmarkService interface {
	ListBookmarks(ctx context.Context, tx db.DBTX, userID uuid.UUID, filters []string, query string, limit, offset int32) ([]bookmarks.BookmarkDTO, int64, error)
	CreateBookmark(ctx context.Context, tx db.DBTX, userId uuid.UUID, dto *bookmarks.BookmarkContentDTO) (*bookmarks.BookmarkDTO, error)
	GetBookmarkWithContent(ctx context.Context, tx db.DBTX, userId, id uuid.UUID) (*bookmarks.BookmarkDTO, error)
	UpdateBookmark(ctx context.Context, tx db.DBTX, userId, id uuid.UUID, bookmak *bookmarks.BookmarkDTO) (*bookmarks.BookmarkDTO, error)
	DeleteBookmark(ctx context.Context, tx db.DBTX, id, userID uuid.UUID) error
	DeleteBookmarksByUser(ctx context.Context, tx db.DBTX, userID uuid.UUID) error

	FetchContent(ctx context.Context, tx db.DBTX, id, userID uuid.UUID, opts fetcher.FetchOptions) (*bookmarks.BookmarkContentDTO, error)
	SummarierContent(ctx context.Context, tx db.DBTX, id, userID uuid.UUID) (*bookmarks.BookmarkContentDTO, error)

	ListTags(ctx context.Context, tx db.DBTX, userID uuid.UUID) ([]bookmarks.TagDTO, error)
	ListDomains(ctx context.Context, tx db.DBTX, userID uuid.UUID) ([]bookmarks.DomainDTO, error)

	GetBookmarkShare(ctx context.Context, tx db.DBTX, userID, bookmarkID uuid.UUID) (*bookmarks.BookmarkShareDTO, error)
	CreateBookmarkShare(ctx context.Context, tx db.DBTX, userID, bookmarkID uuid.UUID, expiresAt time.Time) (*bookmarks.BookmarkShareDTO, error)
	UpdateBookmarkShare(ctx context.Context, tx db.DBTX, userID, bookmarkID uuid.UUID, expiresAt time.Time) (*bookmarks.BookmarkShareDTO, error)
	DeleteBookmarkShare(ctx context.Context, tx db.DBTX, userID, bookmarkID uuid.UUID) error
}

// bookmarkServiceImpl implements BookmarkService.
type bookmarksHandler struct {
	service BookmarkService
	queue   *queue.Queue
}

func registerBookmarkHandlers(e *echo.Group, s *Service) {
	h := &bookmarksHandler{service: bookmarks.NewService(s.llm), queue: s.queue}
	g := e.Group("/bookmarks", authUserMiddleware())
	g.GET("", h.listBookmarks)
	g.POST("", h.createBookmark)
	g.DELETE("/", h.deleteUserBookmarks)
	g.GET("/:bookmark-id", h.getBookmark)
	g.PUT("/:bookmark-id", h.updateBookmark)
	g.DELETE("/:bookmark-id", h.deleteBookmark)
	g.POST("/:bookmark-id/refresh", h.refreshBookmark)
	g.GET("/tags", h.listTags)
	g.GET("/domains", h.listDomains)

	// Updated sharing endpoints
	g.GET("/:bookmark-id/share", h.getBookmarkShare)
	g.POST("/:bookmark-id/share", h.createBookmarkShare)
	g.PUT("/:bookmark-id/share", h.updateBookmarkShare)
	g.DELETE("/:bookmark-id/share", h.deleteBookmarkShare)
}

type listBookmarksRequest struct {
	Limit  int32    `query:"limit" validate:"min=1,max=100"`
	Offset int32    `query:"offset" validate:"min=0"`
	Filter []string `query:"filter"` // filter=category:article;type:rss
	Query  string   `query:"query"`  // query=keyword
}

type listBookmarksResponse struct {
	Bookmarks []bookmarks.BookmarkDTO `json:"bookmarks"`
	Total     int64                   `json:"total"`
	Limit     int32                   `json:"limit"`
	Offset    int32                   `json:"offset"`
}

// listBookmarks handles GET /bookmarks
//
//	@Summary		List Bookmarks
//	@Description	Lists bookmarks for a user with pagination
//	@Tags			Bookmarks
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int										false	"Number of items per page"	default(10)
//	@Param			offset	query		int										false	"Number of items to skip"	default(0)
//	@Success		200		{object}	JSONResult{data=listBookmarksResponse}	"Success"
//	@Failure		401		{object}	JSONResult{data=nil}					"Unauthorized"
//	@Failure		500		{object}	JSONResult{data=nil}					"Internal Server Error"
//	@Router			/bookmarks [get]
func (h *bookmarksHandler) listBookmarks(c echo.Context) error {
	ctx := c.Request().Context()

	// Get pagination parameters with defaults
	req := new(listBookmarksRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	bookmarks, total, err := h.service.ListBookmarks(ctx, tx, user.ID, req.Filter, req.Query, req.Limit, req.Offset)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, listBookmarksResponse{
		Bookmarks: bookmarks,
		Total:     total,
		Limit:     req.Limit,
		Offset:    req.Offset,
	})
}

// listTags handles GET /bookmarks/tags
//
//	@Summary		List Tags
//	@Description	Lists all tags for a user's bookmarks
//	@Tags			Bookmarks
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	JSONResult{data=[]bookmarks.TagDTO}	"Success - Returns array of tags with counts"
//	@Failure		401	{object}	JSONResult{data=nil}				"Unauthorized"
//	@Failure		500	{object}	JSONResult{data=nil}				"Internal Server Error"
//	@Router			/bookmarks/tags [get]
func (h *bookmarksHandler) listTags(c echo.Context) error {
	ctx := c.Request().Context()

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	tags, err := cache.RunInCache(ctx, cache.MemCache,
		cache.NewCacheKey("bookmarks", "tags"),
		time.Minute,
		func() (*[]bookmarks.TagDTO, error) {
			t, err := h.service.ListTags(ctx, tx, user.ID)
			if err != nil {
				return nil, err
			}

			return &t, nil
		})
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, tags)
}

// listDomains handles GET /bookmarks/domains
//
//	@Summary		List Domains
//	@Description	Lists all domains from user's bookmarks
//	@Tags			Bookmarks
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	JSONResult{data=[]bookmarks.DomainDTO}	"Success - Returns array of domains with counts"
//	@Failure		401	{object}	JSONResult{data=nil}					"Unauthorized"
//	@Failure		500	{object}	JSONResult{data=nil}					"Internal Server Error"
//	@Router			/bookmarks/domains [get]
func (h *bookmarksHandler) listDomains(c echo.Context) error {
	ctx := c.Request().Context()

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	domains, err := cache.RunInCache(ctx, cache.MemCache,
		cache.NewCacheKey("bookmarks", "domains"),
		time.Minute,
		func() (*[]bookmarks.DomainDTO, error) {
			d, err := h.service.ListDomains(ctx, tx, user.ID)
			if err != nil {
				return nil, err
			}

			return &d, nil
		})
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, domains)
}

type createBookmarkRequest struct {
	URL         string   `json:"url,omitempty" validate:"omitempty,url"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Content     string   `json:"content,omitempty"`
	HTML        string   `json:"html,omitempty"`

	Type     bookmarks.ContentType              `json:"type" validate:"required,oneof=bookmark pdf epub image audio video"`
	S3Key    string                             `json:"s3_key,omitempty"`
	Metadata *bookmarks.BookmarkContentMetadata `json:"metadata,omitempty"`
}

// createBookmark handles POST /bookmarks
//
//	@Summary		Create Bookmark
//	@Description	Creates a new bookmark for the user
//	@Tags			Bookmarks
//	@Accept			json
//	@Produce		json
//	@Param			bookmark	body		createBookmarkRequest					true	"Bookmark to create"
//	@Success		201			{object}	JSONResult{data=bookmarks.BookmarkDTO}	"Created"
//	@Failure		400			{object}	JSONResult{data=nil}					"Bad Request"
//	@Failure		401			{object}	JSONResult{data=nil}					"Unauthorized"
//	@Failure		500			{object}	JSONResult{data=nil}					"Internal Server Error"
//	@Router			/bookmarks [post]
func (h *bookmarksHandler) createBookmark(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(createBookmarkRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	bookmarkContent := &bookmarks.BookmarkContentDTO{
		UserID:   user.ID,
		URL:      req.URL,
		Type:     req.Type,
		Title:    req.Title,
		Tags:     req.Tags,
		Content:  req.Content,
		Html:     req.HTML,
		S3Key:    req.S3Key,
		Metadata: req.Metadata,
	}

	bookmark, err := h.service.CreateBookmark(ctx, tx, user.ID, bookmarkContent)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if req.Type != bookmarks.ContentTypeBookmark && req.Type != bookmarks.ContentTypeImage {
		return JsonResponse(c, http.StatusCreated, bookmark)
	}

	result, err := h.queue.InsertTx(ctx, tx, queue.CrawlerWorkerArgs{
		ID:           bookmark.ID,
		UserID:       bookmark.UserID,
		FetchOptions: fetcher.FetchOptions{FecherType: fetcher.TypeHttp},
	}, nil)
	if err != nil {
		logger.FromContext(ctx).Error("failed to insert job", "err", err)
	} else {
		logger.FromContext(ctx).Info("success inserted job", "result", result, "err", err)
	}

	return JsonResponse(c, http.StatusCreated, bookmark)
}

type getBookmarkRequest struct {
	BookmarkID uuid.UUID `param:"bookmark-id" validate:"required,uuid4"`
}

// getBookmark handles GET /bookmarks/:bookmark-id
//
//	@Summary		Get Bookmark
//	@Description	Gets a specific bookmark by ID
//	@Tags			Bookmarks
//	@Accept			json
//	@Produce		json
//	@Param			bookmark-id	path		string									true	"Bookmark ID"
//	@Success		200			{object}	JSONResult{data=bookmarks.BookmarkDTO}	"Success"
//	@Failure		400			{object}	JSONResult{data=nil}					"Bad Request"
//	@Failure		401			{object}	JSONResult{data=nil}					"Unauthorized"
//	@Failure		404			{object}	JSONResult{data=nil}					"Not Found"
//	@Failure		500			{object}	JSONResult{data=nil}					"Internal Server Error"
//	@Router			/bookmarks/{bookmark-id} [get]
func (h *bookmarksHandler) getBookmark(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(getBookmarkRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	bookmark, err := h.service.GetBookmarkWithContent(ctx, tx, user.ID, req.BookmarkID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if bookmark == nil {
		return ErrorResponse(c, http.StatusNotFound, fmt.Errorf("bookmark not found"))
	}

	return JsonResponse(c, http.StatusOK, bookmark)
}

type updateBookmarkRequest struct {
	BookmarkID uuid.UUID                         `param:"bookmark-id" validate:"required,uuid4"`
	Summary    string                            `json:"summary"`
	Content    string                            `json:"content"`
	HTML       string                            `json:"html"`
	Metadata   bookmarks.BookmarkContentMetadata `json:"metadata"`
}

// updateBookmark handles PUT /bookmarks/:bookmark-id
//
//	@Summary		Update Bookmark
//	@Description	Updates an existing bookmark
//	@Tags			Bookmarks
//	@Accept			json
//	@Produce		json
//	@Param			bookmark-id	path		string									true	"Bookmark ID"
//	@Param			bookmark	body		updateBookmarkRequest					true	"Updated bookmark data"
//	@Success		200			{object}	JSONResult{data=bookmarks.BookmarkDTO}	"Success"
//	@Failure		400			{object}	JSONResult{data=nil}					"Bad Request"
//	@Failure		401			{object}	JSONResult{data=nil}					"Unauthorized"
//	@Failure		404			{object}	JSONResult{data=nil}					"Not Found"
//	@Failure		500			{object}	JSONResult{data=nil}					"Internal Server Error"
//	@Router			/bookmarks/{bookmark-id} [put]
func (h *bookmarksHandler) updateBookmark(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(updateBookmarkRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	bookmarkContent := &bookmarks.BookmarkContentDTO{
		ID:     req.BookmarkID,
		UserID: user.ID,
	}

	if req.Summary != "" {
		bookmarkContent.Summary = req.Summary
	}

	if req.Content != "" {
		bookmarkContent.Content = req.Content
	}

	if req.HTML != "" {
		bookmarkContent.Html = req.HTML
	}

	bookmark := &bookmarks.BookmarkDTO{
		Content: bookmarkContent,
	}

	updated, err := h.service.UpdateBookmark(ctx, tx, user.ID, req.BookmarkID, bookmark)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if updated == nil {
		return ErrorResponse(c, http.StatusNotFound, fmt.Errorf("bookmark not found"))
	}

	return JsonResponse(c, http.StatusOK, updated)
}

// deleteBookmark handles DELETE /bookmarks/:bookmark-id
//
//	@Summary		Delete Bookmark
//	@Description	Deletes a bookmark
//	@Tags			Bookmarks
//	@Accept			json
//	@Produce		json
//	@Param			bookmark-id	path		string					true	"Bookmark ID"
//	@Success		204			{object}	JSONResult{data=nil}	"No Content"
//	@Failure		400			{object}	JSONResult{data=nil}	"Bad Request"
//	@Failure		401			{object}	JSONResult{data=nil}	"Unauthorized"
//	@Failure		500			{object}	JSONResult{data=nil}	"Internal Server Error"
//	@Router			/bookmarks/{bookmark-id} [delete]
func (h *bookmarksHandler) deleteBookmark(c echo.Context) error {
	ctx := c.Request().Context()

	bookmarkID, err := uuid.Parse(c.Param("bookmark-id"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if err := h.service.DeleteBookmark(ctx, tx, user.ID, bookmarkID); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusNoContent, nil)
}

type deleteUserBookmarksRequest struct {
	UserID uuid.UUID `query:"user-id" validate:"required,uuid4"`
}

// deleteUserBookmarks handles DELETE /bookmarks/:user-id/bookmarks
//
//	@Summary		Delete User Bookmarks
//	@Description	Deletes all bookmarks for a user
//	@Tags			Bookmarks
//	@Accept			json
//	@Produce		json
//	@Param			user-id	query		string					true	"User ID"
//	@Success		204		{object}	JSONResult{data=nil}	"No Content"
//	@Failure		400		{object}	JSONResult{data=nil}	"Bad Request"
//	@Failure		401		{object}	JSONResult{data=nil}	"Unauthorized"
//	@Failure		500		{object}	JSONResult{data=nil}	"Internal Server Error"
//	@Router			/bookmarks [delete]
func (h *bookmarksHandler) deleteUserBookmarks(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(deleteUserBookmarksRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	// Security check - ensure user can only delete their own bookmarks
	if user.ID != req.UserID {
		return ErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
	}

	if err := h.service.DeleteBookmarksByUser(ctx, tx, req.UserID); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusNoContent, nil)
}

type refreshBookmarkRequest struct {
	BookmarkID        uuid.UUID `param:"bookmark-id" validate:"required,uuid4"`
	Fetcher           string    `json:"fetcher" validate:"omitempty,oneof=http jinaReader browser"`
	IsProxyImage      bool      `json:"is_proxy_image"`
	RegenerateSummary bool      `json:"regenerate_summary"`
}

// refreshBookmark handles POST /bookmarks/:bookmark-id/refresh
//
//	@Summary		Refresh Bookmark
//	@Description	Refreshes bookmark content and/or regenerates AI summary
//	@Tags			Bookmarks
//	@Accept			json
//	@Produce		json
//	@Param			bookmark-id	path		string									true	"Bookmark ID"
//	@Param			request		body		refreshBookmarkRequest					true	"Refresh options"
//	@Success		200			{object}	JSONResult{data=bookmarks.BookmarkDTO}	"Success"
//	@Failure		400			{object}	JSONResult{data=nil}					"Bad Request"
//	@Failure		401			{object}	JSONResult{data=nil}					"Unauthorized"
//	@Failure		404			{object}	JSONResult{data=nil}					"Not Found"
//	@Failure		500			{object}	JSONResult{data=nil}					"Internal Server Error"
//	@Router			/bookmarks/{bookmark-id}/refresh [post]
func (h *bookmarksHandler) refreshBookmark(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(refreshBookmarkRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	var content *bookmarks.BookmarkContentDTO

	if req.Fetcher != "" {
		content, err = h.service.FetchContent(ctx, tx, req.BookmarkID, user.ID, fetcher.FetchOptions{
			FecherType:   fetcher.FecherType(req.Fetcher),
			IsProxyImage: req.IsProxyImage,
			Force:        true,
		})
		if err != nil {
			return ErrorResponse(c, http.StatusInternalServerError, err)
		}
	}

	if req.RegenerateSummary {
		content, err = h.service.SummarierContent(ctx, tx, req.BookmarkID, user.ID)
		if err != nil {
			return ErrorResponse(c, http.StatusInternalServerError, err)
		}
	}

	return JsonResponse(c, http.StatusOK, content)
}
