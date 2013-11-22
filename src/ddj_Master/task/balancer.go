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

func NewBalancer(workersCount int) *Balancer {
	b := new(Balancer)
	done := make(chan *Worker)
	tasks := make(map[int64]Task)
	p := NewWorkersPool(workersCount, done)
	b.pool = p
	b.done = done
	b.tasks = tasks
	return b
}

func (b *Balancer) balance(work <-chan restApi.Request) {
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

func (b *Balancer) dispatch(req restApi.Request) {
	w := heap.Pop(&b.pool).(*Worker)
	log.Fine("Dispach request to ", w)
	w.ReqChan <- req
	w.Pending++
	heap.Push(&b.pool, w)
}

func (b *Balancer) completed(w *Worker) {
	w.Pending--
	heap.Remove(&b.pool, w.Index)
	heap.Push(&b.pool, w)
}
