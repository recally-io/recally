package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

type createBookmarkShareRequest struct {
	BookmarkID uuid.UUID `param:"bookmark-id" validate:"required,uuid4"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// createBookmarkShare handles POST /bookmarks/:bookmark-id/share
//
//	@Summary		Share Bookmark
//	@Description	Creates a shareable link for a bookmark
//	@Tags			Bookmarks
//	@Accept			json
//	@Produce		json
//	@Param			bookmark-id	path		string										true	"Bookmark ID"
//	@Param			request		body		createBookmarkShareRequest					true	"Share options"
//	@Success		200			{object}	JSONResult{data=bookmarks.BookmarkShareDTO}	"Success"
//	@Failure		400			{object}	JSONResult{data=nil}						"Bad Request"
//	@Failure		401			{object}	JSONResult{data=nil}						"Unauthorized"
//	@Failure		404			{object}	JSONResult{data=nil}						"Not Found"
//	@Failure		500			{object}	JSONResult{data=nil}						"Internal Server Error"
//	@Router			/bookmarks/{bookmark-id}/share [post]
func (h *bookmarksHandler) createBookmarkShare(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(createBookmarkShareRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	shared, err := h.service.CreateBookmarkShare(ctx, tx, user.ID, req.BookmarkID, req.ExpiresAt)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, shared)
}

// updateBookmarkShare handles PUT /bookmarks/:bookmark-id/share
//
//	@Summary		Update Shared Bookmark
//	@Description	Updates sharing settings for a bookmark
//	@Tags			Bookmarks
//	@Accept			json
//	@Produce		json
//	@Param			bookmark-id	path		string										true	"Bookmark ID"
//	@Param			request		body		createBookmarkShareRequest					true	"Update options"
//	@Success		200			{object}	JSONResult{data=bookmarks.BookmarkShareDTO}	"Success"
//	@Failure		400			{object}	JSONResult{data=nil}						"Bad Request"
//	@Failure		404			{object}	JSONResult{data=nil}						"Not Found"
//	@Failure		500			{object}	JSONResult{data=nil}						"Internal Server Error"
//	@Router			/bookmarks/{bookmark-id}/share [put]
func (h *bookmarksHandler) updateBookmarkShare(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(createBookmarkShareRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	content, err := h.service.UpdateBookmarkShare(ctx, tx, user.ID, req.BookmarkID, req.ExpiresAt)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, content)
}

type getBookmarkShareRequest struct {
	BookmarkID uuid.UUID `param:"bookmark-id" validate:"required,uuid"`
}

// deleteBookmarkShare handles DELETE /bookmarks/:bookmark-id/share
//
//	@Summary		Delete Shared Bookmark
//	@Description	Revokes sharing access for a bookmark
//	@Tags			Bookmarks
//	@Accept			json
//	@Produce		json
//	@Param			bookmark-id	path		string					true	"Bookmark ID"
//	@Success		204			{object}	JSONResult{data=nil}	"No Content"
//	@Failure		400			{object}	JSONResult{data=nil}	"Bad Request"
//	@Failure		500			{object}	JSONResult{data=nil}	"Internal Server Error"
//	@Router			/bookmarks/{bookmark-id}/share [delete]
func (h *bookmarksHandler) deleteBookmarkShare(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(getBookmarkShareRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if err := h.service.DeleteBookmarkShare(ctx, tx, user.ID, req.BookmarkID); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusNoContent, nil)
}

// getBookmarkShare handles GET /bookmarks/:bookmark-id/share
//
//	@Summary		Get Shared Bookmark
//	@Description	Gets sharing information for a bookmark
//	@Tags			Bookmarks
//	@Accept			json
//	@Produce		json
//	@Param			bookmark-id	path		string										true	"Bookmark ID"
//	@Success		200			{object}	JSONResult{data=bookmarks.BookmarkShareDTO}	"Success"
//	@Failure		400			{object}	JSONResult{data=nil}						"Bad Request"
//	@Failure		404			{object}	JSONResult{data=nil}						"Not Found"
//	@Failure		500			{object}	JSONResult{data=nil}						"Internal Server Error"
//	@Router			/bookmarks/{bookmark-id}/share [get]
func (h *bookmarksHandler) getBookmarkShare(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(getBookmarkShareRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}
	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	cs, err := h.service.GetBookmarkShare(ctx, tx, user.ID, req.BookmarkID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	if cs == nil {
		return ErrorResponse(c, http.StatusNotFound, fmt.Errorf("shared bookmark not found"))
	}

	return JsonResponse(c, http.StatusOK, cs)
}
