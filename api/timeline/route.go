package timeline

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/pichuchen/hatsuaki/activitypub"
	"github.com/pichuchen/hatsuaki/api/auth"
	"github.com/pichuchen/hatsuaki/datastore/actor"
)

func RouteTimeline(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		GetTimeline(w, r)
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// GetTimeline 會根據送進來的參數回傳相對應的 Timeline
func GetTimeline(w http.ResponseWriter, r *http.Request) {

	authHdr := r.Header.Get("Authorization")
	if authHdr == "" {
		slog.Warn("api.GetTimeline", "warn", "no authorization header")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := authHdr[len("Bearer "):]
	username, err := auth.VerifyJWT(token)
	if err != nil {
		slog.Warn("api.GetTimeline", "warn", "invalid token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	a, err := actor.FindActorByUsername(username)
	if err != nil {
		slog.Warn("api.GetTimeline", "warn", "actor not found")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if a == nil {
		slog.Warn("api.GetTimeline", "warn", "actor not found")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/activity+json")
	// 這版的 Timeline 演算法先只採用 Inbox + Outbox 然後已發布時間排序的方式

	idList := []string{}
	inboxIDs, err := a.GetInboxObjects()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	outboxIDs, err := a.GetOutboxObjects()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	idList = append(idList, inboxIDs...)
	idList = append(idList, outboxIDs...)

	list := make([]interface{}, len(idList))
	wg := sync.WaitGroup{}
	wg.Add(1)
	// 這邊先不考慮分頁
	for i, id := range idList {
		ii := i
		iid := id
		wg.Add(1)
		go func() {
			defer wg.Done()
			o, err := activitypub.FetchObject(iid, username, false)
			if err != nil {
				slog.Warn("api.GetTimeline.FetchObject", "id", id, "error", err.Error())
				return
			}
			list[ii] = o
		}()
	}
	wg.Done()
	wg.Wait()
	// clean up nil
	cleanList := []interface{}{}
	for _, v := range list {
		if v != nil {
			cleanList = append(cleanList, v)
		}
	}
	list = cleanList

	createActivityList := []interface{}{}
	for _, v := range list {
		srcObj := v.(map[string]interface{})
		o := map[string]interface{}{}
		o["type"] = "Create"
		o["actor"] = a.GetFullID()
		o["published"] = srcObj["published"]
		o["object"] = v
		createActivityList = append(createActivityList, o)

	}
	list = createActivityList

	// sort by published
	// 這邊先不考慮分頁
	sort.Slice(list, func(i, j int) bool {
		iTime, err := time.Parse(time.RFC3339, list[i].(map[string]interface{})["published"].(string))
		if err != nil {
			return false
		}
		jTime, err := time.Parse(time.RFC3339, list[j].(map[string]interface{})["published"].(string))
		if err != nil {
			return false
		}
		return iTime.After(jTime)
	})

	m := map[string]interface{}{}
	c := []interface{}{}
	c = append(c, "https://www.w3.org/ns/activitystreams")
	c = append(c, "https://w3id.org/security/v1")
	m["@context"] = c

	id := "https://hatsuaki-dev2.pichuchen.tw/1/timeline"
	m["id"] = id
	m["type"] = "OrderedCollection"
	m["totalItems"] = len(list)
	m["orderedItems"] = list

	json.NewEncoder(w).Encode(m)

}
