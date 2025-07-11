package httpserver

import (
	"context"
	"net/http"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/tools"

	"github.com/labstack/echo/v4"
)

type llmService interface {
	ListModels(ctx context.Context) ([]llms.Model, error)
	ListTools(ctx context.Context) ([]tools.BaseTool, error)
}

type llmHandler struct {
	service llmService
}

func registerLLMHandlers(e *echo.Group, s *Service) {
	h := &llmHandler{service: s.llm}
	g := e.Group("/llm", authUserMiddleware())
	g.GET("/models", h.listModels)
	g.GET("/tools", h.listTools)
}

// @Router	/llm/models [get].
func (h *llmHandler) listModels(c echo.Context) error {
	models, err := h.service.ListModels(c.Request().Context())
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, models)
}

// @Router	/llm/tools [get].
func (h *llmHandler) listTools(c echo.Context) error {
	tools, err := h.service.ListTools(c.Request().Context())
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, tools)
}
