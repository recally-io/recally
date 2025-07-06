package docs

import (
	"embed"
	"fmt"
	"net/http"
	"os"

	"recally/internal/pkg/config"

	"github.com/labstack/echo/v4"
)

const distDir = ".vitepress/dist"

//go:embed all:.vitepress/dist
var StaticFiles embed.FS

var StaticHttpFS http.FileSystem

func init() {
	// dynamic file system
	if config.Settings.Debug {
		StaticHttpFS = http.FS(os.DirFS(fmt.Sprintf("docs/%s", distDir)))

		return
	}

	// static file system
	PublicDirFS := echo.MustSubFS(StaticFiles, distDir)
	StaticHttpFS = http.FS(PublicDirFS)
}
