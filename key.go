package main

import (
	"encoding/json"
	"net/http"
)

func keyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "" {
			json.NewEncoder(w).Encode(values)
			return
		}
		w.Write([]byte(values[r.URL.Path]))
	case http.MethodPost:
		var value string
		json.NewDecoder(r.Body).Decode(&value)
		values[r.URL.Path] = value
		go notifyInstances(r.URL.Path, value)
	}
}
