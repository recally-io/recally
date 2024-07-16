package handlers

import (
	"fmt"
	"net/http"
	"vibrain/internal/pkg/auth"

	"github.com/labstack/echo/v4"
)

func LoginHandler(c echo.Context) error {
	// return web/login.html page
	return c.File("web/login.html")
}

func OAuthLoginHandler(c echo.Context) error {
	provider := c.Param("provider")
	ctx := c.Request().Context()
	redirectUrl, err := auth.GetOAuth2RedirectURL(ctx, provider)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to get oauth redirect url: %w", err))
	}

	return c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func OAuthCallbackHandler(c echo.Context) error {
	provider := c.Param("provider")
	code := c.QueryParam("code")
	ctx := c.Request().Context()

	token, err := auth.GetOAuth2Token(ctx, provider, code)
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
