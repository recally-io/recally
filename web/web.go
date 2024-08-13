package web

import (
	"embed"
	"net/http"
	"os"
	"vibrain/internal/pkg/config"

	"github.com/labstack/echo/v4"
)

//go:embed all:dist
var StaticFiles embed.FS

var StaticHttpFS http.FileSystem

func init() {
	if config.Settings.Debug {
		StaticHttpFS = http.FS(os.DirFS("web/dist"))
	} else {
		PublicDirFS := echo.MustSubFS(StaticFiles, "dist")
		StaticHttpFS = http.FS(PublicDirFS)
	}
}
