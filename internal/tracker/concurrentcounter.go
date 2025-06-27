package tracker

import (
	"encoding/json"
	"sync"
)

type ConcurrentCounter[T comparable] struct {
	mu     sync.Mutex
	counts map[T]int
}

func NewConcurrentCounter[T comparable]() *ConcurrentCounter[T] {
	return &ConcurrentCounter[T]{counts: make(map[T]int)}
}

func (c *ConcurrentCounter[T]) Increment(key T) {

	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[key]++
}

func (c *ConcurrentCounter[T]) Snapspot() map[T]int {

	c.mu.Lock()
	defer c.mu.Unlock()
	snapspot := make(map[T]int)
	for k, v := range c.counts {
		snapspot[k] = v
	}
	return snapspot
}

func (c *ConcurrentCounter[T]) ToJson() ([]byte, error) {

	c.mu.Lock()
	defer c.mu.Unlock()
	return json.Marshal(c.counts)
}
