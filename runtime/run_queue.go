package runtime

import "sync"

type RunQueue struct {
	mu      sync.Mutex
	gStore  []*G // Pointer to always ensure unique G
	maxSize int
}

func newRunQueue() *RunQueue {
	return &RunQueue{
		gStore:  make([]*G, 0),
		maxSize: 10,
	}
}

func (rq *RunQueue) add(g *G) bool {
	rq.mu.Lock()
	defer rq.mu.Unlock()

	if len(rq.gStore) >= rq.maxSize {
		return false
	}

	rq.gStore = append(rq.gStore, g)
	return true
}

func (rq *RunQueue) get() *G {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	if len(rq.gStore) == 0 {
		return nil
	}
	return rq.gStore[0]
}

func (rq *RunQueue) peek() *G {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	if len(rq.gStore) == 0 {
		return nil
	}
	return rq.gStore[0]
}

func (rq *RunQueue) pop() *G {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	if len(rq.gStore) == 0 {
		return nil
	}
	firstElement := rq.gStore[0]
	rq.gStore = rq.gStore[1:]
	return firstElement
}
