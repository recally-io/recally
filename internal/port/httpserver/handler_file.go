package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"recally/internal/core/files"
	"recally/internal/pkg/config"
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
	files.DELETE("/:id", h.deleteFile)
	files.GET("/file/presigned", h.getPresignedURLs)
	files.GET("/file/url", h.getFilePublicURLByObjectKey)
	files.GET("/file/content", h.getFileContentByObjectKey)
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

// @Router	/files/presigned-urls [get].
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

	publicURL := ""
	if config.Settings.S3.PublicURL != "" {
		publicURL, _ = h.service.GetPublicURL(ctx, objectKey)
	}

	// Return the presigned URL
	return JsonResponse(c, http.StatusOK, getPresignedURLsResponse{
		PresignedURL: presignedURL,
		ObjectKey:    objectKey,
		PublicURL:    publicURL,
	})
}

type getPublicURLResponse struct {
	URL string `json:"url"`
}

func (h *fileHandler) getFileUrlByObjectKey(c echo.Context) (string, error) {
	ctx := c.Request().Context()

	req := new(getFileRequest)
	if err := bindAndValidate(c, req); err != nil {
		return "", ErrorResponse(c, http.StatusBadRequest, err)
	}

	publicURL, err := h.service.GetPublicURL(ctx, req.ObjectKey)
	if err != nil {
		return "", ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to get public URL: %w", err))
	}

	return publicURL, nil
}

func (h *fileHandler) getFilePublicURLByObjectKey(c echo.Context) error {
	publicURL, err := h.getFileUrlByObjectKey(c)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to get public URL: %w", err))
	}

	return JsonResponse(c, http.StatusOK, getPublicURLResponse{
		URL: publicURL,
	})
}

func (h *fileHandler) getFileContentByObjectKey(c echo.Context) error {
	publicURL, err := h.getFileUrlByObjectKey(c)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to get public URL: %w", err))
	}

	return c.Redirect(http.StatusTemporaryRedirect, publicURL)
}

type deleteFileRequest struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

// @Router	/files/{id} [delete].
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
