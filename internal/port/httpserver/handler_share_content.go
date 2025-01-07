package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"recally/internal/core/bookmarks"
	"recally/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BookmarkShareService interface {
	GetSharedContent(ctx context.Context, tx db.DBTX, sharedID uuid.UUID) (*bookmarks.ContentDTO, error)
}

// bookmarkServiceImpl implements BookmarkService
type bookmarkShareHandler struct {
	service BookmarkShareService
}

func registerBookmarkShareHandlers(e *echo.Group, s *Service) {
	// no auth middleware
	h := &bookmarkShareHandler{service: bookmarks.NewService(s.llm)}
	e.GET("/shared/:token", h.getSharedBookmark)
}

type getSharedBookmarkRequest struct {
	Token uuid.UUID `param:"token" validate:"required,uuid"`
}

// getSharedBookmark handles GET /bookmarks/:bookmark-id/share
// @Summary Get Shared Bookmark
// @Description Gets sharing information for a bookmark
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param token path string true "Bookmark ID"
// @Success 200 {object} JSONResult{data=bookmarks.ContentDTO} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 404 {object} JSONResult{data=nil} "Not Found"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /bookmarks/{bookmark-id}/share [get]
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

	bookmark, err := h.service.GetSharedContent(ctx, tx, req.Token)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	if bookmark == nil {
		return ErrorResponse(c, http.StatusNotFound, fmt.Errorf("shared bookmark not found"))
	}

	return JsonResponse(c, http.StatusOK, bookmark)
}
