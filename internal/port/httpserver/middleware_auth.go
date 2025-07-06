package httpserver

import (
	"fmt"
	"net/http"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/contexts"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const authLookupKey = "cookie:token,header:Authorization,header:X-Api-Key"

func authUserMiddleware() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:    authLookupKey,
		Validator:    authValidation,
		ErrorHandler: authErrorHandler,
	})
}

func authAdminMiddleware() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:    authLookupKey,
		Validator:    authValidation,
		ErrorHandler: authErrorHandler,
	})
}

func authValidation(key string, c echo.Context) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("missing key")
	}

	ctx := c.Request().Context()

	tx, err := loadTx(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to load transaction: %w", err)
	}

	// validate key
	authService := auth.New()

	// check if it's a JWT token
	user, _, err := authService.ValidateJWT(ctx, tx, key)
	if err == nil {
		setContext(c, contexts.ContextKeyUser, user)

		return true, nil
	}

	// check if it's an API key
	user, err = authService.ValidateApiKey(ctx, tx, key)
	if err == nil {
		setContext(c, contexts.ContextKeyUser, user)

		return true, nil
	}

	return false, fmt.Errorf("invalid key: %w", err)
}

func authErrorHandler(err error, c echo.Context) error {
	return ErrorResponse(c, http.StatusUnauthorized, err)
}
