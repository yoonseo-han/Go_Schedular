package runtime

type P struct {
	localQueue *RunQueue
}

func NewP() *P {
	return &P{
		localQueue: newRunQueue(),
	}
}
