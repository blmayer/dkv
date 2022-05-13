package main

import (
	"encoding/json"
	"mime"
	"net/http"
	"strings"
)

func notificationHandler(w http.ResponseWriter, r *http.Request) {
	println("got notification for", r.URL.Path)
	var value string
	json.NewDecoder(r.Body).Decode(&value)
	values[r.URL.Path] = value
}

func notifyInstances(key, value string) {
	for _, i := range instances {
		println("sending", key, value, "to", i)
		go http.Post(i+"/notify/"+key, mime.TypeByExtension(".json"), strings.NewReader(`"`+value+`"`))
	}
}
