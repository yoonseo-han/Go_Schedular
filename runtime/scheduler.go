package runtime

import (
	"math/rand/v2"
	"sync"
)

const GOMAXPROCS = 10

type Scheduler struct {
	GOMAXPROCS     int
	globalRunQueue *GlobalQueue
	idlePs         []*P // P's not bound to any M
	idleMu         sync.Mutex
	wg             sync.WaitGroup // tracks number of G's not yet completed
}

// NewScheduler creates the global scheduler with GOMAXPROCS P's in the idle pool.
func NewScheduler() *Scheduler {
	s := &Scheduler{
		GOMAXPROCS:     GOMAXPROCS,
		globalRunQueue: newGlobalQueue(),
		idlePs:         make([]*P, 0, GOMAXPROCS),
	}
	for range GOMAXPROCS {
		s.idlePs = append(s.idlePs, NewP(rand.Int64()))
	}
	return s
}

// AcquireP returns a P from the idle pool. Caller (M) should bind it.
// Returns nil if no P is available.
func (s *Scheduler) AcquireP() *P {
	s.idleMu.Lock()
	defer s.idleMu.Unlock()
	if len(s.idlePs) == 0 {
		return nil
	}
	last := len(s.idlePs) - 1
	p := s.idlePs[last]
	s.idlePs = s.idlePs[:last]
	return p
}

// ReleaseP puts P back into the idle pool (e.g. when M blocks or exits).
func (s *Scheduler) ReleaseP(p *P) {
	if p == nil {
		return
	}
	s.idleMu.Lock()
	defer s.idleMu.Unlock()
	s.idlePs = append(s.idlePs, p)
}

// Add puts a runnable G into the global run queue. M's will pick it up when
// they need work (from their P's local queue or from global).
func (s *Scheduler) Add(g *G) {
	s.wg.Add(1)
	g.state = RUNNABLE
	s.globalRunQueue.add(g)
}

// GCompleted is called by an M when it finishes executing a G (once per G).
func (s *Scheduler) GCompleted() {
	s.wg.Done()
}

// Wait blocks until all G's added via Add have completed.
func (s *Scheduler) Wait() {
	s.wg.Wait()
}
