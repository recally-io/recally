package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"recally/internal/pkg/auth"
	"recally/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type usersService interface {
	UpdateUserInfo(ctx context.Context, tx db.DBTX, username, email, phone *string) (*auth.UserDTO, error)
	UpdateUserSettings(ctx context.Context, tx db.DBTX, settings auth.UserSettings) (*auth.UserDTO, error)
	UpdateUserPassword(ctx context.Context, tx db.DBTX, currentPassword, password string) (*auth.UserDTO, error)
}

type usersHandler struct {
	service usersService
}

func registerUsersHandlers(e *echo.Group) {
	h := &usersHandler{
		service: auth.New(),
	}

	// Add new routes
	users := e.Group("/users", authUserMiddleware())
	users.PUT("/:id/settings", h.updateUserSettings)
	users.PUT("/:id/info", h.updateUserInfo)
	users.PUT("/:id/password", h.updateUserPassword)
}

type updateUserInfoRequest struct {
	Id       uuid.UUID `param:"id" validate:"required,uuid4"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
}

// @Router	/auth/user/info [put].
func (h *usersHandler) updateUserInfo(c echo.Context) error {
	req := new(updateUserInfoRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err))
	}

	ctx := c.Request().Context()

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if user.ID != req.Id {
		return ErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("unauthorized to update user info"))
	}

	user, err = h.service.UpdateUserInfo(ctx, tx, &req.Username, &req.Email, &req.Phone)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to update user info: %w", err))
	}

	return JsonResponse(c, http.StatusOK, toUserResponse(user))
}

type updateUserSettingsRequest struct {
	Id       uuid.UUID         `param:"id" validate:"required,uuid4"`
	Settings auth.UserSettings `json:"settings"`
}

// @Router	/auth/user/settings [put].
func (h *usersHandler) updateUserSettings(c echo.Context) error {
	req := new(updateUserSettingsRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err))
	}

	ctx := c.Request().Context()

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if user.ID != req.Id {
		return ErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("unauthorized to update user info"))
	}

	user, err = h.service.UpdateUserSettings(ctx, tx, req.Settings)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to update user settings: %w", err))
	}

	return JsonResponse(c, http.StatusOK, toUserResponse(user))
}

type updateUserPasswordRequest struct {
	Id              uuid.UUID `param:"id" validate:"required,uuid4"`
	CurrentPassword string    `json:"current_password"`
	Password        string    `json:"password"`
}

// @Router	/auth/user/password [put].
func (h *usersHandler) updateUserPassword(c echo.Context) error {
	req := new(updateUserPasswordRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err))
	}

	ctx := c.Request().Context()

	tx, user, err := initContext(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if user.ID != req.Id {
		return ErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("unauthorized to update user info"))
	}

	user, err = h.service.UpdateUserPassword(ctx, tx, req.CurrentPassword, req.Password)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to update user password: %w", err))
	}

	return JsonResponse(c, http.StatusOK, toUserResponse(user))
}
