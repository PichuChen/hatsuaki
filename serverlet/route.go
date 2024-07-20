package main

import (
	"net/http"

	"github.com/pichuchen/hatsuaki/activitypub"
	"github.com/pichuchen/hatsuaki/api"
	"github.com/pichuchen/hatsuaki/fetcher"
	"github.com/pichuchen/hatsuaki/web"
	"github.com/pichuchen/hatsuaki/webfinger"
)

func Route(mux *http.ServeMux) *http.ServeMux {
	// 在 webfinger 裡面實作的主要是公開必要的 webfinger 資訊
	mux.HandleFunc("GET /.well-known/webfinger", webfinger.Route)

	// 在 .activitypub 裡面實作的主要是處理 activitypub 的請求
	mux.HandleFunc("/.activitypub/", activitypub.Route)

	// 在 web 裡面實作的主要是處理網頁的請求
	mux.HandleFunc("/", web.Route)

	// 在 /1/ 處理比較傳統的 API 呼叫，例如登入註冊以及本站的抓取時間軸等
	// 通常如果是既有網站的話，也許不需要另外實作這部分的 API
	mux.HandleFunc("/1/", api.Route)

	mux.HandleFunc("/fetcher", fetcher.Route)
	// mux.HandleFunc("/world", world)

	return mux
}
