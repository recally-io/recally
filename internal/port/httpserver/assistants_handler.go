package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"vibrain/internal/core/assistants"
	"vibrain/internal/pkg/contexts"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type assistantService interface {
	ListAssistants(ctx context.Context, tx db.DBTX, userId uuid.UUID) ([]assistants.AssistantDTO, error)
	CreateAssistant(ctx context.Context, tx db.DBTX, assistant *assistants.AssistantDTO) (*assistants.AssistantDTO, error)
	GetAssistant(ctx context.Context, tx db.DBTX, id uuid.UUID) (*assistants.AssistantDTO, error)
	ListThreads(ctx context.Context, tx db.DBTX, assistantID uuid.UUID) ([]assistants.ThreadDTO, error)
	ListThreadMessages(ctx context.Context, tx db.DBTX, threadID uuid.UUID) ([]assistants.ThreadMessageDTO, error)
}

type assistantHandler struct {
	service assistantService
}

func registerAssistantHandlers(e *echo.Group, s *Service) {
	h := &assistantHandler{service: assistants.NewService(s.llm)}
	g := e.Group("/assistants")
	g.GET("/", h.listAssistants)
	g.GET("/:assistant-id/threads", h.listThreads)
	g.GET("/:assistant-id/threads/:thread-id/messages", h.listThreadMessages)
}

// listAssistants is a handler function that lists the assistants for a user.
// It retrieves the user ID from the request context and uses it to fetch the assistants.
// If the user ID is not found in the context, it returns an error with status code 401 (Unauthorized).
// If there is an error while fetching the assistants, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it returns a JSON response with status code 200 (OK) and the list of assistants.
//
// @Summary List Assistants
// @Description Lists the assistants for a user
// @Tags Assistants
// @Accept json
// @Produce json
// @success 200 {object} JSONResult{data=[]assistants.Assistant} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants [get]
func (h *assistantHandler) listAssistants(c echo.Context) error {
	ctx := c.Request().Context()
	// userId
	userId, ok := contexts.Get[string](ctx, contexts.ContextKeyUserID)
	logger.FromContext(ctx).Info("list assistants", "user_id", userId)
	if !ok {
		return ErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("user not found"))
	}
	tx, ok := contexts.Get[db.DBTX](ctx, contexts.ContextKeyTx)
	if !ok {
		return ErrorResponse(c, http.StatusInternalServerError, errors.New("missing transaction"))
	}
	assistants, err := h.service.ListAssistants(ctx, tx, uuid.MustParse(userId))
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, assistants)
}

// listThreads is a handler function that lists the threads for an assistant.
// It retrieves the assistant ID from the request parameters and uses it to fetch the threads.
// If the assistant ID is not found in the parameters, it returns an error with status code 400 (Bad Request).
// If there is an error while fetching the threads, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it returns a JSON response with status code 200 (OK) and the list of threads.
//
// @Summary List Threads
// @Description Lists the threads for an assistant
// @Tags Assistants
// @Accept json
// @Produce json
// @PathParam assistant-id path string true "Assistant ID"
// @success 200 {object} JSONResult{data=[]assistants.Thread} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/{assistant-id}/threads [get]
func (h *assistantHandler) listThreads(c echo.Context) error {
	ctx := c.Request().Context()
	assistantId := c.Param("assistant-id")
	if assistantId == "" {
		return ErrorResponse(c, http.StatusBadRequest, errors.New("missing assistant-id"))
	}
	tx, ok := contexts.Get[db.DBTX](ctx, contexts.ContextKeyTx)
	if !ok {
		return ErrorResponse(c, http.StatusInternalServerError, errors.New("missing transaction"))
	}
	threads, err := h.service.ListThreads(ctx, tx, uuid.MustParse(assistantId))
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, threads)
}

// listThreadMessages is a handler function that lists the messages for a thread.
// It retrieves the thread ID from the request parameters and uses it to fetch the messages.
// If the thread ID is not found in the parameters, it returns an error with status code 400 (Bad Request).
// If there is an error while fetching the messages, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it returns a JSON response with status code 200 (OK) and the list of messages.
//
// @Summary List Thread Messages
// @Description Lists the messages for a thread
// @Tags Assistants
// @Accept json
// @Produce json
// @PathParam assistant-id path string true "Assistant ID"
// @PathParam thread-id path string true "Thread ID"
// @success 200 {object} JSONResult{data=[]assistants.ThreadMessage} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/{assistant-id}/threads/{thread-id}/messages [get]
func (h *assistantHandler) listThreadMessages(c echo.Context) error {
	ctx := c.Request().Context()
	assistantId := c.Param("assistant-id")
	threadId := c.Param("thread-id")
	if assistantId == "" || threadId == "" {
		return ErrorResponse(c, http.StatusBadRequest, errors.New("missing assistant-id or thread-id"))
	}
	tx, ok := contexts.Get[db.DBTX](ctx, contexts.ContextKeyTx)
	if !ok {
		return ErrorResponse(c, http.StatusInternalServerError, errors.New("missing transaction"))
	}
	messages, err := h.service.ListThreadMessages(ctx, tx, uuid.MustParse(threadId))
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, messages)
}
