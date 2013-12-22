package task

import (
	"ddj_Master/common"
	"ddj_Master/node"
	"ddj_Master/restApi"
	"fmt"
)

type MockWorker TaskWorker

func (w *MockWorker) String() string {
	return fmt.Sprintf("Worker #%d pending:%d", w.index, w.pending)
}

func (w *MockWorker) IncrementPending() {
	w.pending++
}

func (w *MockWorker) DecrementPending() {
	w.pending--
}

func (w *MockWorker) RequestChan() chan restApi.RestRequest {
	return w.reqChan
}

func (w *MockWorker) Work(done chan Worker, idGen common.Int64Generator, balancer *node.LoadBalancer) {
	for {
		<-w.reqChan
		done <- w
	}
}

func (w *MockWorker) Id() int {
	return w.index
}

func MockWorkersPool(size int, jobsPerWorker int32, done chan Worker, loadBal *node.LoadBalancer) Pool {
	pool := make(Pool, size)
	idGen := common.NewTaskIdGenerator()
	s := int(size)
	for i := 0; i < s; i++ {
		worker := NewTaskWorker(i, jobsPerWorker)
		go worker.Work(done, idGen, loadBal)
		pool[i] = worker
	}
	return pool
}

func NewMockWorker(idx int, jobsPerWorker int32) Worker {
	w := new(MockWorker)
	w.index = idx
	w.pending = 0
	w.reqChan = make(chan restApi.RestRequest, jobsPerWorker)
	return w
}
