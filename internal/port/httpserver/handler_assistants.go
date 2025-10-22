package httpserver

import (
	"context"
	"net/http"

	"recally/internal/core/assistants"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/rag/document"
	"recally/internal/pkg/tools"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type assistantService interface {
	ListAssistants(ctx context.Context, tx db.DBTX, userId uuid.UUID) ([]assistants.AssistantDTO, error)
	CreateAssistant(ctx context.Context, tx db.DBTX, assistant *assistants.AssistantDTO) (*assistants.AssistantDTO, error)
	UpdateAssistant(ctx context.Context, tx db.DBTX, assistant *assistants.AssistantDTO) (*assistants.AssistantDTO, error)
	GetAssistant(ctx context.Context, tx db.DBTX, id uuid.UUID) (*assistants.AssistantDTO, error)
	DeleteAssistant(ctx context.Context, tx db.DBTX, id uuid.UUID) error

	ListThreads(ctx context.Context, tx db.DBTX, assistantID uuid.UUID) ([]assistants.ThreadDTO, error)
	CreateThread(ctx context.Context, tx db.DBTX, thread *assistants.ThreadDTO) (*assistants.ThreadDTO, error)
	UpdateThread(ctx context.Context, tx db.DBTX, thread *assistants.ThreadDTO) (*assistants.ThreadDTO, error)
	GetThread(ctx context.Context, tx db.DBTX, id uuid.UUID) (*assistants.ThreadDTO, error)
	RunThread(ctx context.Context, tx db.DBTX, id uuid.UUID, streamingFunc func(*assistants.MessageDTO, error))
	DeleteThread(ctx context.Context, tx db.DBTX, id uuid.UUID) error
	GenerateThreadTitle(ctx context.Context, tx db.DBTX, id uuid.UUID) (string, error)

	ListThreadMessages(ctx context.Context, tx db.DBTX, threadID uuid.UUID) ([]assistants.MessageDTO, error)
	CreateThreadMessage(ctx context.Context, tx db.DBTX, threadId uuid.UUID, message *assistants.MessageDTO) (*assistants.MessageDTO, error)
	GetThreadMessage(ctx context.Context, tx db.DBTX, id uuid.UUID) (*assistants.MessageDTO, error)
	AddThreadMessage(ctx context.Context, tx db.DBTX, thread *assistants.ThreadDTO, role, text string, metadata assistants.MessageMetadata) (*assistants.MessageDTO, error)
	DeleteThreadMessage(ctx context.Context, tx db.DBTX, id uuid.UUID) error

	CreateAttachment(ctx context.Context, tx db.DBTX, attachment *assistants.AttachmentDTO, docs []document.Document) (*assistants.AttachmentDTO, error)
	ListAttachmentsByAssistant(ctx context.Context, tx db.DBTX, assistantId uuid.UUID) ([]assistants.AttachmentDTO, error)
	ListAttachmentsByThread(ctx context.Context, tx db.DBTX, threadId uuid.UUID) ([]assistants.AttachmentDTO, error)

	ListModels(ctx context.Context) ([]llms.Model, error)
	ListTools(ctx context.Context) ([]tools.BaseTool, error)
}

type assistantHandler struct {
	service assistantService
}

func registerAssistantHandlers(e *echo.Group, s *Service) {
	h := &assistantHandler{service: assistants.NewService(s.llm, s.queue)}
	g := e.Group("/assistants", authUserMiddleware())
	g.GET("", h.listAssistants)
	g.POST("", h.createAssistant)
	g.GET("/:assistant-id", h.getAssistant)
	g.PUT("/:assistant-id", h.updateAssistant)
	g.DELETE("/:assistant-id", h.deleteAssistant)

	g.GET("/:assistant-id/threads", h.listThreads)
	g.POST("/:assistant-id/threads", h.createThread)
	g.GET("/:assistant-id/threads/:thread-id", h.getThread)
	g.PUT("/:assistant-id/threads/:thread-id", h.updateThread)
	g.DELETE("/:assistant-id/threads/:thread-id", h.deleteThread)
	g.POST("/:assistant-id/threads/:thread-id/generate-title", h.generateThreadTitle)

	g.GET("/:assistant-id/threads/:thread-id/messages", h.listThreadMessages)
	g.POST("/:assistant-id/threads/:thread-id/messages", h.createThreadMessage)
	g.PUT("/:assistant-id/threads/:thread-id/messages/:message-id", h.updateThreadMessage)
	g.DELETE("/:assistant-id/threads/:thread-id/messages/:message-id", h.deleteThreadMessage)

	g.GET("/:assistant-id/attachments", h.listAttachmentsByAssistant)
	g.POST("/:assistant-id/attachments", h.uploadAssistantAttachment)
	g.GET("/:assistant-id/threads/:thread-id/attachments", h.listAttachmentsByThread)
	g.POST("/:assistant-id/threads/:thread-id/attachments", h.uploadThreadAttachment)

	g.GET("/models", h.listModels)
	g.GET("/tools", h.listTools)
}

// listAssistants is a handler function that lists the assistants for a user.
// It retrieves the user ID from the request context and uses it to fetch the assistants.
// If the user ID is not found in the context, it returns an error with status code 401 (Unauthorized).
// If there is an error while fetching the assistants, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it returns a JSON response with status code 200 (OK) and the list of assistants.
//
//	@Summary		List Assistants
//	@Description	Lists the assistants for a user
//	@Tags			Assistants
//	@Accept			json
//	@Produce		json
//	@success		200	{object}	JSONResult{data=[]assistants.AssistantDTO}	"Success"
//	@Failure		400	{object}	JSONResult{data=nil}						"Bad Request"
//	@Failure		401	{object}	JSONResult{data=nil}						"Unauthorized"
//	@Failure		500	{object}	JSONResult{data=nil}						"Internal Server Error"
//	@Router			/assistants [get]
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

// @Router			/assistants/{assistant-id} [get].
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
	Metadata     struct {
		Tools []string `json:"tools,omitempty"`
	} `json:"metadata,omitempty"`
}

