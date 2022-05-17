package main

import (
	"io"
	"net"
	"os"
)

var (
	values = map[string]string{}
)

func main() {
	// register itself on root
	r := os.Getenv("DKV_ROOT")
	if r == "" {
		r = "localhost:1234"
	}

	conn, err := net.Dial("tcp", r)
	if err != nil {
		panic(err.Error())
	}

	for {
		data := make([]byte, 1024)
		n, err := conn.Read(data)

		if err == net.ErrClosed || err == io.EOF {
			println("connection closed from root")
			return
		} else if err != nil {
			println(err.Error())
			continue
		}

		go func(d []byte) {
			resp := handleRequest(d)
			conn.Write(resp)
			println("returned", string(resp))
		}(data[:n])
	}
}
