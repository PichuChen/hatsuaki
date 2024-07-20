package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/pichuchen/hatsuaki/datastore/actor"
)

func Route(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		RouteGet(w, r)
		return
	} else if r.Method == "POST" {
		RoutePost(w, r)
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

}

func RouteGet(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func RoutePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/1/register" {
		PostRegister(w, r)
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func PostRegister(w http.ResponseWriter, r *http.Request) {
	slog.Info("api.PostRegister", "info", "register")
	r.ParseForm()
	// 先確定使用者名稱是否存在
	username := r.FormValue("username")
	if username == "" {
		slog.Warn("api.PostRegister", "warn", "username is empty")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	password := r.FormValue("password")
	if password == "" {
		slog.Warn("api.PostRegister", "warn", "password is empty")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	_, err := actor.FindActorByUsername(username)
	if err == nil {
		http.Error(w, "Conflict", http.StatusConflict)
		return
	}

	// 如果不存在，就建立一個新的 actor
	actor.NewActor(username)
	actor.UpdatePassword(username, password)

	// 儲存 actor
	err = actor.SaveActor("./actor.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	m := map[string]interface{}{
		"success": true,
	}
	json.NewEncoder(w).Encode(m)
}
