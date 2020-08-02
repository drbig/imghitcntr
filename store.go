package main

import (
	"sync"
)

type localStore struct {
	mu   sync.Mutex
	data map[string]int
}

var (
	db localStore
)

func init() {
	db.data = make(map[string]int, 32)
}

func getCount(key string) int {
	db.mu.Lock()
	defer db.mu.Unlock()

	count, ok := db.data[key]
	if !ok {
		count = 0
	}
	count++
	db.data[key] = count

	return count
}
