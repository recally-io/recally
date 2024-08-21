package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"vibrain/internal/core/assistants"
	"vibrain/internal/pkg/contexts"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type assistantService interface {
	ListAssistants(ctx context.Context, tx db.DBTX, userId uuid.UUID) ([]assistants.AssistantDTO, error)
	CreateAssistant(ctx context.Context, tx db.DBTX, assistant *assistants.AssistantDTO) (*assistants.AssistantDTO, error)
	UpdateAssistant(ctx context.Context, tx db.DBTX, assistant *assistants.AssistantDTO) (*assistants.AssistantDTO, error)
	GetAssistant(ctx context.Context, tx db.DBTX, id uuid.UUID) (*assistants.AssistantDTO, error)
	ListThreads(ctx context.Context, tx db.DBTX, assistantID uuid.UUID) ([]assistants.ThreadDTO, error)
	CreateThread(ctx context.Context, tx db.DBTX, thread *assistants.ThreadDTO) (*assistants.ThreadDTO, error)
	GetThread(ctx context.Context, tx db.DBTX, id uuid.UUID) (*assistants.ThreadDTO, error)
	ListThreadMessages(ctx context.Context, tx db.DBTX, threadID uuid.UUID) ([]assistants.ThreadMessageDTO, error)
	CreateThreadMessage(ctx context.Context, tx db.DBTX, threadId uuid.UUID, message *assistants.ThreadMessageDTO) (*assistants.ThreadMessageDTO, error)
	AddThreadMessage(ctx context.Context, tx db.DBTX, thread *assistants.ThreadDTO, role, text string) (*assistants.ThreadMessageDTO, error)
	RunThread(ctx context.Context, tx db.DBTX, thread *assistants.ThreadDTO) (*assistants.ThreadMessageDTO, error)
	ListModels(ctx context.Context) ([]string, error)
}

type assistantHandler struct {
	service assistantService
}

