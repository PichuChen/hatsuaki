package activitypub

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/pichuchen/hatsuaki/datastore/actor"
	"github.com/pichuchen/hatsuaki/datastore/config"
	"github.com/pichuchen/hatsuaki/datastore/object"
)

// 相關文件請參閱: https://www.w3.org/TR/activitypub/#outbox

// 這邊會接收所有 /.activitypub/actor/{actor}/outbox 開頭的請求
// 舉例來說會像是 GET /.activitypub/actor/alice/outbox
// Outbox 的用途是讓沒有收到先前消息的人可以查詢某個 Actor 的所有發送過的消息。
func RouteActorOutbox(w http.ResponseWriter, r *http.Request) {

	// 在 Get 的部分標準中並沒有要求一定要驗證簽章
	// 然而在 Mastdon 的實作當中因為支援封鎖伺服器功能，
	// 因此需要鑑別 (Authentication) 來判斷是否有權限查看。

	username := r.PathValue("actor")
	a, err := actor.FindActorByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "actor not found"})
		return
	}

	page := r.URL.Query().Get("page")
	if page == "true" {
		// 如果有 page=true 的參數，則回傳一個 OrderedCollection
		// 這個 OrderedCollection 會包含所有的 Object
		RouteActorOutboxPage(w, r, a)
		return
	}

	w.Header().Set("Content-Type", "application/activity+json")
	m := map[string]interface{}{}

	c := []interface{}{}
	c = append(c, "https://www.w3.org/ns/activitystreams")
	c = append(c, "https://w3id.org/security/v1")
	m["@context"] = c

	id := "https://" + config.GetDomain() + "/.activitypub/actor/" + a.GetUsername() + "/outbox"

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	m["id"] = id
	m["type"] = "OrderedCollection"
	m["totalItems"] = a.GetOutboxObjectsCount()
	m["first"] = id + "?page=true"
	m["last"] = id + "?page=true"

	json.NewEncoder(w).Encode(m)

}

// RouteActorOutboxPage 會回傳一個 OrderedCollection
// 這個 OrderedCollection 會包含所有的 Object
func RouteActorOutboxPage(w http.ResponseWriter, r *http.Request, a *actor.Actor) {
	w.Header().Set("Content-Type", "application/activity+json")
	m := map[string]interface{}{}

	c := []interface{}{}
	c = append(c, "https://www.w3.org/ns/activitystreams")
	c = append(c, "https://w3id.org/security/v1")
	m["@context"] = c

	id := "https://" + config.GetDomain() + "/.activitypub/actor/" + a.GetUsername() + "/outbox"

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	m["id"] = id
	m["type"] = "OrderedCollection"
	m["totalItems"] = a.GetOutboxObjectsCount()
	m["first"] = id + "?page=true"
	m["last"] = id + "?page=true"

	objectIDs, err := a.GetOutboxObjects()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}
	orderedItems := []interface{}{}

	for _, oid := range objectIDs {
		o, err := object.FindObjectByID(oid)
		if err != nil {
			slog.Warn("activitypub.RouteActorOutboxPage", "error", err.Error())
			continue
		}
		activityMap := map[string]interface{}{}
		actor := "https://" + config.GetDomain() + "/.activitypub/actor/" + a.GetUsername()
		activityMap["id"] = o.GetFullID() + "/activity"
		activityMap["type"] = "Create"
		activityMap["published"] = o.GetPublished()
		activityMap["actor"] = actor

		objectMap := map[string]interface{}{}

		objectMap["id"] = o.GetFullID()
		objectMap["type"] = "Note"
		objectMap["published"] = o.GetPublished()
		objectMap["attributedTo"] = o.GetAttributedTo()
		objectMap["content"] = o.GetContent()

		activityMap["object"] = objectMap

		orderedItems = append(orderedItems, activityMap)
	}

	m["orderedItems"] = orderedItems

	json.NewEncoder(w).Encode(m)
}
