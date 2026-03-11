package runtime

import "sync"

// Global Queue to store all the go routines that are not bound to any P
type GlobalQueue struct {
	gStore []*G
	mu     sync.Mutex
}

func newGlobalQueue() *GlobalQueue {
	return &GlobalQueue{
		gStore: make([]*G, 0),
		mu:     sync.Mutex{}, // Should return pointer or not? Check with netty
	}
}

func (gq *GlobalQueue) add(g *G) {
	gq.mu.Lock()
	defer gq.mu.Unlock()
	gq.gStore = append(gq.gStore, g)
}

func (gq *GlobalQueue) pop() *G {
	gq.mu.Lock()
	defer gq.mu.Unlock()
	if len(gq.gStore) == 0 {
		return nil
	}
	firstElement := gq.gStore[0]
	gq.gStore = gq.gStore[1:]
	return firstElement
}
