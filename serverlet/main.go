package main

import (
	"log/slog"
	"net/http"

	"github.com/pichuchen/hatsuaki/datastore/actor"
)

// 這個檔案的用途是整個系統的最初進入點
// 包括聆聽 HTTP 埠口以及呼叫 router 的功能

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	actor.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		XForwardedFor := r.Header.Get("X-Forwarded-For")
		slog.Debug("main", "method", r.Method, "url", r.URL.String(), "remote", r.RemoteAddr, "X-Forwarded-For", XForwardedFor)
		mux2 := http.NewServeMux()
		mux2 = Route(mux2)
		mux2.ServeHTTP(w, r)
	})
	// 這裡的 Route 是在 route.go 中定義的函數

	// 在這邊已明文聆聽 HTTP 埠口 8080
	http.ListenAndServe("0.0.0.0:8083", mux)

}
