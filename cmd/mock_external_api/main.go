package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func main() {
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		group := r.URL.Query().Get("group")
		song := r.URL.Query().Get("song")

		if group == "" || song == "" {
			http.Error(w, "missing group or song parameter", http.StatusBadRequest)
			return
		}

		response := SongDetail{
			ReleaseDate: "12.02.2005",
			Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
			Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("Mock external API server running on http://localhost:9090")
	http.ListenAndServe(":9090", nil)
}
