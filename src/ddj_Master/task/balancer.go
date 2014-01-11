//this package is responsible for handling tasks.
package task

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
	"ddj_Master/node"
	"time"
)

//Balancer is responsible for dispaching task to different worker and taking care of sending node info request.
type Balancer struct {
	pool Pool
	done chan Worker
}

func NewBalancer(workersCount int32, jobForWorkerCount int32, loadBal *node.LoadBalancer) *Balancer {
	log.Info("Creating new task balancer with %d workers and %d pending jobs per each", workersCount, jobForWorkerCount)
	b := new(Balancer)
	done := make(chan Worker, workersCount+1)
	p := NewWorkersPool(workersCount, jobForWorkerCount, done, loadBal)
	b.pool = p
	b.done = done
	return b
}

func (b *Balancer) Balance(work <-chan dto.RestRequest, timeoutSeconds int32) {
	log.Info("Task manager balancer started")
	index := 0

	timeoutDuration := time.Duration(timeoutSeconds) * time.Second
	timeout := time.After(timeoutDuration)

	for {
		select {
		case w := <-b.done:
			b.completed(w)
		default:
			select {
			case req := <-work:
				b.dispatch(req, index)
				index = (index + 1) % b.pool.Len()
			case <-timeout:
				b.dispatch(dto.RestRequest{Type: common.TASK_INFO, Data: new(dto.EmptyElement), Response: nil}, index)
				index = (index + 1) % b.pool.Len()
				timeout = time.After(timeoutDuration)
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
