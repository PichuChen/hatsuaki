package index

import (
	"html/template"
	"log/slog"
	"net/http"

	_ "embed"
)

//go:embed timeline.html
var TimelineTemplate string

func RouteIndex(w http.ResponseWriter, r *http.Request) {
	slog.Debug("web.RouteIndex", "request", r.URL.String())
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.New("timeline").Parse(TimelineTemplate)
	if err != nil {
		slog.Error("web.RouteIndex", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
