package httpserver

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// Serve starts the HTTP server
func Serve() {
	e := echo.New()
	registerMiddlewares(e)
	registerRouters(e)

	// Health check
	e.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "1323"
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
