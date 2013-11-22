package task

import "ddj_Master/restApi"

type Worker struct {
	ReqChan  <-chan restApi.Request
	Pending  int
	Index    int
}

func (w *Worker) Work(done chan *Worker) {
	for {
		req := <-w.requests
		req.c <- req.fn()
		done <- w
	}
}
