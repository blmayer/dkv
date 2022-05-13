package main

import (
	"encoding/json"
	"mime"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	instances []string
	addr      string
	values    = map[string]string{}
)

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	// get own ip
	addrs, err := net.LookupHost("localhost")
	if err != nil || len(addrs) < 1 {
		panic(err.Error())
	}
	addr = "http://" + addrs[0] + ":" + port

	// register itself on other instances
	if r := os.Getenv("DKV_ROOT"); r != "" {
		resp, err := http.Post(r+"/instances", mime.TypeByExtension(".json"), strings.NewReader(`"`+addr+`"`))
		if err != nil {
			println("failed to register on root node", err.Error())
		}
		json.NewDecoder(resp.Body).Decode(&instances)
		resp.Body.Close()
		println("discovered", len(instances), "instances")

		for _, i := range instances {
			_, err = http.Post(i+"/instances", mime.TypeByExtension(".json"), strings.NewReader(`"`+addr+`"`))
			if err != nil {
				println("failed to register on node", i, err.Error())
			}
		}
		instances = append(instances, r)

		resp, err = http.Get(r + "/key/")
		if err != nil {
			println("failed to connect to root node", err.Error())
		}
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&values)
		println("got", len(values), "from", r)
	}

	// three handlers: one for clients and two for other instances
	http.Handle("/instances", http.StripPrefix("/instances", http.HandlerFunc(instanceHandler)))
	http.Handle("/key/", http.StripPrefix("/key/", http.HandlerFunc(keyHandler)))
	http.Handle("/notify/", http.StripPrefix("/notify/", http.HandlerFunc(notificationHandler)))

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
