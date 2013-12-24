package task

import (
	"ddj_Master/common"
	"ddj_Master/node"
)

type Pool []Worker

func NewWorkersPool(size int32, jobsPerWorker int32, done chan Worker, loadBal *node.LoadBalancer) Pool {
	pool := make(Pool, size)
	idGen := common.NewTaskIdGenerator()
	s := int(size)
	for i := 0; i < s; i++ {
		worker := NewTaskWorker(i, jobsPerWorker, node.NodeManager.GetChan)
		go worker.Work(done, idGen, loadBal)
		pool[i] = worker
	}
	return pool
}

func (p Pool) Len() int { return len(p) }
