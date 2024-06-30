package webfinger

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/pichuchen/hatsuaki/datastore/actor"
	"github.com/pichuchen/hatsuaki/datastore/config"
)

// 在這邊會以 GET /.well-known/webfinger 這個路徑來呼叫這個函數
// webfinger 的詳細定義在 [RFC7033](https://datatracker.ietf.org/doc/html/rfc7033)
// 他的請求範例會像是 GET /.well-known/webfinger?resource=acct%3Aalice%40example.com
func Route(w http.ResponseWriter, r *http.Request) {
	slog.Info("webfinger.Route", "request", r.URL.String())

	// 這邊的 Content-Type 是 application/jrd+json
	w.Header().Set("Content-Type", "application/jrd+json")

	resource := r.URL.Query().Get("resource")
	if resource == "" {
		slog.Warn("webfinger.Route", "error", "resource query parameter is required")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "resource query parameter is required"})
		return
	}

	if !strings.HasPrefix(resource, "acct:") {
		slog.Warn("webfinger.Route", "error", "resource query parameter must start with acct:")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "resource query parameter must start with acct:"})
		return
	}

	acct := strings.TrimPrefix(resource, "acct:")
	if acct == "" {
		slog.Warn("webfinger.Route", "error", "resource query parameter must have a value after acct:")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "resource query parameter must have a value after acct:"})
		return
	}

	// 有些人會是 @alice 這樣的 username，把他轉成 alice
	if strings.HasPrefix(acct, "@") {
		// remove the @
		slog.Info("webfinger.Route", "info", "remove @ prefix")
		acct = strings.TrimPrefix(acct, "@")
	}

	// @ 後面的部分是 domain，事實上可以再進行一次驗證，不過這邊不管他
	if strings.Contains(acct, "@") {
		// remove the domain part
		acct = strings.Split(acct, "@")[0]
	}

	a, err := actor.FindActorByUsername(acct)
	if err != nil {
		slog.Warn("webfinger.Route", "error", "actor not found")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "actor not found"})
		return
	}

	slog.Info("webfinger.Route", "actor", a.GetUsername())

	m := map[string]interface{}{}
	m["subject"] = "acct:" + a.GetUsername() + "@" + config.GetDomain()
	m["aliases"] = []string{"https://" + config.GetDomain() + "/u/" + a.GetUsername()}

	links := []map[string]string{}

	// 這邊是個人頁面的網址，給人類使用者看的
	links = append(links, map[string]string{
		"rel":  "http://webfinger.net/rel/profile-page",
		"type": "text/html",
		"href": "https://" + config.GetDomain() + "/u/" + a.GetUsername(),
	})

	// 這邊這個網址是給機器人使用的，可以直接取得 JSON 格式的資料
	links = append(links, map[string]string{
		"rel":  "self",
		"type": "application/ld+json",
		"href": "https://" + config.GetDomain() + "/.activitypub/actor/" + a.GetUsername(),
	})

	// 這邊是讓自家的使用者如果看到其他站的使用者可以直接點擊後由該站導回自家站進行後續訂閱手續的網址
	// 但是並不是 activitypub 的標準，所以這邊先註解掉
	// links = append(links, map[string]string{
	// 	"rel":  "http://ostatus.org/schema/1.0/subscribe",
	// 	"template": "https://" + config.GetDomain() + "/authorize_interaction=?url={url}"
	// })

	m["links"] = links
	json.NewEncoder(w).Encode(m)
}
