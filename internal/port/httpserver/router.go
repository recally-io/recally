package httpserver

import (
	"net/http"
	"vibrain/internal/port/httpserver/handlers"
	"vibrain/web"

	"github.com/labstack/echo/v4"
)

func registerRouters(e *echo.Echo, handler *handlers.Handler) {
	v1Api := e.Group("/api/v1")

	tools := v1Api.Group("/tools")
	tools.GET("/web/reader", handler.WebReaderHandler)
	tools.POST("/web/reader", handler.WebReaderHandler)
	tools.GET("/web/search", handler.WebSearchHandler)
	tools.POST("/web/search", handler.WebSearchHandler)
	tools.GET("/web/summary", handler.WebSummaryHandler)
	tools.POST("/web/summary", handler.WebSummaryHandler)

	oauth := e.Group("/oauth")
	oauth.GET("/:provider/login", handler.OAuthLoginHandler)
	oauth.GET("/:provider/callback", handler.OAuthCallbackHandler)

	// Health check
	e.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	debugApi := e.Group("/debug")
	debugApi.GET("/routes", func(c echo.Context) error {
		routes := e.Routes()
		return handlers.JsonResponse(c, http.StatusOK, routes)
	})

	// static files
	e.GET("/*", echo.WrapHandler(http.FileServer(web.StaticHttpFS)))
}
