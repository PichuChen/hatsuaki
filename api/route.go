package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/pichuchen/hatsuaki/activitypub"
	"github.com/pichuchen/hatsuaki/datastore/actor"
	"github.com/pichuchen/hatsuaki/datastore/object"
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
	} else if r.URL.Path == "/1/login" {
		PostLogin(w, r)
		return
	} else if r.URL.Path == "/1/note" {
		PostNote(w, r)
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

func PostLogin(w http.ResponseWriter, r *http.Request) {
	slog.Info("api.PostLogin", "info", "login")
	r.ParseForm()
	username := r.FormValue("username")
	if username == "" {
		slog.Warn("api.PostLogin", "warn", "username is empty")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	password := r.FormValue("password")
	if password == "" {
		slog.Warn("api.PostLogin", "warn", "password is empty")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	_, err := actor.FindActorByUsername(username)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if err = actor.VerifyPassword(username, password); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := IssueJWT(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	m := map[string]interface{}{
		"success": true,
		"token":   token,
	}
	json.NewEncoder(w).Encode(m)

}

func PostNote(w http.ResponseWriter, r *http.Request) {
	slog.Info("api.PostPost", "info", "post")

	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := auth[len("Bearer "):]
	username, err := VerifyJWT(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	content := r.FormValue("content")
	if content == "" {
		slog.Warn("api.PostPost", "warn", "content is empty")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	a, err := actor.FindActorByUsername(username)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	o := object.NewNote()
	o.SetContent(content)
	o.SetAttributedTo(a.GetFullID())
	o.AddCC(a.GetFullID() + "/followers")
	o.AddTo("https://www.w3.org/ns/activitystreams#Public")
	a.AppendOutboxObject(o.GetFullID())

	activitypub.SendCreate(a, o)

	err = object.SaveObject("./object.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
