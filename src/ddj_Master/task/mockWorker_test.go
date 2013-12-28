package task

import (
	"ddj_Master/common"
	"ddj_Master/node"
)

func MockWorkersPool(size int, jobsPerWorker int32, done chan Worker, loadBal *node.LoadBalancer) Pool {
	pool := make(Pool, size)
	idGen := common.NewTaskIdGenerator()
	s := int(size)
	for i := 0; i < s; i++ {
		worker := NewTaskWorker(i, jobsPerWorker, nil)
		go worker.Work(done, idGen, loadBal)
		pool[i] = worker
	}
	return pool
}
