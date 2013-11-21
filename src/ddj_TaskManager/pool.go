package ddj_TaskManager

import "ddj_RestApi"

type Pool []*Worker

func NewWorkersPool(size int, done chan *Worker) Pool {
	pool := make(Pool, size)
	for index, worker := range pool {
		worker.index = index
		worker.reqChan = make(chan ddj_RestApi.Request)
		go worker.work(done)
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
