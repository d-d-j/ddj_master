package task

import (
	"ddj_Master/common"
	"ddj_Master/node"
	"ddj_Master/restApi"
)

type Pool []*Worker

func NewWorkersPool(size int32, jobsPerWorker int32, done chan *Worker) Pool {
	pool := make(Pool, size)
	idGen := common.NewTaskIdGenerator()
	loadBal := node.NewLoadBalancer()
	for index, worker := range pool {
		worker = &Worker{}
		worker.index = index
		worker.reqChan = make(chan restApi.RestRequest, jobsPerWorker)
		go worker.Work(done, idGen, loadBal)
	}
	return pool
}

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p Pool) Len() int { return len(p) }

func (p Pool) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = i
	p[j].index = j
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
