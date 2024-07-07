package fetcher

// fetcher 的用途是讓後端伺服器可以替前端取得其他伺服器的 JSON-LD 資料。
// 之所以不讓前端直接取得是因為 CORS 的問題。

import (
	"io"
	"net/http"
)

func Route(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		Get(w, r)
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// get 會傳入一個參數 url，然後回傳 url 的 JSON-LD 資料。
func Get(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 這邊很重要，如果不設定的畫，有些伺服器會回傳 HTML 資料。
	req.Header.Set("Accept", "application/activity+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
