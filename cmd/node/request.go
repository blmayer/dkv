package main

import (
	"strings"

	"dkv/internal/status"
	"dkv/internal/op"
)
func handleRequest(data []byte) []byte {
	println("received", string(data))
	if len(data) == 0 {
		return []byte{status.ErrNoData}
	}

	switch data[0] {
	case op.Get:
		val, ok := values[string(data[1:])]
		if !ok {
			return []byte{status.ErrNoData}
		}
		return append([]byte{status.Ok}, []byte(val)...)
	case op.Post:
		kv := strings.Split(string(data[1:]), "\n")
		values[kv[0]] = kv[1]
		return []byte{status.Ok}
	case op.Delete:
		delete(values, string(data[1:]))
		return []byte{status.Ok}
	}
	return []byte{status.ErrUnknownOp}
}

