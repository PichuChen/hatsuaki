package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/pichuchen/hatsuaki/datastore/actor"
	"github.com/pichuchen/hatsuaki/datastore/config"
)

// 這個檔案的用途是整個系統的最初進入點
// 包括聆聽 HTTP 埠口以及呼叫 router 的功能

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	var err error

	// 檢查 config.json 是否存在，如果不存在就建立一個新的
	// 如果存在就讀取進來
	err = config.LoadConfig("./config.json")
	if errors.Is(err, os.ErrNotExist) {
		slog.Info("main", "config", "config.json not found, creating a new one")
		err = config.SaveConfig("./config.json")
		if err != nil {
			slog.Error("main", "error", err)
		}
	} else if err != nil {
		slog.Error("main", "error", err)
	}

	if config.GetLoginJWTSecret() == "" {
		// 產生 256bit 的隨機字串
		key := make([]byte, 32)
		_, err = rand.Read(key)
		if err != nil {
			slog.Error("main", "error", err)
		}
		config.SetLoginJWTSecret(base64.StdEncoding.EncodeToString(key))
		err = config.SaveConfig("./config.json")
		if err != nil {
			slog.Error("main", "error", err)
		}
	}

	err = actor.LoadActor("./actor.json")
	if errors.Is(err, os.ErrNotExist) {
		slog.Info("main", "actor", "actor.json not found, creating a new one")
		actor.InitActorDatastore()
		err = actor.SaveActor("./actor.json")
		if err != nil {
			slog.Error("main", "error", err)
		}
	} else if err != nil {
		slog.Error("main", "error", err)
	}

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
