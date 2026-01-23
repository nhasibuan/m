package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"Host":    r.Host,
			"Content": w.Header().Get("Content-Type"),
		})
	})
	e := http.ListenAndServe(":8080", nil)
	if e != nil {
		log.Fatal(e)
	}
}
