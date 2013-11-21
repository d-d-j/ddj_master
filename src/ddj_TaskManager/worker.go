package ddj_TaskManager

import "ddj_RestApi"

type Worker struct {
	reqChan  <-chan ddj_RestApi.Request
	pending  int
	index    int
}

func (w *Worker) work(done chan *Worker) {
	for {
		req := <-w.requests
		req.c <- req.fn()
		done <- w
	}
}
