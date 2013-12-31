package task

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
	"ddj_Master/node"
	"fmt"
)

type TaskWorker struct {
	reqChan     chan dto.RestRequest
	getNodeChan chan node.GetNodeRequest
	done        chan Worker
	pending     int
	index       int
	idGenerator common.Int64Generator
	balancer    *node.LoadBalancer
}

type Worker interface {
	common.Int64Generator
	Work()
	RequestChan() chan dto.RestRequest
	Done()
	IncrementPending()
	DecrementPending()
	Id() int
	String() string
	getNodeForInsert() (*node.Node, error)
}

func NewTaskWorker(idx int, jobsPerWorker int32, getNodeChan chan node.GetNodeRequest, done chan Worker, nodeBalancer *node.LoadBalancer, idGen common.Int64Generator) Worker {
	w := new(TaskWorker)
	w.index = idx
	w.pending = 0
	w.reqChan = make(chan dto.RestRequest, jobsPerWorker)
	w.getNodeChan = getNodeChan
	w.balancer = nodeBalancer
	w.done = done
	w.idGenerator = idGen
	return w
}

func (w *TaskWorker) Work() {
	for {
		req := <-w.reqChan // GET REQUEST
		log.Finest(w, "Get request to process")
		j := w.getJob(req.Type)
		j(req)
		w.Done()
	}
}

func (w *TaskWorker) getNodeForInsert() (*node.Node, error) {
	nodeId := w.balancer.CurrentInsertNodeId
	if nodeId == common.CONST_UNINITIALIZED {
		return nil, fmt.Errorf("Balancer is uninitialized")
	}

	nodeChan := make(chan *node.Node)
	w.getNodeChan <- node.GetNodeRequest{NodeId: nodeId, BackChan: nodeChan}
	node := <-nodeChan
	if node == nil {
		return nil, fmt.Errorf("Node does not existing")
	}
	return node, nil
}

func createMessage(req dto.RestRequest, t *dto.Task, deviceId int32) []byte {
	var (
		message []byte
		err     error
	)
	message, err = t.MakeRequest(deviceId).Encode()
	if err != nil {
		log.Error("Error while encoding request - ", err)
		req.Response <- dto.NewRestResponse("Internal server error", 0, nil)
		return nil
	}
	return message
}

//Interface implementation
func (w *TaskWorker) String() string {
	return fmt.Sprintf("Worker #%d (pending:%d)", w.index, w.pending)
}

func (w *TaskWorker) Id() int { return w.index }

func (w *TaskWorker) IncrementPending() { w.pending++ }

func (w *TaskWorker) DecrementPending() {
	log.Finest(w, "Pending decrement")
	w.pending--
}

func (w *TaskWorker) RequestChan() chan dto.RestRequest { return w.reqChan }

func (w *TaskWorker) Done() {
	log.Finest(w, "Worker is done")
	w.done <- w
}

func (w *TaskWorker) GetId() int64 { return w.idGenerator.GetId() }
