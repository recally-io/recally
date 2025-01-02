package web

import (
	"embed"
	"net/http"
	"os"
	"recally/internal/pkg/config"
	"recally/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

//go:embed all:dist
var StaticFiles embed.FS

var StaticHttpFS http.FileSystem

func init() {
	// dynamic file system
	if config.Settings.Debug {
		logger.Default.Debug("Using dynamic file system for web UI")
		StaticHttpFS = http.FS(os.DirFS("web/dist"))
		return
	}

	// static file system
	PublicDirFS := echo.MustSubFS(StaticFiles, "dist")
	StaticHttpFS = http.FS(PublicDirFS)
}
