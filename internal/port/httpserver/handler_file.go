package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"recally/internal/core/files"
	"recally/internal/pkg/db"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type fileService interface {
	DeleteFile(ctx context.Context, tx db.DBTX, id uuid.UUID) error
	GetPublicURL(ctx context.Context, objectKey string) (string, error)
	GetShareURL(ctx context.Context, objectKey string) string
	GetPresignedPutObjectURL(ctx context.Context, userID uuid.UUID, fileName string, expires time.Duration) (string, string, error)
}

type fileHandler struct {
	service fileService
}

func registerFileHandlers(e *echo.Group, s *Service) {
	h := &fileHandler{
		service: files.NewService(s.s3),
	}
	files := e.Group("/files", authUserMiddleware())
	files.GET("/presigned-urls", h.getPresignedURLs)
	files.DELETE("/:id", h.deleteFile)
	files.GET("/file", h.getFile)
}

type getPresignedURLsRequest struct {
	FileName string `query:"file_name" validate:"required"`
	FileType string `query:"file_type" validate:"required"`
	Action   string `query:"action" validate:"required,oneof=PUT GET"`
	// Expiration in seconds
	Expiration int `query:"expiration" validate:"required,min=1,max=604800"`
}

type getPresignedURLsResponse struct {
	PresignedURL string `json:"presigned_url"`
	ObjectKey    string `json:"object_key"`
	PublicURL    string `json:"public_url"`
}

type getFileRequest struct {
	ObjectKey string `query:"object_key" validate:"required"`
}

// @Summary		Get presigned URLs for file operations
// @Description	Get presigned URLs for uploading or downloading files from S3
// @Tags			files
// @Accept			json
// @Produce		json
// @Param			fileName	query		string										true	"Name of the file"
// @Param			fileType	query		string										true	"MIME type of the file"
// @Param			action		query		string										false	"Action to perform (put or get)"	Enums(put, get)
// @Param			expiration	query		int											false	"Expiration time in seconds (max 604800)"
// @Success		200			{object}	JSONResult{data=getPresignedURLsResponse}	"Created"
// @Failure		400			{object}	JSONResult{data=nil}						"Bad Request"
// @Failure		401			{object}	JSONResult{data=nil}						"Unauthorized"
// @Failure		500			{object}	JSONResult{data=nil}						"Internal Server Error"
// @Router			/files/presigned-urls [get]
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

	// Default expiration to 1 hour if not provided
	if req.Expiration == 0 {
		req.Expiration = 3600
	}
	expirationDuration := time.Duration(req.Expiration) * time.Second
	presignedURL, objectKey, err := h.service.GetPresignedPutObjectURL(ctx, user.ID, req.FileName, expirationDuration)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to generate presigned URL: %w", err))
	}

	// Return the presigned URL
	return JsonResponse(c, http.StatusOK, getPresignedURLsResponse{
		PresignedURL: presignedURL,
		ObjectKey:    objectKey,
		PublicURL:    h.service.GetShareURL(ctx, objectKey),
	})
}

type getPublicURLResponse struct {
	URL string `json:"url"`
}

func (h *fileHandler) getFile(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(getFileRequest)
	if err := bindAndValidate(c, req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	publicURL, err := h.service.GetPublicURL(ctx, req.ObjectKey)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to get public URL: %w", err))
	}

	return JsonResponse(c, http.StatusOK, getPublicURLResponse{
		URL: publicURL,
	})
}

type deleteFileRequest struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

// @Summary		Delete a file
// @Description	Delete a file by its ID
// @Tags			files
// @Produce		json
// @Param			id	path		string					true	"File ID"
// @Success		200	{object}	JSONResult{data=nil}	"Created"
// @Failure		400	{object}	JSONResult{data=nil}	"Bad Request"
// @Failure		401	{object}	JSONResult{data=nil}	"Unauthorized"
// @Failure		500	{object}	JSONResult{data=nil}	"Internal Server Error"
// @Router			/files/{id} [delete]
func (h *fileHandler) deleteFile(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(deleteFileRequest)
	if err := bindAndValidate(c, req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}
	tx, err := loadTx(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	err = h.service.DeleteFile(ctx, tx, req.ID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to delete file: %w", err))
	}
	return JsonResponse(c, http.StatusOK, nil)
}
