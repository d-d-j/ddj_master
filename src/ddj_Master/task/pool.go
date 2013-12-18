package task

import (
	"ddj_Master/common"
	"ddj_Master/node"
	"fmt"
)

type Pool []*Worker

func NewWorkersPool(size int32, jobsPerWorker int32, done chan *Worker, loadBal *node.LoadBalancer) Pool {
	pool := make(Pool, size)
	idGen := common.NewTaskIdGenerator()
	s := int(size)
	for i := 0; i < s; i++ {
		worker := NewWorker(i, jobsPerWorker)
		go worker.Work(done, idGen, loadBal)
		pool[i] = worker
	}
	return pool
}

func (p Pool) Less(i, j int) bool {
	if p[i] == nil {
		fmt.Println("ERROR ", i)
	}
	if p[j] == nil {
		fmt.Println("ERROR ", j)
	}
	return p[i].pending < p[j].pending
}

func (p Pool) Len() int { return len(p) }

func (p Pool) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = j
	p[j].index = i
}

func (p *Pool) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*p = old[0 : n-1]
	return item
}

func (p *Pool) Push(x interface{}) {
	n := len(*p)
	item := x.(*Worker)
	item.index = n
	*p = append(*p, item)
}
