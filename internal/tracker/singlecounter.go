package tracker

import "sync"

type SingleCounter struct {
	mu    sync.Mutex
	count int
}

func (sc *SingleCounter) Increment() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.count++
}

func (sc *SingleCounter) Get() int {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.count
}

func (sc *SingleCounter) Reset() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.count = 0
}
