package main

import (
	"io"
	"net/http"
	"math/rand"

	"dkv/internal/status"
	"dkv/internal/op"
)

func keyHandler(w http.ResponseWriter, r *http.Request) {
	k := r.URL.Path
	println(r.Method, k)

	switch r.Method {
	case http.MethodGet:
		// get instances
		is := keys[k]
		if len(is) == 0 {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		for _, i := range is {
			println("requesting", i)
			
			c := *instances[i]
			_, err := c.Write(append([]byte{op.Get}, []byte(k)...))
			if err != nil {
				println(err.Error())
				instances[i] = nil
				// go removeInstance(i)
				continue
			}

			// read response
			println("reading response")
			data := make([]byte, 1024)
			n, err := c.Read(data)
			if err != nil {
				println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				continue
			}

			if n == 0 {
				http.Error(w, "no data from node", http.StatusInternalServerError)
				continue
			}
			println(i, "sent", string(data))
			if data[0] != status.Ok {
				http.Error(w, "node returned not ok", http.StatusInternalServerError)
				continue
			}

			w.Write(data[1:n])
			break
		}

	case http.MethodPost:
		is := chooseInstances(rep)
		data := make([]byte, 1024)
		n, err := r.Body.Read(data)
		if err != nil && err != io.EOF {
			println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Body.Close()
		println("data:", string(data))

		for _, i := range is {
			println("requesting", i)

			c := *instances[i]
			_, err := c.Write(append([]byte{op.Post}, []byte(k+"\n"+string(data[:n]))...))
			if err != nil {
				println(err.Error())
				instances[i] = nil
				// go removeInstance(i)
				continue
			}

			// read response
			resp := make([]byte, 1)
			_, err = c.Read(resp)
			if err != nil {
				println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			println("write returned", string(resp))
			if len(resp) == 0 || resp[0] != status.Ok {
				http.Error(w, "write error", http.StatusInternalServerError)
				return
			}

			// TODO: watch for key already used
			keys[k] = append(keys[k], i) 
			ikeys[i] = append(ikeys[i], k)
			println("done", i)
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func chooseInstances(n int) []int {
	is := []int{}
	for i := 0; i < n; i++ {
		is = append(is, rand.Intn(len(instances)))
	}
	return is
}

