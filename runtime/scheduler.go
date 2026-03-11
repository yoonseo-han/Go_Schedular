package runtime

import (
	"math/rand/v2"
	"sync"
)

const GOMAXPROCS = 10

type Scheduler struct {
	GOMAXPROCS     int
	addPIndex      int
	runPIndex      int
	pStore         []*P
	globalRunQueue *GlobalQueue
	idlePs         []*P   // P's not bound to any M
	idleMu         sync.Mutex
}

// Should be singleton?
func NewScheduler() *Scheduler {
	tempScheduler := &Scheduler{
		GOMAXPROCS: GOMAXPROCS,
		addPIndex:  0,
		runPIndex:  0,
		idlePs:     make([]*P, 0, GOMAXPROCS),
	}

	for range GOMAXPROCS {
		p := NewP(rand.Int64())
		tempScheduler.pStore = append(tempScheduler.pStore, p)
		tempScheduler.idlePs = append(tempScheduler.idlePs, p)
	}

	tempScheduler.globalRunQueue = newGlobalQueue()

	return tempScheduler
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

func (s *Scheduler) Add(g *G) {
	// Select which P to add the go routine to based on round robin
	g.state = RUNNABLE
	if !s.pStore[s.addPIndex].localQueue.add(g) {
		// Add to global run queue
		s.globalRunQueue.add(g)
	}
	s.addPIndex = (s.addPIndex + 1) % s.GOMAXPROCS
}

func (s *Scheduler) Run() {
	// Run go routine one time only
	g := s.pStore[s.runPIndex].localQueue.pop()
	if g == nil {
		return
	}

	if g.state != RUNNABLE {
		return
	}

	g.state = RUNNING
	g.funcToRun()
	g.state = WAITING

	s.pStore[s.runPIndex].localQueue.add(g)
	s.runPIndex = (s.runPIndex + 1) % s.GOMAXPROCS
}

func (s *Scheduler) Schedule() {
	for {
		// Get first runnable go routine
		g := s.pStore[s.runPIndex].localQueue.pop()

		if g == nil {
			// Move to next P if no runnable go routine found for current P
			s.runPIndex = (s.runPIndex + 1) % s.GOMAXPROCS
			continue
		}

		if g.state != RUNNABLE {
			continue
		}

		// Run the go routine
		g.state = RUNNING
		g.funcToRun()
		g.state = WAITING

		// Add the go routine back to the run queue
		s.pStore[s.runPIndex].localQueue.add(g)
		s.runPIndex = (s.runPIndex + 1) % s.GOMAXPROCS
		g.state = RUNNABLE
	}
}
