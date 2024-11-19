package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"vibrain/internal/core/bookmarks"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

// BookmarkService defines operations for managing bookmarks
type BookmarkService interface {
	Create(ctx context.Context, tx db.DBTX, dto *bookmarks.BookmarkDTO) (*bookmarks.BookmarkDTO, error)
	Get(ctx context.Context, tx db.DBTX, id, userID uuid.UUID) (*bookmarks.BookmarkDTO, error)
	List(ctx context.Context, tx db.DBTX, userID uuid.UUID, limit, offset int32) ([]*bookmarks.BookmarkDTO, error)
	Update(ctx context.Context, tx db.DBTX, id, userID uuid.UUID, dto *bookmarks.BookmarkDTO) (*bookmarks.BookmarkDTO, error)
	Delete(ctx context.Context, tx db.DBTX, id, userID uuid.UUID) error
	DeleteUserBookmarks(ctx context.Context, tx db.DBTX, userID uuid.UUID) error
	Refresh(ctx context.Context, tx db.DBTX, id, userID uuid.UUID, fetcher string, regenerateSummary bool) (*bookmarks.BookmarkDTO, error)
}

// bookmarkServiceImpl implements BookmarkService
type bookmarksHandler struct {
	service BookmarkService
}

func registerBookmarkHandlers(e *echo.Group, s *Service) {
	h := &bookmarksHandler{service: bookmarks.NewService(s.llm)}
	g := e.Group("/bookmarks")
	g.GET("", h.listBookmarks)
	g.POST("", h.createBookmark)
	g.DELETE("/", h.deleteUserBookmarks)
	g.GET("/:bookmark-id", h.getBookmark)
	g.PUT("/:bookmark-id", h.updateBookmark)
	g.DELETE("/:bookmark-id", h.deleteBookmark)
	g.POST("/:bookmark-id/refresh", h.refreshBookmark)
}

type listBookmarksRequest struct {
	Limit  int32 `query:"limit" validate:"min=1,max=100"`
	Offset int32 `query:"offset" validate:"min=0"`
}

// listBookmarks handles GET /bookmarks
// @Summary List Bookmarks
// @Description Lists bookmarks for a user with pagination
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page" default(10)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} JSONResult{data=[]bookmarks.BookmarkDTO} "Success"
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

	bookmarks, err := h.service.List(ctx, tx, user.ID, req.Limit, req.Offset)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, bookmarks)
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
// @Success 201 {object} JSONResult{data=bookmarks.BookmarkDTO} "Created"
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

	bookmark := &bookmarks.BookmarkDTO{
		UserID:   user.ID,
		URL:      req.URL,
		Metadata: req.Metadata,
	}

	created, err := h.service.Create(ctx, tx, bookmark)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
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
// @Success 200 {object} JSONResult{data=bookmarks.BookmarkDTO} "Success"
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
	BookmarkID uuid.UUID `param:"bookmark-id" validate:"required,uuid4"`
	Summary    string    `json:"summary"`
	Content    string    `json:"content"`
	HTML       string    `json:"html"`
	// Metadata bookmarks.Metadata `json:"metadata"`
}

// updateBookmark handles PUT /bookmarks/:bookmark-id
// @Summary Update Bookmark
// @Description Updates an existing bookmark
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param bookmark-id path string true "Bookmark ID"
// @Param bookmark body updateBookmarkRequest true "Updated bookmark data"
// @Success 200 {object} JSONResult{data=bookmarks.BookmarkDTO} "Success"
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

	if req.Summary == "" && req.Content == "" && req.HTML == "" {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("no update fields provided"))
	}

	bookmark := &bookmarks.BookmarkDTO{
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
// @Success 200 {object} JSONResult{data=bookmarks.BookmarkDTO} "Success"
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

	bookmark, err := h.service.Refresh(ctx, tx, req.BookmarkID, user.ID, req.Fetcher, req.RegenerateSummary)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, bookmark)
}
