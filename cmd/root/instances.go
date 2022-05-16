package main

import (
	"net"
)

func handleInstance(conn net.Conn) {
	instances = append(instances, &conn)


}
