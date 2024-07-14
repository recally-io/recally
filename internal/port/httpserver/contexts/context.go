package contexts

import (
	"context"
	"vibrain/internal/pkg/config"

	"github.com/labstack/echo/v4"
)

func Set(c echo.Context, key string, value interface{}) {
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, config.ContextKey(key), value)

	// set to context.Context
	c.SetRequest(c.Request().WithContext(ctx))

	// set to echo.Context
	c.Set(key, value)
}

func Get(c echo.Context, key string) interface{} {

	// get from echo.Context
	if v := c.Get(key); v != nil {
		return v
	}
	// get from context.Context
	return c.Request().Context().Value(config.ContextKey(key))
}
