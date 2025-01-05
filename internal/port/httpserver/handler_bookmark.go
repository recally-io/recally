package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"recally/internal/core/bookmarks"
	"recally/internal/core/queue"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader/fetcher"
	"time"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

// BookmarkService defines operations for managing bookmarks
type BookmarkService interface {
	Create(ctx context.Context, tx db.DBTX, dto *bookmarks.ContentDTO) (*bookmarks.ContentDTO, error)
	Get(ctx context.Context, tx db.DBTX, id, userID uuid.UUID) (*bookmarks.ContentDTO, error)
	List(ctx context.Context, tx db.DBTX, userID uuid.UUID, filter, query string, limit, offset int32) ([]*bookmarks.ContentDTO, int64, error)
	ListTags(ctx context.Context, tx db.DBTX, userID uuid.UUID) ([]bookmarks.TagDTO, error)
	ListDomains(ctx context.Context, tx db.DBTX, userID uuid.UUID) ([]bookmarks.DomainDTO, error)
	Update(ctx context.Context, tx db.DBTX, id, userID uuid.UUID, dto *bookmarks.ContentDTO) (*bookmarks.ContentDTO, error)
	Delete(ctx context.Context, tx db.DBTX, id, userID uuid.UUID) error
	DeleteUserBookmarks(ctx context.Context, tx db.DBTX, userID uuid.UUID) error
	Refresh(ctx context.Context, tx db.DBTX, id, userID uuid.UUID, fetcher fetcher.FecherType, regenerateSummary bool) (*bookmarks.ContentDTO, error)
	FetchContent(ctx context.Context, tx db.DBTX, id, userID uuid.UUID, fetcherType fetcher.FecherType) (*bookmarks.ContentDTO, error)
	SummarierContent(ctx context.Context, tx db.DBTX, id, userID uuid.UUID) (*bookmarks.ContentDTO, error)
}

// bookmarkServiceImpl implements BookmarkService
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
}

type listBookmarksRequest struct {
	Limit  int32  `query:"limit" validate:"min=1,max=100"`
	Offset int32  `query:"offset" validate:"min=0"`
	Filter string `query:"filter"` // filter=category:article;type:rss
	Query  string `query:"query"`  // query=keyword
}

type listBookmarksResponse struct {
	Bookmarks []*bookmarks.ContentDTO `json:"bookmarks"`
	Total     int64                   `json:"total"`
	Limit     int32                   `json:"limit"`
	Offset    int32                   `json:"offset"`
}

