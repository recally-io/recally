package httpserver

import (
	"net/http"
	"recally/internal/pkg/logger"
	"recally/web"

	_ "recally/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//	@title			Recally API
//	@version		1.0
//	@description	This is a simple API for Recally project.
//	@termsOfService	https://recally.vaayne.com/terms/

//	@contact.name	Vaayne
//	@contact.url	https://vaayne.com
//	@contact.email	recally@vaayne.com

// @host		localhost:1323
// @BasePath	/api/v1
func (s *Service) registerRouters() {
	e := s.Server
	v1Api := e.Group("/api/v1")

	registerAuthHandlers(v1Api)
	registerAssistantHandlers(v1Api, s)
	registerToolsHandlers(v1Api, s)
	registerFileHandlers(v1Api, s.s3)
	registerBookmarkHandlers(v1Api, s)

	// Health check
	e.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Debug routes
	debugApi := e.Group("/debug", authAdminMiddleware())
	debugApi.GET("/routes", func(c echo.Context) error {
		routes := e.Routes()
		return JsonResponse(c, http.StatusOK, routes)
	})

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// web pages
	logger.Default.Debug("Using static files as frontend")
	e.GET("/manifest.webmanifest", func(c echo.Context) error {
		file, err := web.StaticHttpFS.Open("manifest.webmanifest")
		if err != nil {
			return err
		}
		defer file.Close()
		c.Response().Header().Set("Content-Type", "application/manifest+json")
		return c.Stream(http.StatusOK, "application/manifest+json", file)
	})
	e.GET("/*", echo.WrapHandler(http.FileServer(web.StaticHttpFS)))
}
