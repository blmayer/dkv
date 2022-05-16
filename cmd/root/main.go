package main

import (
	"net"
	"net/http"
	"os"
	"time"
	"math/rand"
)

var (
	instances []*net.Conn
	keys    = map[string][]int{}
	ikeys   = map[int][]string{}
	rep     = 2
)

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	rand.Seed(time.Now().UnixNano())

	// three handlers: one for clients and two for other instances
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
	println("server ready on port 1234")
	for {
		conn, err := server.Accept()
		if err != nil {
			println(err.Error())
		}
		println(conn.RemoteAddr().String(), "connected")
		instances = append(instances, &conn)
	}
}