func registerAssistantHandlers(e *echo.Group, s *Service) {
	h := &assistantHandler{service: assistants.NewService(s.llm)}
	g := e.Group("/assistants")
	g.GET("", h.listAssistants)
	g.POST("", h.createAssistant)
	g.GET("/:assistant-id", h.getAssistant)
	g.PUT("/:assistant-id", h.updateAssistant)

	g.GET("/:assistant-id/threads", h.listThreads)
	g.GET("/:assistant-id/threads/:thread-id", h.getThread)
	g.POST("/:assistant-id/threads", h.createThread)

	g.GET("/:assistant-id/threads/:thread-id/messages", h.listThreadMessages)
	g.POST("/:assistant-id/threads/:thread-id/messages", h.createThreadMessage)

	g.GET("/models", h.listModels)
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
// @success 200 {object} JSONResult{data=[]assistants.AssistantDTO} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants [get]
func (h *assistantHandler) listAssistants(c echo.Context) error {
	ctx := c.Request().Context()
	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	assistants, err := h.service.ListAssistants(ctx, tx, user.ID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, assistants)
}

// GetAssistant is a handler function that retrieves an assistant by ID.
// It retrieves the assistant ID from the request parameters and uses it to fetch the assistant.
// If the assistant ID is not found in the parameters, it returns an error with status code 400 (Bad Request).
// If there is an error while fetching the assistant, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it returns a JSON response with status code 200 (OK) and the assistant.

// @Summary Get Assistant
// @Description Retrieves an assistant by ID
// @Tags Assistants
// @Accept json
// @Produce json
// @Param assistant-id path string true "Assistant ID"
// @Success 200 {object} JSONResult{data=assistants.AssistantDTO} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal
// @Router /assistants/{assistant-id} [get]
func (h *assistantHandler) getAssistant(c echo.Context) error {
	ctx := c.Request().Context()
	assistantId, err := uuid.Parse(c.Param("assistant-id"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid assistant-id: %s", c.Param("assistant-id")))
	}
	tx, _, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	assistant, err := h.service.GetAssistant(ctx, tx, assistantId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, assistant)
}

type createAssistantRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	SystemPrompt string `json:"system_prompt"`
	Model        string `json:"model"`
}

// createAssistant is a handler function that creates a new assistant.
// It retrieves the user ID from the request context and uses it to create the assistant.
// If the user ID is not found in the context, it returns an error with status code 401 (Unauthorized).
// If there is an error while creating the assistant, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it returns a JSON response with status code 201 (Created) and the created assistant.
//
// @Summary Create Assistant
// @Description Creates a new assistant
// @Tags Assistants
// @Accept json
// @Produce json
// @Param assistant body createAssistantRequest true "Assistant"
// @Success 201 {object} JSONResult{data=assistants.AssistantDTO} "Created"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants [post]
func (h *assistantHandler) createAssistant(c echo.Context) error {
	ctx := c.Request().Context()
	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	var req createAssistantRequest
	if err := c.Bind(&req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	assistantDTO := assistants.AssistantDTO{
		Name:         req.Name,
		Description:  req.Description,
		SystemPrompt: req.SystemPrompt,
		UserId:       user.ID,
		Model:        req.Model,
	}

	assistant, err := h.service.CreateAssistant(ctx, tx, &assistantDTO)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusCreated, assistant)
}

func (h *assistantHandler) updateAssistant(c echo.Context) error {
	assistantId, err := uuid.Parse(c.Param("assistant-id"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid assistant-id: %s", c.Param("assistant-id")))
	}
	ctx := c.Request().Context()
	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	var req createAssistantRequest
	if err := c.Bind(&req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	assistantDTO := assistants.AssistantDTO{
		Id:           assistantId,
		Name:         req.Name,
		Description:  req.Description,
		SystemPrompt: req.SystemPrompt,
		UserId:       user.ID,
		Model:        req.Model,
	}

	assistant, err := h.service.UpdateAssistant(ctx, tx, &assistantDTO)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusCreated, assistant)
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
// @Param assistant-id path string true "Assistant ID"
// @success 200 {object} JSONResult{data=[]assistants.ThreadDTO} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/{assistant-id}/threads [get]
func (h *assistantHandler) listThreads(c echo.Context) error {
	ctx := c.Request().Context()
	assistantId, err := uuid.Parse(c.Param("assistant-id"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid assistant-id: %s", c.Param("assistant-id")))
	}
	tx, ok := contexts.Get[db.DBTX](ctx, contexts.ContextKeyTx)
	if !ok {
		return ErrorResponse(c, http.StatusInternalServerError, errors.New("missing transaction"))
	}
	threads, err := h.service.ListThreads(ctx, tx, assistantId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, threads)
}

type createThreadRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Model        string `json:"model"`
	SystemPrompt string `json:"system_prompt"`
}

// createThread is a handler function that creates a new thread for an assistant.
// It retrieves the assistant ID from the request parameters and uses it to create the thread.
// If the assistant ID is not found in the parameters, it returns an error with status code 400 (Bad Request).
// If there is an error while creating the thread, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it returns a JSON response with status code 201 (Created) and the created thread.
//
// @Summary Create Thread
// @Description Creates a new thread under an assistant
// @Tags Assistants
// @Accept json
// @Produce json
// @Param assistant-id path string true "Assistant ID"
// @Param thread body assistants.ThreadDTO true "Thread"
// @Success 201 {object} JSONResult{data=assistants.ThreadDTO} "Created"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/{assistant-id}/threads [post]
func (h *assistantHandler) createThread(c echo.Context) error {
	ctx := c.Request().Context()
	assistantId, err := uuid.Parse(c.Param("assistant-id"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid assistant-id: %s", c.Param("assistant-id")))
	}
	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	var req createThreadRequest
	if err := c.Bind(&req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	// Assign the AssistantID for the new thread
	threadDTO := assistants.ThreadDTO{
		Name:         req.Name,
		Description:  req.Description,
		Model:        req.Model,
		SystemPrompt: req.SystemPrompt,
		AssistantId:  assistantId,
		UserId:       user.ID,
	}

	thread, err := h.service.CreateThread(ctx, tx, &threadDTO)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusCreated, thread)
}

// getThread is a handler function that retrieves a thread by ID.
// It retrieves the thread ID from the request parameters and uses it to fetch the thread.
// If the thread ID is not found in the parameters, it returns an error with status code 400 (Bad Request).
// If there is an error while fetching the thread, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it returns a JSON response with status code 200 (OK) and the thread.

// @Summary Get Thread
// @Description Retrieves a thread by ID
// @Tags Assistants
// @Accept json
// @Produce json
// @Param thread-id path string true "Thread ID"
// @success 200 {object} JSONResult{data=assistants.ThreadDTO} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal
// @Router /assistants/{assistant-id}/threads/{thread-id} [get]
func (h *assistantHandler) getThread(c echo.Context) error {
	ctx := c.Request().Context()
	_, err := uuid.Parse(c.Param("assistant-id"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid assistant-id: %s", c.Param("assistant-id")))
	}
	threadId, err := uuid.Parse(c.Param("thread-id"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid thread-id: %s", c.Param("thread-id")))
	}
	tx, _, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	thread, err := h.service.GetThread(ctx, tx, threadId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, thread)
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
// @Param assistant-id path string true "Assistant ID"
// @Param thread-id path string true "Thread ID"
// @success 200 {object} JSONResult{data=[]assistants.ThreadMessageDTO} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/{assistant-id}/threads/{thread-id}/messages [get]
func (h *assistantHandler) listThreadMessages(c echo.Context) error {
	ctx := c.Request().Context()
	_, err := uuid.Parse(c.Param("assistant-id"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid assistant-id: %s", c.Param("assistant-id")))
	}
	threadId, err := uuid.Parse(c.Param("thread-id"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid thread-id: %s", c.Param("thread-id")))
	}
	tx, ok := contexts.Get[db.DBTX](ctx, contexts.ContextKeyTx)
	if !ok {
		return ErrorResponse(c, http.StatusInternalServerError, errors.New("missing transaction"))
	}
	messages, err := h.service.ListThreadMessages(ctx, tx, threadId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, messages)
}

type createThreadMessageRequest struct {
	Role  string `json:"role"`
	Text  string `json:"text"`
	Model string `json:"model"`
}

// createThreadMessage is a handler function that creates a new message for a thread.
// It retrieves the thread and assistant IDs from the request parameters and uses them to create the message.
// If the assistant ID or thread ID is not found in the parameters, it returns an error with status code 400 (Bad Request).
// If there is an error while creating the message, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it returns a JSON response with status code 201 (Created) and the created message.
//
// @Summary Create Thread Message
// @Description Creates a new message in a specified thread
// @Tags Assistants
// @Accept json
// @Produce json
// @Param assistant-id path string true "Assistant ID"
// @Param thread-id path string true "Thread ID"
// @Param message body createThreadMessageRequest true "Thread Message"
// @Success 201 {object} JSONResult{data=assistants.ThreadMessageDTO} "Created"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/{assistant-id}/threads/{thread-id}/messages [post]
func (h *assistantHandler) createThreadMessage(c echo.Context) error {
	ctx := c.Request().Context()
	_, err := uuid.Parse(c.Param("assistant-id"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid assistant-id: %s", c.Param("assistant-id")))
	}
	threadId, err := uuid.Parse(c.Param("thread-id"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid thread-id: %s", c.Param("thread-id")))
	}

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	thread, err := h.service.GetThread(ctx, tx, threadId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	var req createThreadMessageRequest
	if err := c.Bind(&req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	messageDTO := assistants.ThreadMessageDTO{
		Role:     req.Role,
		Text:     req.Text,
		ThreadID: thread.Id,
		UserID:   user.ID,
		Model:    thread.Model,
	}

	// Create Thread Message
	if _, err := h.service.CreateThreadMessage(ctx, tx, thread.Id, &messageDTO); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	resp, err := h.service.RunThread(ctx, tx, thread)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusCreated, resp)
}

// @Summary List Models
// @Description Lists available language models
// @Tags Assistants
// @Accept json
// @Produce json
// @Success 200 {object} JSONResult{data=[]string} "Success"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/models [get]
func (h *assistantHandler) listModels(c echo.Context) error {
	models, err := h.service.ListModels(c.Request().Context())
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, models)
}
