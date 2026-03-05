package runtime

import "sync"

type RunQueue struct {
	mu     sync.Mutex
	gStore []*G // Pointer to always ensure unique G
}

func newRunQueue() *RunQueue {
	return &RunQueue{
		gStore: make([]*G, 0),
	}
}

func (rq *RunQueue) add(g *G) {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	rq.gStore = append(rq.gStore, g)
}

func (rq *RunQueue) get() *G {
	rq.mu.Lock()
	defer rq.mu.Unlock()
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

func (rq *RunQueue) poll() *G {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	if len(rq.gStore) == 0 {
		return nil
	}

	// Normally G is considered as unique hence return by reference
	return rq.gStore[0]
}
