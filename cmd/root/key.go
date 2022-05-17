package main

import (
	"io"
	"net/http"

	"dkv/internal/op"
	"dkv/internal/status"
)

func keyHandler(w http.ResponseWriter, r *http.Request) {
	k := r.URL.Path
	println(r.Method, k)

	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodPost:
		handlePost(w, r)
	case http.MethodDelete:
		handleDelete(w, r)
	}
	r.Body.Close()
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	k := r.URL.Path

	// get instances
	is := keys[k]
	if len(is) == 0 {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	for _, i := range is {
		println("handleGet: get for", i, k)
		resp, err := writeToInstance(i, op.Get, []byte(k))
		if err != nil {
			println(err.Error())
			go moveInstanceKeys(i)
			continue
		}

		w.Write(resp[1:])
		break
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	k := r.URL.Path

	data := make([]byte, 1024)
	n, err := r.Body.Read(data)
	if err != nil && err != io.EOF {
		println("handlePost read:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body.Close()
	println("handlePost: data:", string(data[:n]))
	postToInstances(rep, k, data[:n])

	w.WriteHeader(http.StatusCreated)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	k := r.URL.Path
	println("handleDelete: key:", k)

	is := keys[k]
	for _, i := range is {
		resp, err := writeToInstance(i, op.Delete, []byte(k))
		if err != nil {
			println("handleDelete: writeToInstance:", err.Error())
		}
		if resp[0] != status.Ok {
			println("handleDelete: status not ok")
		}

		// remove from inverted keys index
		go func(ins int) {
			ks := ikeys[ins]
			for u, key := range ks {
				if key == k {
					ks[u] = ks[len(ks)-1]
					ks = ks[:len(ks)-1]
					ikeys[ins] = ks
					break
				}
			}
		}(i)
	}

	delete(keys, k)

	w.WriteHeader(http.StatusNoContent)
}

func moveInstanceKeys(i int) {
	instances[i] = nil
	println("moveInstanceKeys: moving keys from", i)

	newIs := chooseInstances(rep - 1)
	for _, k := range ikeys[i] {
		comps := getCompanionInstances(i, k)
		if len(comps) == 0 {
			println("moveInstanceKeys: lost key", k)
			return
		}
		println("moveInstanceKeys: moving keys from", comps[0])
		println("moveInstanceKeys: requesting", comps[0])
		resp, err := writeToInstance(comps[0], op.Get, []byte(k))
		if err != nil {
			println("moveInstanceKeys: failed to get old key", k)
			continue
		}
		if len(resp) < 2 {
			println("moveInstanceKeys: no data for old key", k)
			continue
		}

		keys[k] = []int{}
		for _, n := range newIs {
			println("moving", k, "to", n)
			data := []byte(k + "\n" + string(resp[1:]))
			resp, err := writeToInstance(n, op.Post, data)
			if err != nil {
				println("moveInstanceKeys: failed to post old key", k)
				continue
			}
			if len(resp) == 0 || resp[0] != status.Ok {
				println("moveInstanceKeys: not ok for old key", k)
				continue
			}
			keys[k] = append(keys[k], n)
			ikeys[n] = append(ikeys[n], k)
		}

		ikeys[i] = []string{}
		println("moveInstanceKeys: moved key", k)
	}
}
