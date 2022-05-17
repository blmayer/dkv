package main

import (
	"dkv/internal/op"
	"dkv/internal/status"
	"math/rand"
)

func chooseInstances(rep int) []int {
	is := []int{}
	avail := []int{}
	for i, c := range instances {
		if c != nil {
			avail = append(avail, i)
			println("chooseInstances available:", i)
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

func getCompanionInstances(i int, k string) []int {
	co := []int{}
	for _, c := range keys[k] {
		if c != i && instances[c] != nil {
			co = append(co, c)
		}
	}
	return co
}

func writeToInstance(i int, o byte, data []byte) ([]byte, error) {
	println("writeToInstance Write:", i, o, string(data))
	c := *instances[i]
	_, err := c.Write(append([]byte{o}, data...))
	if err != nil {
		return nil, err
	}

	// read response
	resp := make([]byte, 1024)
	n, err := c.Read(resp)
	if err != nil {
		return nil, err
	}
	println("writeToInstance Read:", i, o, string(resp[:n]))
	return resp[:n], nil
}

func postToInstances(rep int, k string, data []byte) {
	is := keys[k]
	if len(is) == 0 {
		// this is a new key
		is = chooseInstances(rep)
	}

	for _, i := range is {
		println("postToInstances: posting", i, "data:", string(data))

		resp, err := writeToInstance(i, op.Post, []byte(k+"\n"+string(data)))
		if err != nil {
			println("postToInstances: ", err.Error(), "removing")
			moveInstanceKeys(i)
			postToInstances(rep-1, k, data)
			return
		}

		println("postToInstances: write returned", string(resp))
		if len(resp) == 0 || resp[0] != status.Ok {
			continue
		}

		keys[k] = append(keys[k], i)
		ikeys[i] = append(ikeys[i], k)
		println("postToInstances: done", i)
	}
}
