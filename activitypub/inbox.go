package activitypub

import (
	"encoding/json"
	"fmt"
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
		// POST 的狀況通常是有人想要發送訊息給該使用者。
		PostActorInbox(w, r, a)
		return
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

	}
}

func RouteSharedInbox(w http.ResponseWriter, r *http.Request) {
	slog.Debug("activitypub.RouteSharedInbox", "request", r.URL.String())

	if r.Method == "POST" {
		// POST 的狀況通常是有人想要發送訊息給該使用者。
		PostSharedInbox(w, r)
		return
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetActorInbox(w http.ResponseWriter, r *http.Request, a *actor.Actor) {

	w.Header().Set("Content-Type", "application/activity+json")
	m := map[string]interface{}{}

	// 這個部分需要驗證進行取 inbox 的人必須要是該使用者

	page := r.URL.Query().Get("page")
	if page == "true" {
		// 如果有 page=true 的參數，則回傳一個 OrderedCollection
		// 這個 OrderedCollection 會包含所有的 Object
		RouteActoInboxPage(w, r, a)
		return
	}

	c := []interface{}{}

	// 這是必要的部分
	c = append(c, "https://www.w3.org/ns/activitystreams")
	c = append(c, "https://w3id.org/security/v1")
	m["@context"] = c

	id := "https://" + config.GetDomain() + "/.activitypub/actor/" + a.GetUsername() + "/inbox"

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	m["id"] = id
	m["type"] = "OrderedCollection"
	m["totalItems"] = a.GetInboxObjectsCount()
	m["first"] = id + "?page=true"
	m["last"] = id + "?page=true"

	json.NewEncoder(w).Encode(m)
}

func RouteActoInboxPage(w http.ResponseWriter, r *http.Request, a *actor.Actor) {
	w.Header().Set("Content-Type", "application/activity+json")
	m := map[string]interface{}{}

	c := []interface{}{}
	c = append(c, "https://www.w3.org/ns/activitystreams")
	c = append(c, "https://w3id.org/security/v1")
	m["@context"] = c

	id := "https://" + config.GetDomain() + "/.activitypub/actor/" + a.GetUsername() + "/inbox"

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	m["id"] = id
	m["type"] = "OrderedCollection"
	m["totalItems"] = a.GetInboxObjectsCount()
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
		orderedItems = append(orderedItems, map[string]interface{}{"id": oid})
	}
	m["orderedItems"] = orderedItems

	json.NewEncoder(w).Encode(m)

}

func PostActorInbox(w http.ResponseWriter, r *http.Request, a *actor.Actor) {
	slog.Info("activitypub.PostActorInbox", "info", "inbox")

	// 解碼送入的 JSON
	var requestMap map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestMap)
	if err != nil {
		slog.Warn("activitypub.PostActorInbox", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "bad request"})
		return
	}

	requestType := requestMap["type"].(string)
	if requestType == "Follow" {
		PostActorInboxFollow(w, r, a, requestMap)
		return
	}

	slog.Debug("activitypub.PostActorInbox", "info", requestMap)

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	w.Header().Set("Content-Type", "application/activity+json")

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	m := map[string]interface{}{}
	m["id"] = "https://" + config.GetDomain() + "/.activitypub/actor/" + a.GetUsername() + "/inbox"

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	m["type"] = "OrderedCollection"

	json.NewEncoder(w).Encode(m)
}

func PostActorInboxFollow(w http.ResponseWriter, r *http.Request, a *actor.Actor, requestMap map[string]interface{}) {
	slog.Info("activitypub.PostActorInboxFollow", "info", "follow")

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	w.Header().Set("Content-Type", "application/activity+json")

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	m := map[string]interface{}{}
	m["id"] = "https://" + config.GetDomain() + "/.activitypub/actor/" + a.GetUsername() + "/inbox"

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	m["type"] = "OrderedCollection"

	json.NewEncoder(w).Encode(m)
}

func PostSharedInbox(w http.ResponseWriter, r *http.Request) {
	slog.Info("activitypub.PostSharedInbox", "info", "shared inbox")

	// 解碼送入的 JSON
	var requestMap map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestMap)
	if err != nil {
		slog.Warn("activitypub.PostSharedInbox", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "bad request"})
		return
	}

	requestType := requestMap["type"].(string)
	if requestType == "Create" {
		PostSharedInboxCreate(w, r, requestMap)
		return
	}

	slog.Debug("activitypub.PostSharedInbox", "info", requestMap)

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	w.Header().Set("Content-Type", "application/activity+json")

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	m := map[string]interface{}{}
	m["id"] = "https://" + config.GetDomain() + "/.activitypub/inbox"

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	m["type"] = "OrderedCollection"

	json.NewEncoder(w).Encode(m)
}

func PostSharedInboxCreate(w http.ResponseWriter, r *http.Request, requestMap map[string]interface{}) {
	slog.Info("activitypub.PostSharedInboxCreate", "info", "create", "requestMap", requestMap)

	encoded, _ := json.MarshalIndent(requestMap, "", "  ")
	slog.Debug("activitypub.PostSharedInboxCreate", "info", "create", "encoded", string(encoded))
	fmt.Printf("%s\n", string(encoded))

	o := requestMap["object"].(map[string]interface{})
	oid := o["id"].(string)
	// 這邊需要驗證 oid 的 id 是否和簽署的 key 的 domain 相同
	// 在這邊的驗證我們沒辦法信任來源 IP, 能信任的只有簽發的 Key 而已。

	prefix := "https://" + config.GetDomain() + "/.activitypub/actor/"
	toList := requestMap["to"].([]interface{})
	for _, v := range toList {
		to := v.(string)
		if to[:len(prefix)] != prefix {
			slog.Warn("activitypub.PostSharedInboxCreate", "skip", "to", "to", to)
			continue
		}
		actorName := to[len(prefix):]
		a, err := actor.FindActorByUsername(actorName)
		if err != nil {
			slog.Warn("activitypub.PostSharedInboxCreate", "error", "actor not found", "actorName", actorName)
			continue
		}

		a.AppendInboxObject(oid)
	}

	err := actor.SaveActor("./actor.json")
	if err != nil {
		slog.Warn("activitypub.PostSharedInboxCreate", "error", "actor save error", "err", err)
	}

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	w.Header().Set("Content-Type", "application/activity+json")

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	m := map[string]interface{}{}
	m["id"] = "https://" + config.GetDomain() + "/.activitypub/inbox"

	// 這邊是在 ActivityPub 中的必要 (MUST) 欄位
	m["type"] = "OrderedCollection"

	json.NewEncoder(w).Encode(m)
}
