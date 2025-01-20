package httpserver

import (
	"net/http"
	"recally/docs"
	"recally/web"

	_ "recally/docs/swagger"

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
	registerFileHandlers(v1Api, s)
	registerBookmarkHandlers(v1Api, s)
	registerBookmarkShareHandlers(v1Api, s)
	registerLLMHandlers(v1Api, s)
	registerUsersHandlers(v1Api)

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

	// Docs
	e.GET("/docs/*", func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=86400") // Cache for 1 day
		return echo.WrapHandler(http.StripPrefix("/docs", http.FileServer(docs.StaticHttpFS)))(c)
	})

	// Web UI
	s.registerWebUIRouters(e)
}

func (s *Service) registerWebUIRouters(e *echo.Echo) {
	// manifest.webmanifest for PWA
	e.GET("/manifest.webmanifest", func(c echo.Context) error {
		file, err := web.StaticHttpFS.Open("manifest.webmanifest")
		if err != nil {
			return err
		}
		defer file.Close()
		c.Response().Header().Set("Content-Type", "application/manifest+json")
		c.Response().Header().Set("Cache-Control", "public, max-age=86400") // Cache for 1 day
		return c.Stream(http.StatusOK, "application/manifest+json", file)
	})

	e.GET("/assets/*", func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=86400") // Cache for 1 day
		return echo.WrapHandler(http.FileServer(web.StaticHttpFS))(c)
	})

	// Serve static files for SPA, if there is a static file, serve it, otherwise serve index.html
	e.GET("/*", func(c echo.Context) error {
		path := c.Param("*")
		if path != "" {
			// Try to open the requested file
			if file, err := web.StaticHttpFS.Open(path); err == nil {
				defer file.Close()
				// Add cache control for static assets
				return echo.WrapHandler(http.FileServer(web.StaticHttpFS))(c)
			}
		}

		// If file not found, serve index.html
		file, err := web.StaticHttpFS.Open("index.html")
		if err != nil {
			return err
		}
		defer file.Close()
		// No cache for index.html to ensure latest version
		c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		return c.Stream(http.StatusOK, "text/html", file)
	})
}
