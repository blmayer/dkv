package main

import (
	"io"
	"math/rand"
	"net"
	"os"
	"time"
)

var (
	me string

	// first part is 8 bytes time and then raw data
	values = map[string][]byte{}
	factor = 4
)

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	rand.Seed(time.Now().UnixNano())

	server, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	println("main: server ready on port ", port)
	me = server.Addr().String()

	for {
		conn, err := server.Accept()
		if err != nil {
			println("main: conn accept", err.Error())
			continue
		}
		println("main:", conn.RemoteAddr().String(), "connected")

		go func() {
			data := make([]byte, 1024)
			n, err := conn.Read(data)
			if err == net.ErrClosed || err == io.EOF {
				println("main: connection closed from", conn.RemoteAddr().String())
				return
			} else if err != nil {
				println("main: read error:", err.Error())
				return
			}

			resp := handleRequest(data[:n])
			conn.Write(resp)
			conn.Close()

			println("mail: handler returned", string(resp))
		}()
	}
}
