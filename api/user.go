package api

import (
	"encoding/json"
	"net/http"

	"github.com/pichuchen/hatsuaki/datastore/actor"
)

// 這個部分是傳統的取得使用者資料相關的後端
// 在 ActivityPub 中的 User 通常是 Actor
// 但是在一般的慣常設計中，我們會用 User 來代表使用者
// 所以為了方便理解，這邊我們用 User。
func RouteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		RouteUserGet(w, r)
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

}

func RouteUserGet(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path[len("/1/user/"):]
	u, err := actor.FindActorByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}
