package httpserver

import (
	"context"
	"net/http"
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

func newToolsHandler(service toolService) *toolsHandler {
	return &toolsHandler{
		service: service,
	}
}

func (api *toolsHandler) Register(g *echo.Group) {
	tools := g.Group("/tools")
	tools.GET("/web/reader", api.webReaderHandler)
	tools.POST("/web/reader", api.webReaderHandler)
	tools.GET("/web/search", api.webSearchHandler)
	tools.POST("/web/search", api.webSearchHandler)
	tools.GET("/web/summary", api.webSummaryHandler)
	tools.POST("/web/summary", api.webSummaryHandler)
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
// @success 200 {object} handlers.JSONResult{data=jinareader.Content} "Success"
// @Failure 400 {object} handlers.JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} handlers.JSONResult{data=nil} "Internal Server Error"
// @Router /web/reader [get]
// @Router /web/reader [post]
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
// @success 200 {object} handlers.JSONResult{data=jinasearcher.Content} "Success"
// @Failure 400 {object} handlers.JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} handlers.JSONResult{data=nil} "Internal Server Error"
// @Router /web/search [get]
// @Router /web/search [post]
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
// @success 200 {object} handlers.JSONResult{data=string} "Success"
// @Failure 400 {object} handlers.JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} handlers.JSONResult{data=nil} "Internal Server Error"
// @Router /web/summary [get]
// @Router /web/summary [post]
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
