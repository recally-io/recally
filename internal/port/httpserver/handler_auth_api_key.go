package httpserver

import (
	"fmt"
	"net/http"
	"recally/internal/pkg/auth"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type createApiKeyRequest struct {
	Name      string    `json:"name" validate:"required,min=3,max=255"`
	Prefix    string    `json:"prefix"`
	Scopes    []string  `json:"scopes"`
	ExpiresAt time.Time `json:"expires_at" validate:"required,gt=now"`
}

// @Summary Create API Key
// @Description Create a new API key for the authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body createApiKeyRequest true "API key creation details"
// @Success 200 {object} JSONResult{data=auth.ApiKeyDTO} "API key created successfully"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal server error"
// @Router /keys [post]
func (h *authHandler) createApiKey(c echo.Context) error {
	ctx := c.Request().Context()
	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	req := new(createApiKeyRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err))
	}

	apiKey, err := h.service.CreateApiKey(ctx, tx, &auth.ApiKeyDTO{
		UserID:    user.ID,
		Name:      req.Name,
		Prefix:    req.Prefix,
		Scopes:    req.Scopes,
		ExpiresAt: req.ExpiresAt,
	})
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to create API key: %w", err))
	}

	return JsonResponse(c, http.StatusOK, apiKey)
}

type listApiKeysRequest struct {
	Prefix   string `query:"prefix"`
	IsActive bool   `query:"is_active"`
}

// @Summary List API Keys
// @Description List all API keys for the authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param prefix query string false "Filter by key prefix"
// @Param is_active query bool false "Filter by active status"
// @Success 200 {object} JSONResult{data=[]auth.ApiKeyDTO} "API keys retrieved successfully"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal server error"
// @Router /keys [get]
func (h *authHandler) listApiKeys(c echo.Context) error {
	ctx := c.Request().Context()
	tx, _, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	req := new(listApiKeysRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err))
	}

	keys, err := h.service.ListApiKeys(ctx, tx, req.Prefix, req.IsActive)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to list API keys: %w", err))
	}
	return JsonResponse(c, http.StatusOK, keys)
}

type deleteApiKeyRequest struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

// @Summary Delete API Key
// @Description Delete an API key by ID
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "API key ID"
// @Success 200 {object} JSONResult{data=nil} "API key deleted successfully"
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal server error"
// @Router /keys/{id} [delete]
func (h *authHandler) deleteApiKey(c echo.Context) error {
	ctx := c.Request().Context()
	tx, _, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	req := new(deleteApiKeyRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err))
	}

	if err := h.service.DeleteApiKey(ctx, tx, req.ID); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to delete API key: %w", err))
	}

	return JsonResponse(c, http.StatusOK, nil)
}
