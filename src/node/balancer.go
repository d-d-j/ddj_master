package node

import (
	log "code.google.com/p/log4go"
	"container/heap"
	"dto"
)

type Request struct {
	fn func() dto.Result
	c  chan dto.Result
}

type Balancer struct {
	pool Pool
	done chan *Worker
}

func (b *Balancer) balance(work chan Request) {
	log.Info("Balancer started")
	for {
		select {
		case req := <-work:
			b.dispach(req)
		case w := <-b.done:
			b.completed(w)
		}
	}
}

func (b *Balancer) dispach(req Request) {
	w := heap.Pop(&b.pool).(*Worker)
	log.Fine("Dispach request to ", w)
	w.requests <- req
	w.pending++
	heap.Push(&b.pool, w)
}

func (b *Balancer) completed(w *Worker) {
	w.pending--
	heap.Remove(&b.pool, w.index)
	heap.Push(&b.pool, w)
}
