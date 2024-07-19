package handlers

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type WebReaderRequest struct {
	URL string `form:"url"`
}

func (h *Handler) WebReaderHandler(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(WebReaderRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	content, err := h.worker.WebReader(ctx, req.URL)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, content)
}

type WebSearchRequest struct {
	Query string `form:"query"`
}

func (h *Handler) WebSearchHandler(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(WebSearchRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}
	slog.InfoContext(c.Request().Context(), "WebSearchHandler", "query", req.Query)
	content, err := h.worker.WebSearcher(ctx, req.Query)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, content)
}

type WebSummaryRequest struct {
	URL string `form:"url"`
}

func (h *Handler) WebSummaryHandler(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(WebSummaryRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	content, err := h.worker.WebSummary(ctx, req.URL)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, content)
}
