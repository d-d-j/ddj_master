package task

import (
	"ddj_Master/restApi"
	"ddj_Master/common"
	log "code.google.com/p/log4go"
)

type Worker struct {
	reqChan  <-chan restApi.Request
	pending  int
	index    int
}

func (w *Worker) Work(done chan *Worker, idGen *Int64Generator) {
	for {
		req := <-w.reqChan
		switch req.Type {
		case common.TASK_INSERT:
			log.Finest("Worker is processing [insert] task")
			// Process insert task
			id := idGen.getId()			// generate id
			t := NewTask(id, req)		// create new task for the request
			TaskManager.AddChan <- t    // add task to dictionary
			                            // create message to send
			var (
				message []byte
				err		error
			)
			message, err = t.MakeRequest().Encode()
										// get node from load balancer
										// send a message

		case common.TASK_SELECT_ALL:
			log.Finest("Worker is processing [select all] task")
			// TODO: Process select task
		default:
			log.Error("Worker can't handle task type ", req.Type)
		}
		done <- w
	}
}
