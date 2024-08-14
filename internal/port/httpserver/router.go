package httpserver

import (
	"net/http"
	"vibrain/internal/core/assistants"
	"vibrain/internal/core/workers"
	"vibrain/internal/pkg/auth"
	"vibrain/web"

	_ "vibrain/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//	@title			Vibrain API
//	@version		1.0
//	@description	This is a simple API for Vibrain project.
//	@termsOfService	https://vibrain.vaayne.com/terms/

//	@contact.name	Vaayne
//	@contact.url	https://vaayne.com
//	@contact.email	vibrain@vaayne.com

// @host		localhost:1323
// @BasePath	/api/v1
func (s *Service) registerRouters() {
	e := s.Server
	v1Api := e.Group("/api/v1")

	// Authentication API
	authApi := newAuthHandler(auth.New())
	authApi.Register(v1Api)

	// Assistant API
	assistantApi := newAssistantHandler(assistants.NewService(s.llm))
	assistantApi.Register(v1Api)

	// Tools API
	toolsApi := newToolsHandler(workers.New(s.cache))
	toolsApi.Register(v1Api)

	// Health check
	e.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Debug routes
	debugApi := e.Group("/debug")
	debugApi.GET("/routes", func(c echo.Context) error {
		routes := e.Routes()
		return JsonResponse(c, http.StatusOK, routes)
	})

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// web pages
	e.GET("/*", echo.WrapHandler(http.FileServer(web.StaticHttpFS)))
}
