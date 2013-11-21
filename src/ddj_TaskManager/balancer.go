package ddj_TaskManager

import (
	log "code.google.com/p/log4go"
	"container/heap"
	"ddj_RestApi"
)

/* TODO: 	Implement stop function for balancer
			it should wait for all workers to finish their jobs
			but it shouldn't delegate more tasks to workers
			after that balance function should return
 */

type Balancer struct {
	pool Pool
	done chan *Worker
}

func NewBalancer(workersCount int) *Balancer {
	b := new(Balancer)
	done := make(chan *Worker)
	p := NewWorkersPool(workersCount, done)
	b.pool = p
	b.done = done
	return b
}

func (b *Balancer) balance(work <-chan ddj_RestApi.Request) {
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

func (b *Balancer) dispach(req ddj_RestApi.Request) {
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
