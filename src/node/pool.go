package node

import (
	log "code.google.com/p/log4go"
	"container/heap"
)

type Pool []*Worker

type Worker struct {
	requests chan Request
	pending  int
	index    int
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
