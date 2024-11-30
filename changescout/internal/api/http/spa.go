package http

import (
	"github.com/gelleson/changescout/changescout/pkg/ui/dist"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/fs"
	"net/http"
)

// RegisterHandlers sets up the required route handlers for serving static files and index.html
func RegisterHandlers(e *echo.Echo) {
	// Convert embed.FS to http.FileSystem
	fsys, err := fs.Sub(dist.DistFS, "dist")
	if err != nil {
		panic(err)
	}

	// Create a filesystem handler for static files
	fsHandler := middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       ".",
		Index:      "index.html",
		HTML5:      true,
		Filesystem: http.FS(fsys),
	})

	// Apply the filesystem handler
	e.Use(fsHandler)
}
