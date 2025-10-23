package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"recally/internal/pkg/auth"
	"recally/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type authService interface {
	GetOAuth2RedirectURL(ctx context.Context, provider string) (string, error)
	HandleOAuth2Callback(ctx context.Context, tx db.DBTX, provider, code string) (*auth.UserDTO, error)
	HandleOAuth2UserLogin(ctx context.Context, tx db.DBTX, provider, code, state string) (*auth.UserDTO, error)
	CreateUser(ctx context.Context, tx db.DBTX, user *auth.UserDTO) (*auth.UserDTO, error)
	AuthByPassword(ctx context.Context, tx db.DBTX, email, password string) (*auth.UserDTO, error)
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
	// Initialize auth service with Goth adapter for OAuth
	dao := db.New()
	oauthAdapter := auth.InitGothAdapter(dao)
	authService := auth.NewWithAdapter(oauthAdapter)

	h := &authHandler{
		service: authService,
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

	// For backward compatibility: use old flow if GetOAuth2RedirectURL exists
	// Otherwise, fall back to the legacy OAuth flow
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
	state := c.QueryParam("state")
	ctx := c.Request().Context()

	tx, err := loadTx(ctx)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	// Use new Goth-based OAuth flow with state validation
	var user *auth.UserDTO
	if state != "" {
		// New flow with CSRF protection via state validation
		user, err = h.service.HandleOAuth2UserLogin(ctx, tx, provider, code, state)
		if err != nil {
			return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to complete oauth login: %w", err))
		}
	} else {
		// Legacy flow for backward compatibility (no state validation)
		user, err = h.service.HandleOAuth2Callback(ctx, tx, provider, code)
		if err != nil {
			return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to get oauth token: %w", err))
		}
	}

	jwtToken, err := h.service.GenerateJWT(user.ID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to generate jwt token: %w", err))
	}

	h.setSecureCookieJwtToken(c, jwtToken)

	return JsonResponse(c, http.StatusOK, toUserResponse(user))
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Router	/auth/login [post].
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

// @Router	/auth/register [post].
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

// setSecureCookieJwtToken sets a secure JWT cookie with enhanced security flags
// This method should be used for OAuth flows to ensure cookies are properly protected
func (h *authHandler) setSecureCookieJwtToken(c echo.Context, token string) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		Path:     "/",
		HttpOnly: true,                    // Prevent JavaScript access (XSS protection)
		Secure:   c.IsTLS(),               // Only send over HTTPS in production
		SameSite: http.SameSiteStrictMode, // CSRF protection (Strict for OAuth callbacks)
	}
	c.SetCookie(cookie)
}

// @Router	/auth/validate-jwt [get].
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

// @Router	/auth/logout [post].
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
