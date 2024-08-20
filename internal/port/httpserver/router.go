package httpserver

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/logger"
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

	registerAuthHandlers(v1Api)
	registerAssistantHandlers(v1Api, s)
	registerToolsHandlers(v1Api, s)

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
	if config.Settings.Env == "dev" {
		// proxy to vite server localhost:5173
		logger.Default.Debug("Using vite server as frontend")
		e.GET("/*", reactReverseProxy)
	} else {
		logger.Default.Debug("Using static files as frontend")
		e.GET("/*", echo.WrapHandler(http.FileServer(web.StaticHttpFS)))
	}
}

// reactReverseProxy is a reverse proxy for vite server
func reactReverseProxy(c echo.Context) error {
	remote, _ := url.Parse("http://localhost:5173")
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request().Header
		req.Host = remote.Host
		req.URL = c.Request().URL
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
	}
	proxy.ServeHTTP(c.Response().Writer, c.Request())
	return nil
}
