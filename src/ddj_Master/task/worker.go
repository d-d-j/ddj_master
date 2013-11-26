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

func getNodeForInsert(req restApi.RestRequest, balancer *node.LoadBalancer) *node.Node {
	nodeId := balancer.CurrentInsertNodeId
	if nodeId == 0 {
		log.Warn("No node connected")
		req.Response <- restApi.NewRestResponse("No node connected", 0, nil)
		return nil
	}
	// get node
	var insertNode *node.Node
	nodeChan := make(chan *node.Node)
	nodeReq := node.GetNodeRequest{nodeId, nodeChan}
	node.NodeManager.GetChan <- nodeReq
	insertNode = <- nodeChan
	return insertNode
}

func createMessage(req restApi.RestRequest, t *Task) []byte {
	var (
		message []byte
		err		error
	)
	message, err = t.MakeRequest().Encode()
	if err != nil {
		log.Error("Error while encoding request - ", err)
		req.Response <- restApi.NewRestResponse("Internal server error", 0, nil)
		return nil
	}
	return message
}

func (w *Worker) Work(done chan *Worker, idGen common.Int64Generator, balancer *node.LoadBalancer) {
	for {
		req := <-w.reqChan
		switch req.Type {
		case common.TASK_INSERT:
			log.Finest("Worker is processing [insert] task")
			insertNode := getNodeForInsert(req, balancer)			// get nodeId from load balancer
			if insertNode == nil {
				break
			}
			id := idGen.GetId()										// generate id
			t := NewTask(id, req)									// create new task for the request
			TaskManager.AddChan <- t    							// add task to dictionary
			message := createMessage(req, t)						// create a message to send
			if message == nil {
				break
			}
			insertNode.Incoming <- message							// send a message to node
			req.Response <- restApi.NewRestResponse("", id, nil)	// send response to client
		case common.TASK_SELECT_ALL:
			log.Finest("Worker is processing [select all] task")
			// TODO: Process select task
		default:
			log.Error("Worker can't handle task type ", req.Type)
		}
		done <- w
	}
}