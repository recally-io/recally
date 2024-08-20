package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"vibrain/internal/pkg/auth"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type authService interface {
	GetOAuth2RedirectURL(ctx context.Context, provider string) (string, error)
	GetOAuth2Token(ctx context.Context, provider, code string) (*oauth2.Token, error)
	CreateUser(ctx context.Context, tx db.DBTX, user *auth.UserDTO) (*auth.UserDTO, error)
	AuthByPassword(ctx context.Context, tx db.DBTX, email string, password string) (*auth.UserDTO, error)
	GenerateJWT(user uuid.UUID) (string, error)
	ValidateJWT(tokenString string) (uuid.UUID, int64, error)
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
	auth.POST("/register", h.register)
	auth.GET("/validate-jwt", h.validateJwtToken)
}

func (h *authHandler) oAuthLogin(c echo.Context) error {
	provider := c.Param("provider")
	ctx := c.Request().Context()
	redirectUrl, err := h.service.GetOAuth2RedirectURL(ctx, provider)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to get oauth redirect url: %w", err))
	}

	return c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func (h *authHandler) oAuthCallback(c echo.Context) error {
	provider := c.Param("provider")
	code := c.QueryParam("code")
	ctx := c.Request().Context()

	_, err := h.service.GetOAuth2Token(ctx, provider, code)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to get oauth token: %w", err))
	}

	userId := uuid.New()
	jwtToken, err := h.service.GenerateJWT(userId)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to generate jwt token: %w", err))
	}

	h.setCookieJwtToken(c, jwtToken)
	return JsonResponse(c, http.StatusOK, userResponse{
		ID: userId,
	})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary Login
// @Description Authenticate a user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "User login details"
// @Success 200 {object} JSONResult{data=userResponse} "User logged in successfully"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} JSONResult{data=nil} "Internal server error"
// @Router /auth/login [post]
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
	return JsonResponse(c, http.StatusOK, userResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

// @Summary Register a new user
// @Description Register a new user with the provided username, email, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "User registration details"
// @Success 200 {object} JSONResult{data=userResponse} "User registered successfully"
// @Failure 400 {object} JSONResult{data=nil} "Bad Request"
// @Failure 500 {object} JSONResult{data=nil} "Internal server error"
// @Router /auth/register [post]
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

	return c.JSON(http.StatusOK, userResponse{
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

// @Summary Validate JWT token
// @Description Validate the JWT token and return the user ID
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} JSONResult{data=userResponse} "
// @Failure 401 {object} JSONResult{data=nil} "Unauthorized"
// @Failure 500 {object} JSONResult{data=nil} "Internal
// @Router /auth/validate-jwt [get]
func (h *authHandler) validateJwtToken(c echo.Context) error {
	token, err := c.Cookie("token")
	if err != nil {
		return ErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("failed to get jwt token: %w", err))
	}

	userId, exp, err := h.service.ValidateJWT(token.Value)
	if err != nil {
		return ErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("failed to validate jwt token: %w", err))
	}

	if time.Now().Add(-time.Hour*4).Unix() < exp {
		jwt, err := h.service.GenerateJWT(userId)
		if err != nil {
			return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to generate jwt token: %w", err))
		}
		h.setCookieJwtToken(c, jwt)
	}

	return c.JSON(http.StatusOK, userResponse{
		ID: userId,
	})
}
