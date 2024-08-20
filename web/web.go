package web

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed all:dist
var StaticFiles embed.FS

var StaticHttpFS http.FileSystem

func init() {
	// dynamic file system
	// StaticHttpFS = http.FS(os.DirFS("web/dist"))

	// static file system
	PublicDirFS := echo.MustSubFS(StaticFiles, "dist")
	StaticHttpFS = http.FS(PublicDirFS)
}
