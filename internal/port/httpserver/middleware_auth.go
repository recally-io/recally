package httpserver

import (
	"fmt"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/contexts"
	"recally/internal/pkg/db"

	"github.com/google/uuid"
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

	err = loadAndSetUser(c, userId)
	return true, err
}

func loadAndSetUser(c echo.Context, userId uuid.UUID) error {
	ctx := c.Request().Context()
	tx, err := loadTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to load transaction: %w", err)
	}

	dao := db.New()
	dbUser, err := dao.GetUserById(ctx, tx, userId)
	if err != nil {
		return fmt.Errorf("failed to load user: %w", err)
	}

	user := new(auth.UserDTO)
	user.Load(&dbUser)
	setContext(c, contexts.ContextKeyUser, user)
	return nil
}
