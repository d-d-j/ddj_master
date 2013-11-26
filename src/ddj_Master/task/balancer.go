package task

import (
	log "code.google.com/p/log4go"
	"container/heap"
	"ddj_Master/restApi"
)

/* TODO: 	Implement stop function for balancer
			it should wait for all workers to finish their jobs
			but it shouldn't delegate more tasks to workers
			after that balance function should return
 */

type Balancer struct {
	pool 		Pool
	done 		chan *Worker
}

func NewBalancer(workersCount int32, jobForWorkerCount int32) *Balancer {
	b := new(Balancer)
	done := make(chan *Worker)
	p := NewWorkersPool(workersCount, jobForWorkerCount, done)
	b.pool = p
	b.done = done
	heap.Init(&p)
	return b
}

func (b *Balancer) Balance(work <-chan restApi.RestRequest) {
	log.Info("Task manager balancer started")
	for {
		select {
		case req := <-work:
			b.dispatch(req)
		case w := <-b.done:
			b.completed(w)
		}
	}
}

func (b *Balancer) dispatch(req restApi.RestRequest) {
	w := heap.Pop(&(b.pool)).(*Worker)
	log.Fine("Dispach request to ", w)
	w.reqChan <- req
	w.pending++
	heap.Push(&(b.pool), w)
}

func (b *Balancer) completed(w *Worker) {
	w.pending--
	heap.Remove(&(b.pool), w.index)
	heap.Push(&(b.pool), w)
}
