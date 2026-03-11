package runtime

type M struct {
	id          int64
	designatedP *P
}

func NewM(id int64) *M {
	return &M{
		id: id,
	}
}

// Run starts the M's schedule loop. If M has no bound P, it acquires one from
// the scheduler's idle P pool.
func (m *M) Run(sched *Scheduler) {
	if m.designatedP == nil {
		m.designatedP = sched.AcquireP()
	}
	// TODO: if m.designatedP still nil (no idle P), block or spin
	// else: run schedule loop using m.designatedP
}
