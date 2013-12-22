package task

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
	"ddj_Master/node"
	"ddj_Master/restApi"
	"time"
)

/* TODO: 	Implement stop function for balancer
it should wait for all workers to finish their jobs
but it shouldn't delegate more tasks to workers
after that balance function should return
*/

type Balancer struct {
	pool Pool
	done chan Worker
}

func NewBalancer(workersCount int32, jobForWorkerCount int32, loadBal *node.LoadBalancer) *Balancer {
	b := new(Balancer)
	done := make(chan Worker)
	p := NewWorkersPool(workersCount, jobForWorkerCount, done, loadBal)
	b.pool = p
	b.done = done
	return b
}

func (b *Balancer) Balance(work <-chan restApi.RestRequest) {
	log.Info("Task manager balancer started")
	index := 0
	for {
		select {
		case req := <-work:
			b.dispatch(req, index)
			index = (index + 1) % b.pool.Len()
		case w := <-b.done:
			b.completed(w)
		case <-time.After(5 * time.Second):
			b.dispatch(restApi.RestRequest{common.TASK_INFO, new(dto.EmptyElement), nil}, 0)
		}
	}
}

func (b *Balancer) dispatch(req restApi.RestRequest, index int) {
	w := b.pool[index]
	log.Fine("Dispatch request to %s", w)
	w.RequestChan() <- req
	w.IncrementPending()
}

func (b *Balancer) completed(w Worker) {
	w.DecrementPending()
}
