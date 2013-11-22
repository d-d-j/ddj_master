package taskManager

import "ddj_Master/restApi"

type Worker struct {
	reqChan  <-chan restApi.Request
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
