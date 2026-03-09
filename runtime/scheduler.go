package runtime

import (
	"math/rand/v2"
)

const GOMAXPROCS = 10

type Scheduler struct {
	GOMAXPROCS int
	addPIndex  int
	runPIndex  int
	pStore     []*P
}

// Should be singleton?
func NewScheduler() *Scheduler {
	tempScheduler := &Scheduler{
		GOMAXPROCS: GOMAXPROCS,
		addPIndex:  0,
		runPIndex:  0,
	}

	for range GOMAXPROCS {
		tempScheduler.pStore = append(tempScheduler.pStore, NewP(rand.Int64()))
	}

	return tempScheduler
}

func (s *Scheduler) Add(g *G) {
	// Select which P to add the go routine to based on round robin
	g.state = RUNNABLE
	s.pStore[s.addPIndex].localQueue.add(g)
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
			break
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
