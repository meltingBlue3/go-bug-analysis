package server

import (
	"io/fs"
	"net/http"
)

// New creates an HTTP handler with static file serving and API route placeholders.
// staticFS should be the embedded filesystem with the "static" prefix already stripped.
func New(staticFS fs.FS) http.Handler {
	mux := http.NewServeMux()

	// Static file server — serves HTML/CSS/JS/ECharts from embedded FS
	fileServer := http.FileServer(http.FS(staticFS))
	mux.Handle("/", fileServer)

	// API routes placeholder — Plan 02 will add /api/upload and other endpoints
	// mux.HandleFunc("POST /api/upload", uploadHandler)

	return mux
}
