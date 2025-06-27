package tracker

import (
	"sync"
)

var (
	mu      sync.RWMutex
	counter = make(map[string]*PathStats)
)

// Increment increases the count for the given path
func Request(path string, method string) {
	mu.Lock()
	defer mu.Unlock()
	ps, ok := counter[path]
	if !ok {
		ps = newPathStats()
		counter[path] = ps
	}
	ps.Methods.Increment(method)
	ps.Total.Increment()

}

func Response(path string, status int) {
	mu.Lock()
	defer mu.Unlock()
	ps, ok := counter[path]
	if !ok {
		return
	}
	ps.Status.Increment(status)

}

// GetCount returns the current count for a given path
// func GetCount(path string) int {
// 	mu.RLock()
// 	defer mu.RUnlock()
// 	return counter[path]
// }

// GetAll returns a copy of the whole map
func GetAll() map[string]interface{} {
	mu.RLock()
	defer mu.RUnlock()

	// Copy to prevent external mutation
	copy := make(map[string]interface{}, len(counter))
	for k, v := range counter {
		copy[k] = v.Snapshot()
	}
	return copy
}
