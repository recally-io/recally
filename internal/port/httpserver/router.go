package httpserver

import (
	"vibrain/internal/port/httpserver/handlers"

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
}
