package runtime

import (
	"fmt"
)

type Scheduler struct {
	runQueue *RunQueue
}

// Should be singleton?
func newScheduler() *Scheduler {
	return &Scheduler{
		runQueue: newRunQueue(),
	}
}

func (s *Scheduler) add(g *G) {
	s.runQueue.add(g)
}

func (s *Scheduler) run() {
	g := s.runQueue.pop()
	if g == nil {
		return
	}

	if g.state != WAITING {
		fmt.Println("G is not waiting")
		return
	}

	g.state = RUNNING
	g.funcToRun()
	g.state = WAITING
	s.runQueue.add(g)
}
