package main

import (
	"net/http"

	"github.com/pichuchen/hatsuaki/activitypub"
	"github.com/pichuchen/hatsuaki/webfinger"
)

func Route(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("GET /.well-known/webfinger", webfinger.Route)
	mux.HandleFunc("/.activitypub/", activitypub.Route)
	// mux.HandleFunc("/world", world)

	return mux
}
