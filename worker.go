package hunkee

import (
	"sync/atomic"
)

type pool struct {
	size int
	free chan *worker
}

func initPool(size int) *pool {
	free := make(chan *worker, size)
	for i := 0; i < size; i++ {
		free <- &worker{}
	}
	return &pool{
		size: size,
		free: free,
	}
}

func (p *pool) get(parent *mapper) *worker {
	w := <-p.free
	w.parent = parent
	return w
}

func (p *pool) put(w *worker) {
	p.free <- w
}

// worker need to safely use any amount of parsers for
// corresponded format without any memory overhead and races.
type worker struct {
	id  uint
	pos uint32
	len uint32

	parent *mapper
}

func (w *worker) seek(i uint32) *field {
	if w.len == 0 {
		w.len = uint32(len(w.parent.tokensSeq))
	}

	if i >= w.len {
		return nil
	}

	f := w.parent.getField(w.parent.tokensSeq[i])
	atomic.StoreUint32(&w.pos, uint32(i+1))
	return f
}

func (w *worker) next() *field {
	return w.seek(atomic.LoadUint32(&w.pos))
}

func (w *worker) first() *field {
	return w.seek(0)
}

func (w *worker) free() {
	w.parent.workerPool.put(w)
}
