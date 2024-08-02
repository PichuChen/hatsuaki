package activitypub

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/pichuchen/hatsuaki/datastore/object"
)

func RouteObject(w http.ResponseWriter, r *http.Request) {
	slog.Debug("activitypub.RouteObject", "request", r.URL.String())

	objectID := r.PathValue("object")
	o, err := object.FindObjectByID(objectID)
	if err != nil {
		slog.Warn("activitypub.RouteObject", "error", "object not found")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "object not found"})
		return
	}

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

	// All objects must have an id and type property
	m["id"] = o.GetFullID()
	m["type"] = o.GetType()

	// 接下來是在 ActivityPub 中的必要 (MUST) 欄位
	m["attributedTo"] = o.GetAttributedTo()
	m["content"] = o.GetContent()

	// 這邊是在 ActivityPub 中的應該 (SHOULD) 欄位
	m["cc"] = o.GetCC()
	m["to"] = o.GetTo()
	m["published"] = o.GetPublished()

	// 這邊是在 ActivityPub 中的也許 (MAY) 欄位
	m["inReplyTo"] = o.GetInReplyTo()

	json.NewEncoder(w).Encode(m)
}
