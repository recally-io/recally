package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"recally/internal/core/bookmarks"
	"recally/internal/core/files"
	"recally/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BookmarkShareService interface {
	GetBookmarkShareContent(ctx context.Context, tx db.DBTX, sharedID uuid.UUID) (*bookmarks.BookmarkContentDTO, error)
}

// bookmarkServiceImpl implements BookmarkService.
type bookmarkShareHandler struct {
	service     BookmarkShareService
	fileService *files.Service
}

func registerBookmarkShareHandlers(e *echo.Group, s *Service) {
	// no auth middleware
	h := &bookmarkShareHandler{service: bookmarks.NewService(s.llm), fileService: files.NewService(s.s3)}
	g := e.Group("/shared")
	g.GET("/files/:key", h.redirectToFile)
	g.HEAD("/files/:key", h.getFileMetadata)
	g.GET("/:token", h.getSharedBookmark)
}

type getSharedBookmarkRequest struct {
	Token uuid.UUID `param:"token" validate:"required,uuid"`
}

type sharedFileRequest struct {
	Key string `param:"key" validate:"required"`
}

// getSharedBookmark handles GET /bookmarks/:bookmark-id/share
//
//	@Summary		Get Shared Bookmark
//	@Description	Gets sharing information for a bookmark
//	@Tags			Bookmarks
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string											true	"Bookmark ID"
//	@Success		200		{object}	JSONResult{data=bookmarks.BookmarkContentDTO}	"Success"
//	@Failure		400		{object}	JSONResult{data=nil}							"Bad Request"
//	@Failure		404		{object}	JSONResult{data=nil}							"Not Found"
//	@Failure		500		{object}	JSONResult{data=nil}							"Internal Server Error"
//	@Router			/bookmarks/{bookmark-id}/share [get]
func (h *bookmarkShareHandler) getSharedBookmark(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(getSharedBookmarkRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, err := loadTx(ctx)
	if err != nil {
		return errors.New("tx not found")
	}

	bookmark, err := h.service.GetBookmarkShareContent(ctx, tx, req.Token)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if bookmark == nil {
		return ErrorResponse(c, http.StatusNotFound, fmt.Errorf("shared bookmark not found"))
	}

	// remove sensitive information
	bookmark.ID = uuid.Nil
	bookmark.UserID = uuid.Nil

	return JsonResponse(c, http.StatusOK, bookmark)
}

// @Router			/files/{id} [get].
func (h *bookmarkShareHandler) redirectToFile(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(sharedFileRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	// Get presigned URL with 1 hour expiration
	presignedURL, err := h.fileService.GetPresignedGetObjectURL(ctx, tx, user.ID, req.Key, time.Hour, nil)
	if err != nil {
		return ErrorResponse(c, http.StatusNotFound, fmt.Errorf("file not found"))
	}

	// Redirect to the presigned URL
	return c.Redirect(http.StatusFound, presignedURL)
}

// @Router			/shared/files/{key} [head].
func (h *bookmarkShareHandler) getFileMetadata(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(sharedFileRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	// Get presigned URL with 1 hour expiration for HEAD request
	presignedURL, err := h.fileService.GetPresignedHeadObjectURL(ctx, tx, user.ID, req.Key, time.Hour, nil)
	if err != nil {
		return ErrorResponse(c, http.StatusNotFound, fmt.Errorf("file not found"))
	}

	// perform HEAD request
	resp, err := http.Head(presignedURL)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Copy relevant headers from S3 response to our response
	for key, values := range resp.Header {
		for _, value := range values {
			c.Response().Header().Add(key, value)
		}
	}

	return c.NoContent(resp.StatusCode)
}
