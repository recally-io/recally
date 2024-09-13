package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
	"vibrain/internal/core/assistants"
	"vibrain/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

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
	req := new(getAssistantRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	tx, _, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	threads, err := h.service.ListThreads(ctx, tx, req.AssistantId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, threads)
}

type createThreadRequest struct {
	AssistantId  uuid.UUID `param:"assistant-id" validate:"required,uuid4"`
	Id           uuid.UUID `json:"id,omitempty" validate:"omitempty,uuid"`
	Name         string    `json:"name" validate:"required"`
	Description  string    `json:"description,omitempty"`
	Model        string    `json:"model,omitempty"`
	SystemPrompt string    `json:"system_prompt,omitempty"`
	Metadata     struct {
		Tools []string `json:"tools,omitempty"`
	} `json:"metadata,omitempty"`
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
	req := new(createThreadRequest)
	if err := bindAndValidate(c, req); err != nil {
		logger.FromContext(ctx).Error("bind request error", "err", err)
		return err
	}
	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	// Assign the AssistantID for the new thread
	threadDTO := assistants.ThreadDTO{
		Id:           req.Id,
		Name:         req.Name,
		Description:  req.Description,
		Model:        req.Model,
		SystemPrompt: req.SystemPrompt,
		AssistantId:  req.AssistantId,
		UserId:       user.ID,
		Metadata: assistants.ThreadMetadata{
			Tools: req.Metadata.Tools,
		},
	}

	thread, err := h.service.CreateThread(ctx, tx, &threadDTO)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to create thread: %w", err))
	}

	return JsonResponse(c, http.StatusCreated, thread)
}

type updateThreadRequest struct {
	AssistantId  uuid.UUID `param:"assistant-id" validate:"required,uuid4"`
	ThreadId     uuid.UUID `param:"thread-id" validate:"required,uuid4"`
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	Model        string    `json:"model,omitempty"`
	SystemPrompt string    `json:"system_prompt,omitempty"`
	Metadata     struct {
		Tools []string `json:"tools,omitempty"`
	} `json:"metadata,omitempty"`
}

// updateThread is a handler function that updates an existing thread for an assistant.
// It retrieves the assistant ID and thread ID from the request parameters and uses them to update the thread.
// If the assistant ID or thread ID is not found in the parameters, it returns an error with status code 400 (Bad Request).
// If there is an error while updating the thread, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it returns a JSON response with status code 200 (OK) and the updated thread.
// @Summary Update Thread
// @Description Updates an existing thread under an assistant
// @Tags Assistants
// @Accept json
// @Produce json
// @Param assistant-id path string true "Assistant ID"
// @Param thread-id path string true "Thread ID"
// @Param thread body assistants.ThreadDTO true "Thread"
// @Success 200 {object} JSONResult{data=assistants.ThreadDTO} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/{assistant-id}/threads/{thread-id} [put]
func (h *assistantHandler) updateThread(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(updateThreadRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}
	tx, _, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	thread, err := h.service.GetThread(ctx, tx, req.ThreadId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if req.Name != "" {
		thread.Name = req.Name
	}
	if req.Description != "" {
		thread.Description = req.Description
	}
	if req.Model != "" {
		thread.Model = req.Model
	}
	if req.SystemPrompt != "" {
		thread.SystemPrompt = req.SystemPrompt
	}
	if req.Metadata.Tools != nil {
		thread.Metadata.Tools = req.Metadata.Tools
	}

	thread, err = h.service.UpdateThread(ctx, tx, thread)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, thread)
}

type getThreadRequest struct {
	AssistantId uuid.UUID `param:"assistant-id" validate:"required,uuid4"`
	ThreadId    uuid.UUID `param:"thread-id" validate:"required,uuid4"`
}

