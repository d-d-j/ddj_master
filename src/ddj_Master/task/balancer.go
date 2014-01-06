package task

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
	"ddj_Master/node"
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
	log.Info("Creating new task balancer with %d workers and %d pending jobs per each", workersCount, jobForWorkerCount)
	b := new(Balancer)
	done := make(chan Worker, workersCount + 1)
	p := NewWorkersPool(workersCount, jobForWorkerCount, done, loadBal)
	b.pool = p
	b.done = done
	return b
}

func (b *Balancer) Balance(work <-chan dto.RestRequest) {
	log.Info("Task manager balancer started")
	index := 0
	timeout := time.After(1*time.Second)

	for {
		select {
		case w := <-b.done:
			b.completed(w)
		default:
			select {
			case req := <-work:
				b.dispatch(req, index)
				index = (index + 1)%b.pool.Len()
			case <-timeout:
				b.dispatch(dto.RestRequest{Type: common.TASK_INFO, Data: new(dto.EmptyElement), Response: nil}, index)
				index = (index + 1)%b.pool.Len()
				timeout = time.After(1*time.Second)
			}
		}
	}
}

func (b *Balancer) dispatch(req dto.RestRequest, index int) {
	w := b.pool[index]
	log.Fine("Dispatch request to %s", w)
	w.RequestChan() <- req
	w.IncrementPending()
}

func (b *Balancer) completed(w Worker) {
	log.Fine("%s finish his job", w)
	w.DecrementPending()
}
