package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/db"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type authService interface {
	GetOAuth2RedirectURL(ctx context.Context, provider string) (string, error)
	HandleOAuth2Callback(ctx context.Context, tx db.DBTX, provider, code string) (*auth.UserDTO, error)
	CreateUser(ctx context.Context, tx db.DBTX, user *auth.UserDTO) (*auth.UserDTO, error)
	AuthByPassword(ctx context.Context, tx db.DBTX, email string, password string) (*auth.UserDTO, error)
	GenerateJWT(user uuid.UUID) (string, error)
	ValidateJWT(ctx context.Context, tx db.DBTX, tokenString string) (*auth.UserDTO, int64, error)

	CreateApiKey(ctx context.Context, tx db.DBTX, key *auth.ApiKeyDTO) (*auth.ApiKeyDTO, error)
	DeleteApiKey(ctx context.Context, tx db.DBTX, id uuid.UUID) error
	ListApiKeys(ctx context.Context, tx db.DBTX, prefix string, isActive bool) ([]*auth.ApiKeyDTO, error)
}

type authHandler struct {
	service authService
}

func registerAuthHandlers(e *echo.Group) {
	h := &authHandler{
		service: auth.New(),
	}
	oauth := e.Group("/oauth")
	oauth.GET("/:provider/login", h.oAuthLogin)
	oauth.GET("/:provider/callback", h.oAuthCallback)

	auth := e.Group("/auth")
	auth.POST("/login", h.login)
	auth.POST("/logout", h.logout)
	auth.POST("/register", h.register)
	auth.GET("/validate-jwt", h.validateJwtToken)

	// Register API key handlers
	apiKeys := auth.Group("/keys", authUserMiddleware())
	apiKeys.POST("", h.createApiKey)
	apiKeys.GET("", h.listApiKeys)
	apiKeys.DELETE("/:id", h.deleteApiKey)
}

func (h *authHandler) oAuthLogin(c echo.Context) error {
	provider := c.Param("provider")
	ctx := c.Request().Context()
	redirectUrl, err := h.service.GetOAuth2RedirectURL(ctx, provider)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to get oauth redirect url: %w", err))
	}

	return JsonResponse(c, http.StatusOK, map[string]string{
		"url": redirectUrl,
	})
}

func (h *authHandler) oAuthCallback(c echo.Context) error {
	provider := c.Param("provider")
	code := c.QueryParam("code")
	ctx := c.Request().Context()
	tx, err := loadTx(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	user, err := h.service.HandleOAuth2Callback(ctx, tx, provider, code)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to get oauth token: %w", err))
	}

	jwtToken, err := h.service.GenerateJWT(user.ID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to generate jwt token: %w", err))
	}

	h.setCookieJwtToken(c, jwtToken)
	return JsonResponse(c, http.StatusOK, toUserResponse(user))
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary		Login
// @Description	Authenticate a user with email and password
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			request	body		loginRequest					true	"User login details"
// @Success		200		{object}	JSONResult{data=userResponse}	"User logged in successfully"
// @Failure		400		{object}	JSONResult{data=nil}			"Bad Request"
// @Failure		500		{object}	JSONResult{data=nil}			"Internal server error"
// @Router			/auth/login [post]
func (h *authHandler) login(c echo.Context) error {
	req := new(loginRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err))
	}

	ctx := c.Request().Context()
	tx, err := loadTx(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	user, err := h.service.AuthByPassword(ctx, tx, req.Email, req.Password)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to login: %w", err))
	}

	token, err := h.service.GenerateJWT(user.ID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to generate jwt token: %w", err))
	}

	h.setCookieJwtToken(c, token)
	return JsonResponse(c, http.StatusOK, toUserResponse(user))
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID       uuid.UUID         `json:"id"`
	Avatar   string            `json:"avatar"`
	Username string            `json:"username"`
	Email    string            `json:"email"`
	Phone    string            `json:"phone"`
	Status   string            `json:"status"`
	Settings auth.UserSettings `json:"settings"`
}

// @Summary		Register a new user
// @Description	Register a new user with the provided username, email, and password
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			request	body		registerRequest					true	"User registration details"
// @Success		200		{object}	JSONResult{data=userResponse}	"User registered successfully"
// @Failure		400		{object}	JSONResult{data=nil}			"Bad Request"
// @Failure		500		{object}	JSONResult{data=nil}			"Internal server error"
// @Router			/auth/register [post]
func (h *authHandler) register(c echo.Context) error {
	req := new(registerRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err))
	}

	ctx := c.Request().Context()
	tx, err := loadTx(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	user, err := h.service.CreateUser(ctx, tx, &auth.UserDTO{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to create user: %w", err))
	}

	token, err := h.service.GenerateJWT(user.ID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to generate jwt token: %w", err))
	}

	h.setCookieJwtToken(c, token)
	return JsonResponse(c, http.StatusOK, userResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}

func (h *authHandler) setCookieJwtToken(c echo.Context, token string) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
}

// @Summary		Validate JWT token
// @Description	Validate the JWT token and return the user ID
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Success		200	{object}	JSONResult{data=userResponse}	"
// @Failure		401	{object}	JSONResult{data=nil}			"Unauthorized"
// @Failure		500	{object}	JSONResult{data=nil}			"Internal
// @Router			/auth/validate-jwt [get]
func (h *authHandler) validateJwtToken(c echo.Context) error {
	ctx := c.Request().Context()
	tx, err := loadTx(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	token, err := c.Cookie("token")
	if err != nil {
		return ErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("failed to get jwt token: %w", err))
	}

	user, exp, err := h.service.ValidateJWT(ctx, tx, token.Value)
	if err != nil {
		return ErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("failed to validate jwt token: %w", err))
	}

	if time.Now().Add(-time.Hour*4).Unix() < exp {
		jwt, err := h.service.GenerateJWT(user.ID)
		if err != nil {
			return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to generate jwt token: %w", err))
		}
		h.setCookieJwtToken(c, jwt)
	}
	return JsonResponse(c, http.StatusOK, toUserResponse(user))
}

// @Summary		User logout
// @Description	Clear user session by removing JWT token
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Success		200	{object}	JSONResult{data=nil}	"Successfully logged out"
// @Failure		401	{object}	JSONResult{data=nil}	"Unauthorized"
// @Failure		500	{object}	JSONResult{data=nil}	"Internal server error"
// @Router			/auth/logout [post]
func (h *authHandler) logout(c echo.Context) error {
	// Remove the token cookie by setting its expiry to a past time
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour), // Set expiry to the past
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	return JsonResponse(c, http.StatusOK, nil)
}

func toUserResponse(user *auth.UserDTO) userResponse {
	return userResponse{
		ID:       user.ID,
		Avatar:   strings.ToUpper(user.Username[:1]),
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Status:   user.Status,
		Settings: user.Settings,
	}
}
