package httpserver

import (
	"vibrain/internal/port/httpserver/handlers"

	"github.com/labstack/echo/v4"
)

func registerRouters(e *echo.Echo) {
	v1Api := e.Group("/api/v1")
	apiToolsRouters(v1Api)
	authRouters(e.Group("/oauth"))
}

func apiToolsRouters(e *echo.Group) {
	tools := e.Group("/tools")

	tools.GET("/web/reader", handlers.WebReaderHandler)
	tools.POST("/web/reader", handlers.WebReaderHandler)
	tools.GET("/web/search", handlers.WebSearchHandler)
	tools.POST("/web/search", handlers.WebSearchHandler)
	tools.GET("/web/summary", handlers.WebSummaryHandler)
	tools.POST("/web/summary", handlers.WebSummaryHandler)
}

func authRouters(oauth *echo.Group) {
	oauth.GET("/:provider/login", handlers.OAuthLoginHandler)
	oauth.GET("/:provider/callback", handlers.OAuthCallbackHandler)
}
