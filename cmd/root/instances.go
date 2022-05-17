package main

import (
	"dkv/internal/op"
	"dkv/internal/status"
	"math/rand"
	"net"
)

func chooseInstances(rep int) []int {
	is := []int{}
	avail := []int{}
	for i, c := range instances {
		if c != nil {
			avail = append(avail, i)
			println("available:", i)
		}
	}
	if len(avail) == 1 {
		return []int{avail[0]}
	}

	rand.Shuffle(len(avail), func(i, j int) {
		avail[i], avail[j] = avail[j], avail[i]
	})
	for i := 0; i < rep; i++ {
		is = append(is, avail[i])
	}
	return is
}

func moveInstanceKeys(i int, key string) {
	instances[i] = nil

	comps := getCompanionInstances(i, key)
	if len(comps) == 0 {
		println("lost keys from", i)
		return
	}
	println("moving keys from", comps[0])

	from := instances[comps[0]]
	newIs := chooseInstances(rep - 1)

	for _, k := range ikeys[i] {
		resp, err := writeToInstance(*from, op.Get, []byte(k))
		if err != nil {
			println("failed to get old key", k)
			continue
		}
		if len(resp) < 2 {
			println("no data for old key", k)
			continue
		}

		keys[k] = []int{}
		for _, n := range newIs {
			data := []byte(k + "\n" + string(resp[1:]))
			resp, err := writeToInstance(*instances[n], op.Post, data)
			if err != nil {
				println("failed to post old key", k)
				continue
			}
			if len(resp) == 0 || resp[0] != status.Ok {
				println("no t ok for old key", k)
				continue
			}
			keys[k] = append(keys[k], n)
			ikeys[n] = append(ikeys[n], k)
		}

		// TODO: remove old conn from ikeys and to the key
		println("moved key", k)
	}
}

func getCompanionInstances(i int, k string) []int {
	co := []int{}
	for _, c := range keys[k] {
		if c != i {
			co = append(co, c)
		}
	}
	return co
}

func writeToInstance(c net.Conn, o byte, data []byte) ([]byte, error) {
	_, err := c.Write(append([]byte{o}, data...))
	if err != nil {
		return nil, err
	}

	// read response
	resp := make([]byte, 1)
	_, err = c.Read(resp)
	if err != nil {
		return nil, err
	}
	return data, nil
}
