package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"vibrain/internal/pkg/contexts"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ListAssistants is a handler function that lists the assistants for a user.
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
// @success 200 {object} handlers.JSONResult{data=[]assistants.Assistant} "Success"
// @Failure 400 {object} handlers.JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} handlers.JSONResult{data=nil} "Internal Server Error"
// @Router /assistants [get]
func (h *Handler) ListAssistants(c echo.Context) error {
	ctx := c.Request().Context()
	// userId
	userId, ok := contexts.Get[string](ctx, contexts.ContextKeyUserID)
	if !ok {
		return ErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("user not found"))
	}

	assistants, err := h.assistant.ListAssistants(ctx, h.Pool, uuid.MustParse(userId))
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, assistants)
}

// ListThreads is a handler function that lists the threads for an assistant.
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
// @success 200 {object} handlers.JSONResult{data=[]assistants.Thread} "Success"
// @Failure 400 {object} handlers.JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} handlers.JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/{assistant-id}/threads [get]
func (h *Handler) ListThreads(c echo.Context) error {
	ctx := c.Request().Context()
	assistantId := c.Param("assistant-id")
	if assistantId == "" {
		return ErrorResponse(c, http.StatusBadRequest, errors.New("missing assistant-id"))
	}
	threads, err := h.assistant.ListThreads(ctx, h.Pool, uuid.MustParse(assistantId))
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, threads)
}

// ListThreadMessages is a handler function that lists the messages for a thread.
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
// @success 200 {object} handlers.JSONResult{data=[]assistants.ThreadMessage} "Success"
// @Failure 400 {object} handlers.JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} handlers.JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/{assistant-id}/threads/{thread-id}/messages [get]
func (h *Handler) ListThreadMessages(c echo.Context) error {
	ctx := c.Request().Context()
	assistantId := c.Param("assistant-id")
	threadId := c.Param("thread-id")
	if assistantId == "" || threadId == "" {
		return ErrorResponse(c, http.StatusBadRequest, errors.New("missing assistant-id or thread-id"))
	}
	messages, err := h.assistant.ListThreadMessages(ctx, h.Pool, uuid.MustParse(threadId))
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, messages)
}
