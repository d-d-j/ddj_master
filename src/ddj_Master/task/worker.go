package task

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/node"
	"ddj_Master/restApi"
	"fmt"
)

type TaskWorker struct {
	reqChan chan restApi.RestRequest
	pending int
	index   int
}

type Worker interface {
	Work(done chan Worker, idGen common.Int64Generator, balancer *node.LoadBalancer)
	RequestChan() chan restApi.RestRequest
	IncrementPending()
	DecrementPending()
	Id() int
	String() string
}

func NewTaskWorker(idx int, jobsPerWorker int32) Worker {
	w := new(TaskWorker)
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

func (w *TaskWorker) String() string {
	return fmt.Sprintf("Worker #%d pending:%d", w.index, w.pending)
}

func (w *TaskWorker) Id() int {
	return w.index
}

func (w *TaskWorker) IncrementPending() {
	w.pending++
}

func (w *TaskWorker) DecrementPending() {
	w.pending--
}

func (w *TaskWorker) RequestChan() chan restApi.RestRequest {
	return w.reqChan
}

func (w *TaskWorker) Work(done chan Worker, idGen common.Int64Generator, balancer *node.LoadBalancer) {
Loop:
	for {
		req := <-w.reqChan
		log.Debug("Worker is working")
		switch req.Type {
		case common.TASK_INSERT:
			log.Finest("Worker is processing [insert] task")
			insertNode := getNodeForInsert(req, balancer) // get nodeId from load balancer
			if insertNode == nil {
				done <- w
				continue Loop
			}
			id := idGen.GetId()
			t := NewTask(id, req, nil)
			log.Debug("Created new task with: id=", t.Id, " type=", t.Type, " size=", t.DataSize)
			TaskManager.AddChan <- t         // add task to dictionary
			message := createMessage(req, t) // create a message to send
			if message == nil {
				done <- w
				continue Loop
			}
			log.Finest("Sending message [%d] to node #%d", id, insertNode.Id)
			insertNode.Incoming <- message
			req.Response <- restApi.NewRestResponse("", id, nil)
			log.Finest("Worker finish task [%d]", id)
		case common.TASK_SELECT_ALL:
			nodes := node.NodeManager.GetNodes()
			avaliableNodes := len(nodes)
			log.Debug("Worker is processing [select all] task for all %d nodes", avaliableNodes)
			responseChan := make(chan *restApi.RestResponse, avaliableNodes)

			for _, n := range nodes {
				log.Finest("Sending [select all] task to #%d", n.Id)
				id := idGen.GetId()
				t := NewTask(id, req, responseChan)
				log.Fine("Created new task with: id=", t.Id, " type=", t.Type, " size=", t.DataSize)
				log.Finest(t)
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
			log.Debug("Waiting for data")
			for i := 0; i < avaliableNodes; i++ {
				result := <-responseChan
				log.Finest("Select: ", result)
			}
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
				log.Finest(t)
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
			if avaliableNodes != 0 {
				log.Debug("Waiting for status infos")
			}
			for i := 0; i < avaliableNodes; i++ {
				result := <-responseChan
				log.Finest("Get info", result)
			}
		default:
			log.Error("Worker can't handle task type ", req.Type)
		}
		log.Debug("Worker is done")
		done <- w
	}
}
