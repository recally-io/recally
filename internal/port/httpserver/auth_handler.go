package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"vibrain/internal/pkg/auth"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type authService interface {
	GetOAuth2RedirectURL(ctx context.Context, provider string) (string, error)
	GetOAuth2Token(ctx context.Context, provider, code string) (*oauth2.Token, error)
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

	token, err := h.service.GetOAuth2Token(ctx, provider, code)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to get oauth token: %w", err))
	}

	// TODO: get user info from database
	userId := "userId"

	jwtUser := auth.JwtUser{
		UserID:        userId,
		OAuthProvider: provider,
		AccessToken:   token.AccessToken,
		TokenType:     token.TokenType,
		Expiry:        token.Expiry,
	}

	jwtToken, err := auth.GenerateJWT(jwtUser)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to generate jwt token: %w", err))
	}

	// write jwt token to cookie
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = jwtToken
	cookie.Expires = token.Expiry
	c.SetCookie(cookie)
	return JsonResponse(c, http.StatusOK, map[string]interface{}{
		"jwt_token": jwtToken,
		"jwt_user":  jwtUser,
	})
}
