package contexts

import (
	"context"
	"vibrain/internal/pkg/constant"

	"github.com/labstack/echo/v4"
)

func Set(c echo.Context, key string, value interface{}) {
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, constant.ContextKey(key), value)

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
	return c.Request().Context().Value(constant.ContextKey(key))
}
