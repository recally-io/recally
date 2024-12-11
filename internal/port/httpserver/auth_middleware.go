package httpserver

import (
	"fmt"
	"vibrain/internal/pkg/auth"
	"vibrain/internal/pkg/contexts"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func authUserMiddleware() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "cookie:token,header:Authorization",
		Validator: authValidation,
	})
}

func authAdminMiddleware() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "cookie:token,header:Authorization",
		Validator: authValidation,
	})
}

func authValidation(key string, c echo.Context) (bool, error) {
	// validate key
	userId, _, err := auth.ValidateJWT(key)
	if err != nil {
		return false, fmt.Errorf("invalid token: %w", err)
	}
	setContext(c, contexts.ContextKeyUserID, userId)
	return true, nil
}
