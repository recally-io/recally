package httpserver

import (
	"context"
	"net/http"
	"vibrain/internal/core/assistants"
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
	UpdateThread(ctx context.Context, tx db.DBTX, thread *assistants.ThreadDTO) (*assistants.ThreadDTO, error)
	GetThread(ctx context.Context, tx db.DBTX, id uuid.UUID) (*assistants.ThreadDTO, error)
	RunThread(ctx context.Context, tx db.DBTX, thread *assistants.ThreadDTO) (*assistants.ThreadMessageDTO, error)

	ListThreadMessages(ctx context.Context, tx db.DBTX, threadID uuid.UUID) ([]assistants.ThreadMessageDTO, error)
	CreateThreadMessage(ctx context.Context, tx db.DBTX, threadId uuid.UUID, message *assistants.ThreadMessageDTO) (*assistants.ThreadMessageDTO, error)
	AddThreadMessage(ctx context.Context, tx db.DBTX, thread *assistants.ThreadDTO, role, text string) (*assistants.ThreadMessageDTO, error)

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
	g.PUT("/:assistant-id/threads/:thread-id", h.updateThread)

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

type getAssistantRequest struct {
	AssistantId uuid.UUID `param:"assistant-id" validate:"required,uuid4"`
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
	req := new(getAssistantRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}
	tx, _, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	assistant, err := h.service.GetAssistant(ctx, tx, req.AssistantId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, assistant)
}

type createAssistantRequest struct {
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description,omitempty"`
	SystemPrompt string `json:"system_prompt,omitempty"`
	Model        string `json:"model,omitempty"`
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
	req := new(createAssistantRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}
	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	assistantDTO := assistants.AssistantDTO{
		UserId:       user.ID,
		Name:         req.Name,
		Description:  req.Description,
		SystemPrompt: req.SystemPrompt,
		Model:        req.Model,
	}

	assistant, err := h.service.CreateAssistant(ctx, tx, &assistantDTO)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusCreated, assistant)
}

type updateAssistantRequest struct {
	AssistantId  uuid.UUID `param:"assistant-id" validate:"required,uuid4"`
	Name         string    `json:"name" validate:"required"`
	Description  string    `json:"description,omitempty"`
	SystemPrompt string    `json:"system_prompt,omitempty"`
	Model        string    `json:"model,omitempty"`
}

// updateAssistant is a handler function that updates an existing assistant.
// It retrieves the assistant ID from the request parameters and uses it to fetch the assistant.
// If the assistant ID is not found in the parameters, it returns an error with status code 400 (Bad Request).
// If there is an error while fetching the assistant, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it updates the assistant with the new values and returns a JSON response with status code 200 (OK) and the updated assistant.
//
// @Summary Update Assistant
// @Description Updates an existing assistant
// @Tags Assistants
// @Accept json
// @Produce json
// @Param assistant-id path string true "Assistant ID"
// @Param assistant body updateAssistantRequest true "Assistant"
// @Success 200 {object} JSONResult{data=assistants.AssistantDTO} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/{assistant-id} [put]
func (h *assistantHandler) updateAssistant(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(updateAssistantRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, _, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	assistant, err := h.service.GetAssistant(ctx, tx, req.AssistantId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if req.Name != "" {
		assistant.Name = req.Name
	}
	if req.Description != "" {
		assistant.Description = req.Description
	}
	if req.SystemPrompt != "" {
		assistant.SystemPrompt = req.SystemPrompt
	}
	if req.Model != "" {
		assistant.Model = req.Model
	}

	assistant, err = h.service.UpdateAssistant(ctx, tx, assistant)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusCreated, assistant)
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
