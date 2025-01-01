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

// @Summary Update user info
// @Description Update user's username, email, and phone
// @Tags Auth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body updateUserInfoRequest true "User info update details"
// @Success 200 {object} JSONResult{data=userResponse} "User info updated successfully"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} JSONResult{data=nil} "Internal server error"
// @Router /auth/user/info [put]
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

// @Summary Update user settings
// @Description Update user's settings
// @Tags Auth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body updateUserSettingsRequest true "User settings update"
// @Success 200 {object} JSONResult{data=userResponse} "User settings updated successfully"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} JSONResult{data=nil} "Internal server error"
// @Router /auth/user/settings [put]
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

// @Summary Update user password
// @Description Update user's password
// @Tags Auth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body updateUserPasswordRequest true "User password update"
// @Success 200 {object} JSONResult{data=userResponse} "User password updated successfully"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} JSONResult{data=nil} "Internal server error"
// @Router /auth/user/password [put]
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
