package runtime

type M struct {
	id          int64
	designatedP *P
	scheduler   *Scheduler
}

// Inject schedular into M
func NewM(id int64, sched *Scheduler) *M {
	return &M{
		id:        id,
		scheduler: sched,
	}
}

// Run starts the M's schedule loop. If M has no bound P, it acquires one from
// the scheduler's idle P pool.
func (m *M) Run(sched *Scheduler) {
	if m.designatedP == nil {
		m.designatedP = sched.AcquireP()
	}
	if m.designatedP != nil {
		m.scheduleLoop()
	}
}

func (m *M) scheduleLoop() {
	// Assume M has a designated P
	for {
		// Get next G from P's local run queue
		g := m.designatedP.localQueue.pop()
		if g == nil {
			g = m.scheduler.globalRunQueue.pop()
		}
		if g == nil {
			// No more Gs to run, release P and exit
			break
		}

		g.state = RUNNING
		g.funcToRun()
		g.state = DEAD
		m.scheduler.GCompleted()
	}
	m.scheduler.ReleaseP(m.designatedP)
	m.designatedP = nil
}
