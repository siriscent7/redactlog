package web

import (
	"encoding/json"
	"net/http"

	"github.com/siriscent7/redactlog/redactor"
)

type redactRequest struct {
	Text string `json:"text"`
	Mode string `json:"mode"`
}

type redactResponse struct {
	Redacted string `json:"redacted"`
	Count    int    `json:"count"`
}

// Handler returns the HTTP mux for the RedactLog web UI.
func Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))
	})

	mux.HandleFunc("/api/redact", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		var req redactRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		mode := redactor.Mask
		switch req.Mode {
		case "hash":
			mode = redactor.Hash
		case "drop":
			mode = redactor.Drop
		}
		out, n := redactor.New(mode).Redact(req.Text)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(redactResponse{Redacted: out, Count: n})
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(indexHTML))
	})

	return mux
}
