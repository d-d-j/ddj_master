package task

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/node"
	"ddj_Master/restApi"
	"fmt"
)

type Worker struct {
	reqChan chan restApi.RestRequest
	pending int
	index   int
}

func NewWorker(idx int, jobsPerWorker int32) *Worker {
	w := new(Worker)
	w.index = idx
	w.pending = 0
	w.reqChan = make(chan restApi.RestRequest, jobsPerWorker)
	return w
}

func getNodeForInsert(req restApi.RestRequest, balancer *node.LoadBalancer) *node.Node {
	nodeId := balancer.CurrentInsertNodeId
	if nodeId == common.CONST_UNINITIALIZED {
		log.Warn("No node connected")
		req.Response <- restApi.NewRestResponse("No node connected", 0, nil)
		return nil
	}
	// get node
	var insertNode *node.Node
	nodeChan := make(chan *node.Node)
	nodeReq := node.GetNodeRequest{nodeId, nodeChan}
	node.NodeManager.GetChan <- nodeReq
	insertNode = <-nodeChan
	return insertNode
}

func createMessage(req restApi.RestRequest, t *Task) []byte {
	var (
		message []byte
		err     error
	)
	message, err = t.MakeRequest().Encode()
	if err != nil {
		log.Error("Error while encoding request - ", err)
		req.Response <- restApi.NewRestResponse("Internal server error", 0, nil)
		return nil
	}
	return message
}

func (w *Worker) String() string {
	return fmt.Sprintf("Worker #%d pending:%d", w.index, w.pending)
}

func (w *Worker) Work(done chan *Worker, idGen common.Int64Generator, balancer *node.LoadBalancer) {
Loop:
	for {
		req := <-w.reqChan
		log.Debug("Worker is working")
		switch req.Type {
		case common.TASK_INSERT:
			log.Finest("Worker is processing [insert] task")
			insertNode := getNodeForInsert(req, balancer) // get nodeId from load balancer
			if insertNode == nil {
				continue Loop
			}
			id := idGen.GetId()
			t := NewTask(id, req, nil)
			log.Debug("Created new task with: id=", t.Id, " type=", t.Type, " size=", t.DataSize)
			TaskManager.AddChan <- t         // add task to dictionary
			message := createMessage(req, t) // create a message to send
			if message == nil {
				continue Loop
			}
			log.Finest("Sending message [%d] to node #%d", id, insertNode.Id)
			insertNode.Incoming <- message
			req.Response <- restApi.NewRestResponse("", id, nil)
			log.Finest("Worker finish task [%d]", id)
		case common.TASK_SELECT_ALL:
			log.Finest("Worker is processing [select all] task")
			// TODO: Process select task
		case common.TASK_INFO:
			nodes := node.NodeManager.GetNodes()
			avaliableNodes := len(nodes)
			log.Debug("Worker is processing [info] task for all %d nodes", avaliableNodes)
			responseChan := make(chan *restApi.RestResponse, avaliableNodes)
			for _, n := range nodes {
				log.Finest("Sending [info] task to #%d", n.Id)
				id := idGen.GetId()
				t := NewTask(id, req, responseChan)
				log.Fine("Created new task with: id=", t.Id, " type=", t.Type, " size=", t.DataSize)
				TaskManager.AddChan <- t         // add task to dictionary
				message := createMessage(req, t) // create a message to send
				message, err := t.MakeRequest().Encode()
				if err != nil {
					log.Error("Cannot parse request: %s", err)
					continue
				}
				log.Finest("Sending message [%d] to node #%d", id, n.Id)
				nodeChan := make(chan *node.Node)
				nodeReq := node.GetNodeRequest{n.Id, nodeChan}
				node.NodeManager.GetChan <- nodeReq
				currentNode := <-nodeChan
				currentNode.Incoming <- message
				log.Finest("Worker send task [%d]", id)
			}
			log.Debug("Waiting for status infos")
			for i := 0; i < avaliableNodes; i++ {
				result := <-responseChan
				log.Finest("Get info", result)
			}
		default:
			log.Error("Worker can't handle task type ", req.Type)
		}
		done <- w
	}
}
