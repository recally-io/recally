package httpserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()
	registerMiddlewares(e)
	registerRouters(e)

	// Health check
	e.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	return e
}
