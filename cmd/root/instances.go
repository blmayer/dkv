package main

import (
	"math/rand"
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

