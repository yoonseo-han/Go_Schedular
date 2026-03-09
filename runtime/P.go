package runtime

type P struct {
	id         int64
	localQueue *RunQueue
}

func NewP(id int64) *P {
	return &P{
		id:         id,
		localQueue: newRunQueue(),
	}
}
