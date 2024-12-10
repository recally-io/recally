package httpserver

import (
	"context"
	"net/http"
	"vibrain/internal/core/workers"
	"vibrain/internal/pkg/logger"
	"vibrain/internal/pkg/tools/jinareader"
	"vibrain/internal/pkg/tools/jinasearcher"

	"github.com/labstack/echo/v4"
)

type toolService interface {
	WebReader(ctx context.Context, url string) (*jinareader.Content, error)
	WebSearcher(ctx context.Context, query string) ([]*jinasearcher.Content, error)
	WebSummary(ctx context.Context, url string) (string, error)
}

type toolsHandler struct {
	service toolService
}

func registerToolsHandlers(e *echo.Group, s *Service) {
	h := &toolsHandler{
		service: workers.New(s.cache),
	}
	g := e.Group("/tools", authUserMiddleware())
	g.GET("/web/reader", h.webReaderHandler)
	g.POST("/web/reader", h.webReaderHandler)
	g.GET("/web/search", h.webSearchHandler)
	g.POST("/web/search", h.webSearchHandler)
	g.GET("/web/summary", h.webSummaryHandler)
	g.POST("/web/summary", h.webSummaryHandler)
}

type WebReaderRequest struct {
	URL string `form:"url"`
}

// @Summary Read web content
// @Description Read the content of a web page
// @Tags tools
// @Accept json
// @Produce json
// @Param url query string true "URL of the web page"
// @success 200 {object} JSONResult{data=jinareader.Content} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /tools/web/reader [get]
// @Router /tools/web/reader [post]
func (h *toolsHandler) webReaderHandler(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(WebReaderRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	content, err := h.service.WebReader(ctx, req.URL)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, content)
}

type WebSearchRequest struct {
	Query string `form:"query"`
}

// @Summary Search web content
// @Description Search the content of a web page
// @Tags tools
// @Accept json
// @Produce json
// @Param query query string true "Query string"
// @success 200 {object} JSONResult{data=jinasearcher.Content} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /tools/web/search [get]
// @Router /tools/web/search [post]
func (h *toolsHandler) webSearchHandler(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(WebSearchRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}
	logger.FromContext(ctx).Info("WebSearchHandler", "query", req.Query)
	content, err := h.service.WebSearcher(ctx, req.Query)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, content)
}

type WebSummaryRequest struct {
	URL string `form:"url"`
}

// @Summary Get web summary
// @Description Get the summary of a web page
// @Tags tools
// @Accept json
// @Produce json
// @Produce plain
// @Param url query string true "URL of the web page"
// @success 200 {object} JSONResult{data=string} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /tools/web/summary [get]
// @Router /tools/web/summary [post]
func (h *toolsHandler) webSummaryHandler(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(WebSummaryRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	content, err := h.service.WebSummary(ctx, req.URL)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, content)
}
