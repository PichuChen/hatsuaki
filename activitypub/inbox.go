package activitypub

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/pichuchen/hatsuaki/datastore/actor"
	"github.com/pichuchen/hatsuaki/datastore/config"
)

// 相關文件請參閱: https://www.w3.org/TR/activitypub/#inbox

// 這邊會接收所有 /.activitypub/actor/inbox 開頭的請求
// 舉例來說會像是 GET /.activitypub/actor/alice
func RouteActorInbox(w http.ResponseWriter, r *http.Request) {
	slog.Debug("activitypub.RouteActorInbox", "request", r.URL.String())

	username := r.PathValue("actor")
	a, err := actor.FindActorByUsername(username)
	if err != nil {
		slog.Warn("activitypub.RouteActor", "error", "actor not found")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "actor not found"})
		return
	}

	if r.Method == "GET" {
		// GET 的狀況通常是該使用者想要查看自己的 inbox
		GetActorInbox(w, r, a)
		return
	} else if r.Method == "POST" {
		// POST 的狀況通常是有人想要發送訊息給該使用者，值得一提的是。
		PostActorInbox(w, r, a)
		return
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

	}
}
func GetActorInbox(w http.ResponseWriter, r *http.Request, a *actor.Actor) {

	w.Header().Set("Content-Type", "application/activity+json")
	m := map[string]interface{}{}

	// 在 JSON-LD 的回應中分為兩個大部分，@context 和其他的
	// @context 理論上是必須，但是實際上實作中大家通常都不會去讀取他，所以比較偏向會給工程師除錯用的。
	// 另外如果在 JSON 中有新增自己站自定義的欄位時，請記得補充 context 內容。

	c := []interface{}{}

	// 這是必要的部分
	c = append(c, "https://www.w3.org/ns/activitystreams")
	c = append(c, "https://w3id.org/security/v1")
	m["@context"] = c

	baseURL := "https://" + config.GetDomain() + "/.activitypub/actor/" + a.GetUsername()

	// All objects must have an id and type property
	m["id"] = baseURL
	m["type"] = "Person"

	// 接下來是在 ActivityPub 中的必要 (MUST) 欄位
	m["inbox"] = baseURL + "/inbox"
	m["outbox"] = baseURL + "/outbox"

	// 這邊是在 ActivityPub 中的應該 (SHOULD) 欄位
	m["following"] = baseURL + "/following"
	m["followers"] = baseURL + "/followers"

	// 這邊是在 ActivityPub 中的也許 (MAY) 欄位
	m["liked"] = baseURL + "/liked"
	// m["streams"] = baseURL + "/streams"
	// 在 misskey 2024.05 之前的版本，沒有 perferredUsername 會造成更新錯誤。
	m["preferredUsername"] = a.GetUsername()

	endpoints := map[string]string{}
	// 有 sharedInbox 的話，可以講低同個 instance follow 同個外部使用者時的訊息量。
	// 另外在 misskey 2024.05 之前的版本，沒有 sharedInbox 會造成更新錯誤。
	endpoints["sharedInbox"] = "https://" + config.GetDomain() + "/.activitypub/inbox"
	m["endpoints"] = endpoints

	// 此處請依照喜好自由加入。
	// m["published"] = "2023-01-01T00:00:00Z"
	// m["icon"] = nil
	// m["image"] = nil
	// m["url"] = baseURL
	// m["name"] = a.GetUsername()
	// m["manuallyApprovesFollowers"] = false
	// m["discoverable"] = true
	// m["summary"] = "test"

	json.NewEncoder(w).Encode(m)
}

func PostActorInbox(w http.ResponseWriter, r *http.Request, a *actor.Actor) {
	slog.Info("activitypub.PostActorInbox", "info", "inbox")

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	w.Header().Set("Content-Type", "application/activity+json")

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	m := map[string]interface{}{}
	m["id"] = "https://" + config.GetDomain() + "/.activitypub/actor/" + a.GetUsername() + "/inbox"

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	m["type"] = "OrderedCollection"

	json.NewEncoder(w).Encode(m)
}
