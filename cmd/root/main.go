package main

import (
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	instances []*net.Conn

	// TODO: get rid of theses maps
	keys  = map[string][]int{}
	ikeys = map[int][]string{}
	rep   = 2
)

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", keyHandler)
	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			panic(err)
		}
	}()

	server, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	println("main: server ready on port 1234")
	for {
		conn, err := server.Accept()
		if err != nil {
			println(err.Error())
		}
		println("main:", conn.RemoteAddr().String(), "connected")
		instances = append(instances, &conn)
	}
}
