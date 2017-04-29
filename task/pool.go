package task

import (
	"github.com/d-d-j/ddj_master/common"
	"github.com/d-d-j/ddj_master/node"
)

//This is type to handle pool of workers. It is used to implement thread pool pattern since each worker is independent gorutine
type Pool []Worker

//Constructor of worker Pool. Count of workers is specify by size
func NewWorkersPool(size int32, jobsPerWorker int32, done chan Worker, loadBal *node.LoadBalancer) Pool {
	pool := make(Pool, size)
	idGen := common.NewTaskIdGenerator()
	s := int(size)
	for i := 0; i < s; i++ {
		worker := NewTaskWorker(i, jobsPerWorker, node.NodeManager.GetChan, done, loadBal, idGen)
		go worker.Work()
		pool[i] = worker
	}
	return pool
}

//Return size of Pool
func (p Pool) Len() int { return len(p) }
