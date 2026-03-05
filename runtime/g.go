package runtime

import (
	"math/rand/v2"
	"time"
)

type State string

const (
	RUNNING  State = "running"
	WAITING  State = "waiting"
	DEAD     State = "dead"
	NEW      State = "new"
	RUNNABLE State = "runnable"
)

type G struct {
	id          int64
	state       State
	createdAt   time.Time
	scheduledAt time.Time
	priority    int
	funcToRun   func()
}

func newG() *G {
	return &G{
		id:    rand.Int64(),
		state: NEW,
	}
}
