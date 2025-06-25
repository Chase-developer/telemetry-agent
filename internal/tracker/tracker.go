package tracker

import (
	"sync"
)

var (
	mu      sync.RWMutex
	counter = make(map[string]int)
)

// Increment increases the count for the given path
func Increment(path string) {
	mu.Lock()
	defer mu.Unlock()
	counter[path]++
}

// GetCount returns the current count for a given path
func GetCount(path string) int {
	mu.RLock()
	defer mu.RUnlock()
	return counter[path]
}

// GetAll returns a copy of the whole map
func GetAll() map[string]int {
	mu.RLock()
	defer mu.RUnlock()

	// Copy to prevent external mutation
	copy := make(map[string]int, len(counter))
	for k, v := range counter {
		copy[k] = v
	}
	return copy
}
