package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	expand "github.com/openvenues/gopostal/expand"
)

type Expansion struct {
	Address    string   `json:"address"`
	Expansions []string `json:"expansions"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "libpostal rest wrapper")
	})

	http.HandleFunc("/expand", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprintf(w, "expand addresses\nsend a list of addresses to expand in post body")
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var addresses []string
		err := json.NewDecoder(r.Body).Decode(&addresses)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		slog.Debug("expanding addresses", "addresses", addresses)
		expansions := make([]Expansion, len(addresses))

		for i, str := range addresses {
			expanded := expand.ExpandAddress(str)
			expansions[i] = Expansion{Address: str, Expansions: expanded}
			slog.Debug("expanded", "expansions", expansions[i])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expansions)
	})

	slog.Info("starting server", "port", "8080")
	http.ListenAndServe(":8080", nil)
}
