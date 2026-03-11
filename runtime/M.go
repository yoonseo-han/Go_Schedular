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

func (m *M) Run() {
	// Check if there is designated P
	if m.designatedP == nil {
		// Get P from idle

	}
}
