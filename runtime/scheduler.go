package runtime

import (
	"fmt"
	"time"
)

type Scheduler struct {
	runQueue *RunQueue
}

// Should be singleton?
func NewScheduler() *Scheduler {
	return &Scheduler{
		runQueue: newRunQueue(),
	}
}

func (s *Scheduler) Add(g *G) {
	// Change state to runnable
	g.state = RUNNABLE
	s.runQueue.add(g)
}

func (s *Scheduler) Run() {
	// Run go routine one time only
	g := s.runQueue.pop()
	if g == nil {
		return
	}

	if g.state != RUNNABLE {
		fmt.Println("G is not waiting")
		return
	}

	g.state = RUNNING
	g.funcToRun()
	g.state = DEAD
}

func (s *Scheduler) Schedule() {
	for {
		// Get first runnable go routine
		g := s.runQueue.pop()

		if g == nil {
			time.Sleep(1 * time.Millisecond)
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
		s.runQueue.add(g)
		g.state = RUNNABLE
	}
}
