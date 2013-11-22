package task

import (
	"ddj_Master/restApi"
	"ddj_Master/common"
	log "code.google.com/p/log4go"
	"ddj_Master/node"
)

type Worker struct {
	reqChan  <-chan restApi.Request
	pending  int
	index    int
}

func (w *Worker) Work(done chan *Worker, idGen *Int64Generator, balancer *node.LoadBalancer) {
	for {
		req := <-w.reqChan
		switch req.Type {
		case common.TASK_INSERT:
			log.Finest("Worker is processing [insert] task")

			nodeId := balancer.CurrentInsertNodeId			// get nodeId from load balancer

			if nodeId == 0 {
				log.Warn("No node connected")
				req.Response <- nil
			}
															// get node
			var insertNode *Node
			nodeChan := make(chan<- *Node)
			nodeReq := node.GetNodeRequest{nodeId, nodeChan}
			node.NodeManager.GetChan <- nodeReq
			insertNode = <- nodeChan

															// Process insert task:
			id := idGen.getId()								// generate id
			t := NewTask(id, req)							// create new task for the request
			TaskManager.AddChan <- t    					// add task to dictionary
			var (
				message []byte
				err		error
			)
			message, err = t.MakeRequest().Encode()			// create a message to send
			insertNode.Incoming <- message					// send a message

			req.Response <- nil
		case common.TASK_SELECT_ALL:
			log.Finest("Worker is processing [select all] task")
			// TODO: Process select task
		default:
			log.Error("Worker can't handle task type ", req.Type)
		}
		done <- w
	}
}
