package task

import "ddj_Master/restApi"

type Pool []*Worker

func NewWorkersPool(size int, done chan *Worker) Pool {
	pool := make(Pool, size)
	for index, worker := range pool {
		worker.Index = index
		worker.ReqChan = make(chan restApi.Request)
		go worker.Work(done)
	}
	return pool
}

func (p Pool) Less(i, j int) bool {
	return p[i].Pending < p[j].Pending
}

func (p Pool) Len() int { return len(p) }

func (p Pool) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].Index = i
	p[j].Index = j
}

func (p *Pool) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	item.Index = -1 // for safety
	*p = old[0 : n-1]
	return item
}

func (p *Pool) Push(x interface{}) {
	n := len(*p)
	item := x.(*Worker)
	item.Index = n
	*p = append(*p, item)
}
