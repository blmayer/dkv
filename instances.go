package main

import (
	"encoding/json"
	"net/http"
)

func instanceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(instances)
	case http.MethodPost:
		// update instances list
		var addr string
		json.NewDecoder(r.Body).Decode(&addr)
		r.Body.Close()

		println("got new instance", addr)
		json.NewEncoder(w).Encode(instances)
		instances = append(instances, addr)
	case http.MethodDelete:
		var addr string
		json.NewDecoder(r.Body).Decode(&addr)
		r.Body.Close()

		for i, a := range instances {
			if a == addr {
				instances[i] = instances[len(instances)-1]
				instances = instances[:len(instances)-1]
			}
		}
		println("lost instance", addr)
		json.NewEncoder(w).Encode(instances)
	}
}
