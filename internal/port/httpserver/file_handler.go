package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type fileService interface {
	GetPresignedURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error)
	Delete(ctx context.Context, id string) error
	GetPublicURL(objectKey string) string
}

type fileHandler struct {
	service fileService
}

func registerFileHandlers(e *echo.Group, s fileService) {
	h := &fileHandler{
		service: s,
	}
	files := e.Group("/files")
	files.GET("/presigned-urls", h.getPresignedURLs)
	files.DELETE("/:id", h.deleteFile)
}

type getPresignedURLsRequest struct {
	AssistantId uuid.UUID `query:"assistant_id" validate:"required,uuid4"`
	ThreadId    uuid.UUID `query:"thread_id,omitempty" validate:"omitempty,uuid4"`
	FileName    string    `query:"file_name" validate:"required"`
	FileType    string    `query:"file_type" validate:"required"`
	Action      string    `query:"action" validate:"required,oneof=put get"`
	// Expiration in seconds
	Expiration int `query:"expiration" validate:"required,min=1,max=604800"`
}

type getPresignedURLsResponse struct {
	PresignedURL string `json:"presigned_url"`
	PublicURL    string `json:"public_url"`
}

// @Summary Get presigned URLs for file operations
// @Description Get presigned URLs for uploading or downloading files from S3
// @Tags files
// @Accept json
// @Produce json
// @Param fileName query string true "Name of the file"
// @Param fileType query string true "MIME type of the file"
// @Param action query string false "Action to perform (put or get)" Enums(put, get)
// @Param expiration query int false "Expiration time in seconds (max 604800)"
// @Success 200 {object} JSONResult{data=getPresignedURLsResponse} "Created"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /files/presigned-urls [get]
func (h *fileHandler) getPresignedURLs(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(getPresignedURLsRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}

	_, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	objectKey := fmt.Sprintf("%s/%s/%s/%s-%s", user.ID, req.AssistantId, time.Now().Format("2006-01"), uuid.New().String(), req.FileName)
	// objectKey := url.PathEscape(fmt.Sprintf("%s/%s/%s/%s", user.ID, req.AssistantId, time.Now().Format("2006-01"), req.FileName))

	// Default expiration to 1 hour if not provided
	if req.Expiration == 0 {
		req.Expiration = 3600
	}
	expirationDuration := time.Duration(req.Expiration) * time.Second
	presignedURL, err := h.service.GetPresignedURL(c.Request().Context(), objectKey, expirationDuration)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to generate presigned URL: %w", err))
	}

	// Return the presigned URL
	return JsonResponse(c, http.StatusOK, getPresignedURLsResponse{
		PresignedURL: presignedURL,
		PublicURL:    h.service.GetPublicURL(objectKey),
	})
}

// @Summary Delete a file
// @Description Delete a file by its ID
// @Tags files
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} JSONResult{data=nil} "Created"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal Server Error"
// @Router /files/{id} [delete]
func (h *fileHandler) deleteFile(c echo.Context) error {
	id := c.Param("id")
	err := h.service.Delete(c.Request().Context(), id)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to delete file: %w", err))
	}
	return JsonResponse(c, http.StatusOK, nil)
}
