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
		if m.designatedP != nil {
			// ptr shows identity; same id may repeat in the log when P is released and re-acquired (not concurrent sharing).
			// fmt.Printf("M %d acquired P id=%d ptr=%p\n", m.id, m.designatedP.id, m.designatedP)
		}
	}
	if m.designatedP != nil {
		m.scheduleLoop()
	}
}

func (m *M) scheduleLoop() {
	for {
		g := m.designatedP.localQueue.pop()
		if g == nil {
			g = m.scheduler.globalRunQueue.pop()
		}
		if g == nil {
			g = m.stealWork()
		}
		if g == nil {
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

// stealWork pops one G from another P's local queue (toy work-stealing).
func (m *M) stealWork() *G {
	for _, p := range m.scheduler.listOfPs {
		if p == m.designatedP {
			continue
		}
		if g := p.localQueue.pop(); g != nil {
			return g
		}
	}
	return nil
}