// createAssistant is a handler function that creates a new assistant.
// It retrieves the user ID from the request context and uses it to create the assistant.
// If the user ID is not found in the context, it returns an error with status code 401 (Unauthorized).
// If there is an error while creating the assistant, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it returns a JSON response with status code 201 (Created) and the created assistant.
//
//	@Summary		Create Assistant
//	@Description	Creates a new assistant
//	@Tags			Assistants
//	@Accept			json
//	@Produce		json
//	@Param			assistant	body		createAssistantRequest						true	"Assistant"
//	@Success		201			{object}	JSONResult{data=assistants.AssistantDTO}	"Created"
//	@Failure		400			{object}	JSONResult{data=nil}						"Bad Request"
//	@Failure		401			{object}	JSONResult{data=nil}						"Unauthorized"
//	@Failure		401			{object}	JSONResult{data=nil}						"Unauthorized"
//	@Failure		500			{object}	JSONResult{data=nil}						"Internal Server Error"
//	@Router			/assistants [post]
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
		Metadata: assistants.AssistantMetadata{
			Tools: req.Metadata.Tools,
		},
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
	Metadata     struct {
		Tools []string `json:"tools,omitempty"`
	} `json:"metadata,omitempty"`
}

// updateAssistant is a handler function that updates an existing assistant.
// It retrieves the assistant ID from the request parameters and uses it to fetch the assistant.
// If the assistant ID is not found in the parameters, it returns an error with status code 400 (Bad Request).
// If there is an error while fetching the assistant, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it updates the assistant with the new values and returns a JSON response with status code 200 (OK) and the updated assistant.
//
//	@Summary		Update Assistant
//	@Description	Updates an existing assistant
//	@Tags			Assistants
//	@Accept			json
//	@Produce		json
//	@Param			assistant-id	path		string										true	"Assistant ID"
//	@Param			assistant		body		updateAssistantRequest						true	"Assistant"
//	@Success		200				{object}	JSONResult{data=assistants.AssistantDTO}	"Success"
//	@Failure		400				{object}	JSONResult{data=nil}						"Bad Request"
//	@Failure		401				{object}	JSONResult{data=nil}						"Unauthorized"
//	@Failure		500				{object}	JSONResult{data=nil}						"Internal Server Error"
//	@Router			/assistants/{assistant-id} [put]
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

	if req.Metadata.Tools != nil {
		assistant.Metadata.Tools = req.Metadata.Tools
	}

	assistant, err = h.service.UpdateAssistant(ctx, tx, assistant)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusCreated, assistant)
}

// deleteAssistant is a handler function that deletes an existing assistant.
//
//	@Summary		Delete Assistant
//	@Description	Deletes an existing assistant
//	@Tags			Assistants
//	@Accept			json
//	@Produce		json
//	@Param			assistant-id	path		string					true	"Assistant ID"
//	@Success		204				{object}	JSONResult{data=nil}	"No Content"
//	@Failure		400				{object}	JSONResult{data=nil}	"Bad Request"
//	@Failure		401				{object}	JSONResult{data=nil}	"Unauthorized"
//	@Failure		500				{object}	JSONResult{data=nil}	"Internal Server Error"
//	@Router			/assistants/{assistant-id} [delete]
func (h *assistantHandler) deleteAssistant(c echo.Context) error {
	ctx := c.Request().Context()
	assistantId := c.Param("assistant-id")

	uuid, err := uuid.Parse(assistantId)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	tx, _, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	err = h.service.DeleteAssistant(ctx, tx, uuid)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusNoContent, nil)
}