// getThread is a handler function that retrieves a thread by ID.
// It retrieves the thread ID from the request parameters and uses it to fetch the thread.
// If the thread ID is not found in the parameters, it returns an error with status code 400 (Bad Request).
// If there is an error while fetching the thread, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it returns a JSON response with status code 200 (OK) and the thread.
//
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
	req := new(getThreadRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}
	tx, _, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	thread, err := h.service.GetThread(ctx, tx, req.ThreadId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	messages, err := h.service.ListThreadMessages(ctx, tx, req.ThreadId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	thread.Messages = messages
	return JsonResponse(c, http.StatusOK, thread)
}

// deleteThread is a handler function that deletes a thread by ID.
// It retrieves the thread ID from the request parameters and uses it to delete the thread.
// If the thread ID is not found in the parameters, it returns an error with status code 400 (Bad Request).
// If there is an error while deleting the thread, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it returns a JSON response with status code 204 (No Content).

// @Summary Delete Thread
// @Description Deletes a thread by ID
// @Tags Assistants
// @Accept json
// @Produce json
// @Param assistant-id path string true "Assistant ID"
// @Param thread-id path string true "Thread ID"
// @success 204 {object} JSONResult{data=nil} "No Content"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/{assistant-id}/threads/{thread-id} [delete]
func (h *assistantHandler) deleteThread(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(getThreadRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}
	tx, _, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	err = h.service.DeleteThread(ctx, tx, req.ThreadId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusNoContent, nil)
}

// generateThreadTitle is a handler function that generates a title for a thread.
// It retrieves the thread ID from the request parameters and uses it to generate the title.
// If the thread ID is not found in the parameters, it returns an error with status code 400 (Bad Request).
// If there is an error while generating the title, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it returns a JSON response with status code 200 (OK) and the generated title.

// @Summary Generate Thread Title
// @Description Generates a title for a thread based on the conversation
// @Tags Assistants
// @Accept json
// @Produce json
// @Param assistant-id path string true "Assistant ID"
// @Param thread-id path string true "Thread ID"
// @success 200 {object} JSONResult{data=string} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/{assistant-id}/threads/{thread-id}/generate-title [post]
func (h *assistantHandler) generateThreadTitle(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(getThreadRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}
	tx, _, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	title, err := h.service.GenerateThreadTitle(ctx, tx, req.ThreadId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return JsonResponse(c, http.StatusOK, title)
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
// @success 200 {object} JSONResult{data=[]assistants.MessageDTO} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/{assistant-id}/threads/{thread-id}/messages [get]
func (h *assistantHandler) listThreadMessages(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(getThreadRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}
	tx, _, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	messages, err := h.service.ListThreadMessages(ctx, tx, req.ThreadId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusOK, messages)
}

type createThreadMessageRequest struct {
	AssistantId uuid.UUID `param:"assistant-id" validate:"required,uuid4"`
	ThreadId    uuid.UUID `param:"thread-id" validate:"required,uuid4"`
	Role        string    `json:"role" validate:"required"`
	Text        string    `json:"text" validate:"required"`
	Model       string    `json:"model,omitempty"`
	Metadata    struct {
		Tools  []string `json:"tools,omitempty"`
		Images []string `json:"images,omitempty"`
	} `json:"metadata,omitempty"`
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
// @Success 201 {object} JSONResult{data=assistants.MessageDTO} "Created"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/{assistant-id}/threads/{thread-id}/messages [post]
func (h *assistantHandler) createThreadMessage(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(createThreadMessageRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}
	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	thread, err := h.service.GetThread(ctx, tx, req.ThreadId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	messageDTO := assistants.MessageDTO{
		UserID:   user.ID,
		ThreadID: thread.Id,
		Model:    thread.Model,
		Role:     req.Role,
		Text:     req.Text,
		Metadata: assistants.MessageMetadata{
			Tools:  req.Metadata.Tools,
			Images: req.Metadata.Images,
			Stream: true,
		},
	}

	if req.Model != "" {
		messageDTO.Model = req.Model
	}

	// Create Thread Message
	if _, err := h.service.CreateThreadMessage(ctx, tx, thread.Id, &messageDTO); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	msgChan := make(chan *assistants.MessageDTO)
	errChan := make(chan error)

	streamingFunc := func(msg *assistants.MessageDTO, err error) {
		if err != nil {
			errChan <- err
			return
		}
		msgChan <- msg
	}
	go h.service.RunThread(ctx, tx, thread.Id, streamingFunc)

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.Request().Context().Done():
			logger.FromContext(ctx).Info("SSE client disconnected", "ip", c.RealIP())
			return nil
		case err := <-errChan:
			if errors.Is(err, io.EOF) {
				w.Flush()
				return nil
			}
			return ErrorResponse(c, http.StatusInternalServerError, err)
		case msg := <-msgChan:
			data, err := json.Marshal(msg)
			if err != nil {
				return ErrorResponse(c, http.StatusInternalServerError, err)
			}
			event := Event{
				Data: data,
			}
			if err := event.MarshalTo(w); err != nil {
				return ErrorResponse(c, http.StatusInternalServerError, err)
			}
			w.Flush()
		}
	}
}

type getThreadMessageRequest struct {
	AssistantId uuid.UUID `param:"assistant-id" validate:"required,uuid4"`
	ThreadId    uuid.UUID `param:"thread-id" validate:"required,uuid4"`
	MessageId   uuid.UUID `param:"message-id" validate:"required,uuid4"`
}

// deleteThreadMessage is a handler function that deletes a thread message by ID.
// It retrieves the thread message ID from the request parameters and uses it to delete the thread message.
// If the thread message ID is not found in the parameters, it returns an error with status code 400 (Bad Request).
// If there is an error while deleting the thread message, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it returns a JSON response with status code 204 (No Content).

// @Summary Delete Thread Message
// @Description Deletes a thread message by ID
// @Tags Assistants
// @Accept json
// @Produce json
// @Param assistant-id path string true "Assistant ID"
// @Param thread-id path string true "Thread ID"
// @Param message-id path string true "Message ID"
// @success 204 {object} JSONResult{data=nil} "No Content"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/{assistant-id}/threads/{thread-id}/messages/{message-id} [delete]
func (h *assistantHandler) deleteThreadMessage(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(getThreadMessageRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}
	tx, _, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	err = h.service.DeleteThreadMessage(ctx, tx, req.MessageId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return JsonResponse(c, http.StatusNoContent, nil)
}

type updateThreadMessageRequest struct {
	AssistantId uuid.UUID `param:"assistant-id" validate:"required,uuid4"`
	ThreadId    uuid.UUID `param:"thread-id" validate:"required,uuid4"`
	MessageId   uuid.UUID `param:"message-id" validate:"required,uuid4"`
	Model       string    `json:"model,omitempty"`
	Text        string    `json:"text" validate:"required"`
	Metadata    struct {
		Tools []string `json:"tools,omitempty"`
	} `json:"metadata,omitempty"`
}

// update Thread Message
// updateThreadMessage is a handler function that updates a thread message by ID.
// It retrieves the thread message ID from the request parameters and the updated message data from the request body.
// If the thread message ID is not found in the parameters or the request body is invalid, it returns an error with status code 400 (Bad Request).
// If there is an error while updating the thread message, it returns an error with status code 500 (Internal Server Error).
// Otherwise, it returns a JSON response with status code 200 (OK) and the updated thread message.

// @Summary Update Thread Message
// @Description Updates a thread message by ID
// @Tags Assistants
// @Accept json
// @Produce json
// @Param assistant-id path string true "Assistant ID"
// @Param thread-id path string true "Thread ID"
// @Param message-id path string true "Message ID"
// @Param message body assistants.MessageDTO true "Updated Message"
// @Success 200 {object} JSONResult{data=assistants.MessageDTO} "Success"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /assistants/{assistant-id}/threads/{thread-id}/messages/{message-id} [put]
func (h *assistantHandler) updateThreadMessage(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(updateThreadMessageRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	var updatedMessage assistants.MessageDTO
	if err := c.Bind(&updatedMessage); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	message, err := h.service.GetThreadMessage(ctx, tx, req.MessageId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	err = h.service.DeleteThreadMessage(ctx, tx, req.MessageId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	thread, err := h.service.GetThread(ctx, tx, req.ThreadId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	messageDTO := assistants.MessageDTO{
		UserID:   user.ID,
		ThreadID: thread.Id,
		Model:    message.Model,
		Role:     message.Role,
		Text:     req.Text,
		Metadata: assistants.MessageMetadata{
			Tools: req.Metadata.Tools,
		},
	}

	if req.Model != "" {
		messageDTO.Model = req.Model
	}
	return JsonResponse(c, http.StatusOK, messageDTO)
	// Create Thread Message
	// if _, err := h.service.CreateThreadMessage(ctx, tx, thread.Id, &messageDTO); err != nil {
	// 	return ErrorResponse(c, http.StatusInternalServerError, err)
	// }

	// resp, err := h.service.RunThread(ctx, tx, thread.Id)
	// if err != nil {
	// 	return ErrorResponse(c, http.StatusInternalServerError, err)
	// }

	// return JsonResponse(c, http.StatusOK, resp)
}
