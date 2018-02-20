package hunkee

// worker need to safely use any amount of parsers for
// corresponded format without any memory overhead and races.
type worker struct {
	id     uint
	busy   bool
	parent *mapper
	pos    int
}

func (m *mapper) initWorkers(amount int) {
	m.workers = make([]*worker, amount)
	for i := 0; i < amount; i++ {
		m.workers[i] = &worker{parent: m}
	}
}

func (m *mapper) aquireWorker() *worker {
	for i := 0; i < len(m.workers); i++ {
		if !m.workers[i].busy {
			m.workers[i].busy = true
			return m.workers[i]
		}
	}
	// TODO implement worker queue
	return nil
}

func (w *worker) seek(i int) *field {
	if i >= len(w.parent.tokensSeq) {
		return nil
	}

	f := w.parent.fields[w.parent.tokensSeq[i]]
	w.pos = i + 1
	return f
}

func (w *worker) next() *field {
	return w.seek(w.pos)
}

func (w *worker) first() *field {
	return w.seek(0)
}

func (w *worker) release() {
	w.busy = false
}