// @Router			/assistants/{assistant-id}/attachments [get].
func (h *assistantHandler) listAttachmentsByAssistant(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(getAssistantRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, _, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	attachments, err := h.service.ListAttachmentsByAssistant(ctx, tx, req.AssistantId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, attachments)
}

type uploadAssistantAttachmentRequest struct {
	AssistantId uuid.UUID           `param:"assistant-id" validate:"required,uuid4"`
	Type        string              `json:"type" validate:"required"`
	Name        string              `json:"name" validate:"required"`
	URL         string              `json:"url" validate:"required,url"`
	Docs        []document.Document `json:"docs,omitempty"`
}

// uploadAssistantAttachment is a handler function that uploads an attachment to an assistant.
//
//	@Summary		Upload Assistant Attachment
//	@Description	Uploads an attachment to an assistant
//	@Tags			Assistants
//	@Accept			json
//	@Produce		json
//	@Param			assistant-id	path		string										true	"Assistant ID"
//	@Success		201				{object}	JSONResult{data=assistants.AttachmentDTO}	"Created"
//	@Failure		400				{object}	JSONResult{data=nil}						"Bad Request"
//	@Failure		401				{object}	JSONResult{data=nil}						"Unauthorized"
//	@Failure		500				{object}	JSONResult{data=nil}						"Internal Server Error"
//	@Router			/assistants/{assistant-id}/attachments [post]
func (h *assistantHandler) uploadAssistantAttachment(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(uploadAssistantAttachmentRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	attachment := &assistants.AttachmentDTO{
		UserId:      user.ID,
		AssistantId: req.AssistantId,
		Type:        req.Type,
		Name:        req.Name,
		URL:         req.URL,
	}

	attachment, err = h.service.CreateAttachment(ctx, tx, attachment, req.Docs)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusCreated, attachment)
}

// @Router			/assistants/{assistant-id}/threads/{thread-id}/attachments [get].
func (h *assistantHandler) listAttachmentsByThread(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(getThreadRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, _, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	attachments, err := h.service.ListAttachmentsByThread(ctx, tx, req.ThreadId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, attachments)
}

type uploadThreadAttachmentRequest struct {
	AssistantId uuid.UUID           `param:"assistant-id" validate:"required,uuid4"`
	ThreadId    uuid.UUID           `param:"thread-id" validate:"required,uuid4"`
	Type        string              `json:"type" validate:"required"`
	Name        string              `json:"name" validate:"required"`
	URL         string              `json:"url" validate:"required,url"`
	Docs        []document.Document `json:"docs,omitempty"`
}

// uploadThreadAttachment is a handler function that uploads an attachment to a thread.
//
//	@Summary		Upload Thread Attachment
//	@Description	Uploads an attachment to a thread
//	@Tags			Assistants
//	@Accept			json
//	@Produce		json
//	@Param			assistant-id	path		string										true	"Assistant ID"
//	@Param			thread-id		path		string										true	"Thread ID"
//	@Success		201				{object}	JSONResult{data=assistants.AttachmentDTO}	"Created"
//	@Failure		400				{object}	JSONResult{data=nil}						"Bad Request"
//	@Failure		401				{object}	JSONResult{data=nil}						"Unauthorized"
//	@Failure		500				{object}	JSONResult{data=nil}						"Internal Server Error"
//	@Router			/assistants/{assistant-id}/threads/{thread-id}/attachments [post]
func (h *assistantHandler) uploadThreadAttachment(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(uploadThreadAttachmentRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	attachment := &assistants.AttachmentDTO{
		UserId:      user.ID,
		AssistantId: req.AssistantId,
		ThreadId:    req.ThreadId,
		Type:        req.Type,
		Name:        req.Name,
		URL:         req.URL,
	}

	attachment, err = h.service.CreateAttachment(ctx, tx, attachment, req.Docs)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusCreated, attachment)
}

// @Router			/assistants/models [get].
func (h *assistantHandler) listModels(c echo.Context) error {
	models, err := h.service.ListModels(c.Request().Context())
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, models)
}

// add an api to list all available tools.
func (h *assistantHandler) listTools(c echo.Context) error {
	tools, err := h.service.ListTools(c.Request().Context())
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, tools)
}
