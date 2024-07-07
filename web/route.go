package web

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/pichuchen/hatsuaki/web/index"
)

// Web 套件主要是用來呈現給瀏覽器使用的 HTML
func Route(w http.ResponseWriter, r *http.Request) {
	slog.Debug("web.Route", "request", r.URL.String())
	mux := http.NewServeMux()
	mux.HandleFunc("GET /assets/", func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("web.Route.assets", "request", r.URL.String())

		path := strings.TrimPrefix(r.URL.Path, "/assets/")

		http.ServeFile(w, r, "../web/assets/"+path)
	})
	mux.HandleFunc("GET /", index.RouteIndex)

	mux.ServeHTTP(w, r)
}
