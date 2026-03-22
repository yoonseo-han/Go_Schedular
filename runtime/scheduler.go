package runtime

import (
	"sync"
)

const GOMAXPROCS = 10

type Scheduler struct {
	GOMAXPROCS     int
	globalRunQueue *GlobalQueue

	idlePs        []*P // P's not bound to any M
	listOfPs      []*P
	currentPIndex int64 // index of the current P

	idleMu sync.Mutex
	addMu  sync.Mutex
	wg     sync.WaitGroup // tracks number of G's not yet completed
}

// NewScheduler creates the global scheduler with GOMAXPROCS P's in the idle pool.
func NewScheduler() *Scheduler {
	s := &Scheduler{
		GOMAXPROCS:     GOMAXPROCS,
		globalRunQueue: newGlobalQueue(),
		idlePs:         make([]*P, 0, GOMAXPROCS),
		currentPIndex:  0,
		listOfPs:       make([]*P, 0, GOMAXPROCS),
	}
	for i := range GOMAXPROCS {
		p := NewP(int64(i))
		s.listOfPs = append(s.listOfPs, p)
		s.idlePs = append(s.idlePs, p)
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
	// fmt.Printf("Releasing P id=%d ptr=%p\n", p.id, p)
	s.idlePs = append(s.idlePs, p)
}

// Add puts a runnable G into the current P's local run queue.
func (s *Scheduler) Add(g *G) {
	s.addMu.Lock()
	defer s.addMu.Unlock()

	s.wg.Add(1)
	g.state = RUNNABLE

	currentP := s.listOfPs[s.currentPIndex]
	placedOnLocal := currentP.add(g)
	if !placedOnLocal {
		// fmt.Println("Adding G", g.id, "to global run queue")
		s.globalRunQueue.add(g)
	} else {
		// fmt.Println("Added G", g.id, "to local run queue of P", currentP.id)
	}
	s.currentPIndex = (s.currentPIndex + 1) % int64(len(s.listOfPs))
}

// GCompleted is called by an M when it finishes executing a G (once per G).
func (s *Scheduler) GCompleted() {
	s.wg.Done()
}

// Wait blocks until all G's added via Add have completed.
func (s *Scheduler) Wait() {
	s.wg.Wait()
}
