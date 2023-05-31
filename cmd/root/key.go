package main

import (
	"encoding/binary"
	"net"
	"os"
	"sort"
	"strings"
	"time"

	"dkv/internal/op"
	"dkv/internal/status"
)

func handleRequest(data []byte) []byte {
	println("handler: received", string(data))
	if len(data) == 0 {
		return []byte{status.ErrNoData}
	}

	var res []byte
	switch data[0] {
	case op.Get:
		val, ok := values[string(data[1:])]
		if !ok {
			// search other instances using the internal get
			res = distribute(data)
			break
		}
		res = append([]byte{status.Ok}, val[8:]...)
	case op.Post:
		kv := strings.Split(string(data[1:]), "\n")
		values[kv[0]] = append(
			binary.BigEndian.AppendUint64(make([]byte, 8), uint64(time.Now().Unix())),
			kv[1]...,
		)
		res = []byte{status.Ok}
		go distribute(data)
	case op.Delete:
		delete(values, string(data[1:]))
		res = []byte{status.Ok}
	case op.GetInternal:
		val, ok := values[string(data[1:])]
		if !ok {
			res = []byte{status.ErrNoData, 0, 0, 0, 0, 0, 0, 0, 0}
			break
		}
		res = append([]byte{status.Ok}, val...)
	case op.PostInternal:
		kv := strings.Split(string(data[1:]), "\n")
		values[kv[0]] = append(
			binary.BigEndian.AppendUint64(make([]byte, 8), uint64(time.Now().Unix())),
			data[1:]...,
		)
		res = []byte{status.Ok}
	case op.DeleteInternal:
		delete(values, string(data[1:]))
		res = []byte{status.Ok}
	default:
		res = []byte{status.ErrUnknownOp}
	}

	return res
}

func distribute(data []byte) []byte {
	println("distribute: start")

	env := os.Getenv("DKV_INSTANCES")
	ins := strings.Split(env, ",")
	if len(ins) < 3 {
		println("distribute: not enough instances")
		return []byte{status.ErrNoData}
	}

	for i, addr := range ins {
		if me == addr {
			ins[i] = ins[len(ins)-1]
			ins = ins[:len(ins)-1]
			break
		}
	}

	// at least 2 instances of replication
	n := len(ins) / factor
	if n < 2 {
		n = 2
	}
	ins = ins[:n-1]

	switch data[0] {
	case op.Get:
		data[0] = op.GetInternal
	case op.Post:
		data[0] = op.PostInternal
	case op.Delete:
		data[0] = op.DeleteInternal
	}

	// TODO: make async and random
	res := make([][]byte, len(ins))
	for k, i := range ins {
		conn, err := net.Dial("tcp", i)
		if err != nil {
			println("distribute:", err.Error())
			continue
		}

		res[k] = make([]byte, 1024)
		go conn.Write(data)

		bs, err := conn.Read(res[k])
		if err != nil {
			println("distribute:", i, "returned", err.Error())
			continue
		}
		res[k] = res[k][:bs-1]
		conn.Close()
	}

	// check responses
	sort.Slice(
		res,
		func(i, j int) bool {
			ti := binary.BigEndian.Uint64(res[i][1:9])
			tj := binary.BigEndian.Uint64(res[j][1:9])
			return ti > tj
		},
	)

	return res[0]
}
