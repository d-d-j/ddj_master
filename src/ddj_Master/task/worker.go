package task

import (
	"ddj_Master/restApi"
	"ddj_Master/common"
	log "code.google.com/p/log4go"
	"ddj_Master/node"
)

type Worker struct {
	reqChan  chan restApi.RestRequest
	pending  int
	index    int
}

func (w *Worker) Work(done chan *Worker, idGen common.Int64Generator, balancer *node.LoadBalancer) {
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
			var insertNode *node.Node
			nodeChan := make(chan *node.Node)
			nodeReq := node.GetNodeRequest{nodeId, nodeChan}
			node.NodeManager.GetChan <- nodeReq
			insertNode = <- nodeChan

															// Process insert task:
			id := idGen.GetId()								// generate id
			t := NewTask(id, req)							// create new task for the request
			TaskManager.AddChan <- t    					// add task to dictionary
			var (
				message []byte
				err		error
			)
			message, err = t.MakeRequest().Encode()			// create a message to send
			if err != nil {
				log.Error("Error while encoding request - ", err)
			}
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
