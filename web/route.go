package web

import (
	"embed"
	"log/slog"
	"net/http"

	"github.com/pichuchen/hatsuaki/web/index"
)

//go:embed assets
var assetsFS embed.FS

// Web 套件主要是用來呈現給瀏覽器使用的 HTML
func Route(w http.ResponseWriter, r *http.Request) {
	slog.Debug("web.Route", "request", r.URL.String())
	mux := http.NewServeMux()
	mux.HandleFunc("GET /assets/", func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("web.Route.assets", "request", r.URL.String())

		http.FileServer(http.FS(assetsFS)).ServeHTTP(w, r)
	})
	mux.HandleFunc("GET /", index.RouteIndex)

	mux.ServeHTTP(w, r)
}
