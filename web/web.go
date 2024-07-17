package web

import (
	"embed"
	"net/http"
	"os"
	"vibrain/internal/pkg/config"

	"github.com/labstack/echo/v4"
)

//go:embed all:public
var StaticFiles embed.FS

var StaticHttpFS http.FileSystem

func init() {
	if config.Settings.Debug {
		StaticHttpFS = http.FS(os.DirFS("web/public"))
	} else {
		PublicDirFS := echo.MustSubFS(StaticFiles, "public")
		StaticHttpFS = http.FS(PublicDirFS)
	}
}