// listBookmarks handles GET /bookmarks
// @Summary List Bookmarks
// @Description Lists bookmarks for a user with pagination
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page" default(10)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} JSONResult{data=listBookmarksResponse} "Success"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /bookmarks [get]
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

	bookmarks, total, err := h.service.List(ctx, tx, user.ID, req.Filter, req.Query, req.Limit, req.Offset)
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
// @Summary List Tags
// @Description Lists all tags for a user's bookmarks
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Success 200 {object} JSONResult{data=[]bookmarks.TagDTO} "Success - Returns array of tags with counts"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /bookmarks/tags [get]
func (h *bookmarksHandler) listTags(c echo.Context) error {
	ctx := c.Request().Context()

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	tags, err := cache.RunInCache[[]bookmarks.TagDTO](ctx, cache.MemCache,
		cache.NewCacheKey("bookmarks", "tags"),
		5*time.Minute,
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
// @Summary List Domains
// @Description Lists all domains from user's bookmarks
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Success 200 {object} JSONResult{data=[]bookmarks.DomainDTO} "Success - Returns array of domains with counts"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /bookmarks/domains [get]
func (h *bookmarksHandler) listDomains(c echo.Context) error {
	ctx := c.Request().Context()

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	domains, err := cache.RunInCache[[]bookmarks.DomainDTO](ctx, cache.MemCache,
		cache.NewCacheKey("bookmarks", "domains"),
		5*time.Minute,
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
	URL      string             `json:"url" validate:"required,url"`
	Metadata bookmarks.Metadata `json:"metadata"`
}

// createBookmark handles POST /bookmarks
// @Summary Create Bookmark
// @Description Creates a new bookmark for the user
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param bookmark body createBookmarkRequest true "Bookmark to create"
// @Success 201 {object} JSONResult{data=bookmarks.ContentDTO} "Created"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /bookmarks [post]
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

	bookmark := &bookmarks.ContentDTO{
		UserID:   user.ID,
		URL:      req.URL,
		Metadata: req.Metadata,
		Type:     bookmarks.ContentTypeBookmark,
	}

	created, err := h.service.Create(ctx, tx, bookmark)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	result, err := h.queue.Insert(ctx, queue.CrawlerWorkerArgs{
		ID:          created.ID,
		UserID:      created.UserID,
		FetcherName: fetcher.TypeHttp,
	}, nil)
	if err != nil {
		logger.FromContext(ctx).Error("failed to insert job", "err", err)
	} else {
		logger.FromContext(ctx).Info("success inserted job", "result", result, "err", err)
	}
	return JsonResponse(c, http.StatusCreated, created)
}

type getBookmarkRequest struct {
	BookmarkID uuid.UUID `param:"bookmark-id" validate:"required,uuid4"`
}

// getBookmark handles GET /bookmarks/:bookmark-id
// @Summary Get Bookmark
// @Description Gets a specific bookmark by ID
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param bookmark-id path string true "Bookmark ID"
// @Success 200 {object} JSONResult{data=bookmarks.ContentDTO} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 404 {object} JSONResult{data=nil} "Not Found"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /bookmarks/{bookmark-id} [get]
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

	bookmark, err := h.service.Get(ctx, tx, req.BookmarkID, user.ID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	if bookmark == nil {
		return ErrorResponse(c, http.StatusNotFound, fmt.Errorf("bookmark not found"))
	}

	return JsonResponse(c, http.StatusOK, bookmark)
}

type updateBookmarkRequest struct {
	BookmarkID uuid.UUID          `param:"bookmark-id" validate:"required,uuid4"`
	Summary    string             `json:"summary"`
	Content    string             `json:"content"`
	HTML       string             `json:"html"`
	Metadata   bookmarks.Metadata `json:"metadata"`
}

// updateBookmark handles PUT /bookmarks/:bookmark-id
// @Summary Update Bookmark
// @Description Updates an existing bookmark
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param bookmark-id path string true "Bookmark ID"
// @Param bookmark body updateBookmarkRequest true "Updated bookmark data"
// @Success 200 {object} JSONResult{data=bookmarks.ContentDTO} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 404 {object} JSONResult{data=nil} "Not Found"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /bookmarks/{bookmark-id} [put]
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

	bookmark := &bookmarks.ContentDTO{
		ID:     req.BookmarkID,
		UserID: user.ID,
	}

	if req.Summary != "" {
		bookmark.Summary = req.Summary
	}

	if req.Content != "" {
		bookmark.Content = req.Content
	}

	if req.HTML != "" {
		bookmark.HTML = req.HTML
	}

	updated, err := h.service.Update(ctx, tx, req.BookmarkID, user.ID, bookmark)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	if updated == nil {
		return ErrorResponse(c, http.StatusNotFound, fmt.Errorf("bookmark not found"))
	}

	return JsonResponse(c, http.StatusOK, updated)
}

// deleteBookmark handles DELETE /bookmarks/:bookmark-id
// @Summary Delete Bookmark
// @Description Deletes a bookmark
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param bookmark-id path string true "Bookmark ID"
// @Success 204 {object} JSONResult{data=nil} "No Content"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /bookmarks/{bookmark-id} [delete]
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

	if err := h.service.Delete(ctx, tx, bookmarkID, user.ID); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusNoContent, nil)
}

type deleteUserBookmarksRequest struct {
	UserID uuid.UUID `query:"user-id" validate:"required,uuid4"`
}

// deleteUserBookmarks handles DELETE /bookmarks/:user-id/bookmarks
// @Summary Delete User Bookmarks
// @Description Deletes all bookmarks for a user
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param user-id query string true "User ID"
// @Success 204 {object} JSONResult{data=nil} "No Content"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /bookmarks [delete]
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

	if err := h.service.DeleteUserBookmarks(ctx, tx, req.UserID); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusNoContent, nil)
}

type refreshBookmarkRequest struct {
	BookmarkID        uuid.UUID `param:"bookmark-id" validate:"required,uuid4"`
	Fetcher           string    `json:"fetcher" validate:"omitempty,oneof=http jina browser"`
	RegenerateSummary bool      `json:"regenerate_summary"`
}

// refreshBookmark handles POST /bookmarks/:bookmark-id/refresh
// @Summary Refresh Bookmark
// @Description Refreshes bookmark content and/or regenerates AI summary
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param bookmark-id path string true "Bookmark ID"
// @Param request body refreshBookmarkRequest true "Refresh options"
// @Success 200 {object} JSONResult{data=bookmarks.ContentDTO} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 404 {object} JSONResult{data=nil} "Not Found"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /bookmarks/{bookmark-id}/refresh [post]
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

	bookmark, err := h.service.Refresh(ctx, tx, req.BookmarkID, user.ID, fetcher.FecherType(req.Fetcher), req.RegenerateSummary)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, bookmark)
}
