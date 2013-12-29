package task

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
	"ddj_Master/node"
	"fmt"
	"sort"
)

type TaskWorker struct {
	reqChan     chan dto.RestRequest
	getNodeChan chan node.GetNodeRequest
	done        chan Worker
	pending     int
	index       int
	balancer    *node.LoadBalancer
	taskManager *Manager
}

type Worker interface {
	Work(idGen common.Int64Generator)
	RequestChan() chan dto.RestRequest
	Done()
	IncrementPending()
	DecrementPending()
	Id() int
	String() string
	getNodeForInsert() (*node.Node, error)
}

func NewTaskWorker(idx int, jobsPerWorker int32, getNodeChan chan node.GetNodeRequest, done chan Worker, nodeBalancer *node.LoadBalancer) Worker {
	w := new(TaskWorker)
	w.index = idx
	w.pending = 0
	w.reqChan = make(chan dto.RestRequest, jobsPerWorker)
	w.getNodeChan = getNodeChan
	w.balancer = nodeBalancer
	w.done = done
	return w
}

func (w *TaskWorker) Work(idGen common.Int64Generator) {
Loop:
	for {
		req := <-w.reqChan // GET REQUEST

		switch req.Type {
		case common.TASK_INSERT:
			log.Finest("Worker is processing [insert] task")

			// GET NODE FOR INSERT
			insertNode, err := w.getNodeForInsert()
			if err != nil {
				log.Warn("Problem with getting node to insert, ", err)
				req.Response <- dto.NewRestResponse("No node connected", common.TASK_UNINITIALIZED, nil)
				w.Done()
				continue Loop
			}

			// CREATE TASK
			id := idGen.GetId()
			t := dto.NewTask(id, req, nil)
			log.Fine("Created new %s", t)
			TaskManager.AddChan <- t // add task to dictionary

			// CREATE MESSAGE
			message, err := t.MakeRequest(insertNode.PreferredDeviceId).Encode()
			if err != nil {
				log.Error("Error while encoding request - ", err)
				req.Response <- dto.NewRestResponse("Internal server error", 0, nil)
				w.Done()
				continue Loop
			}

			// SEND MESSAGE
			log.Finest("Sending message to node #%d", id, insertNode.Id)
			insertNode.Incoming <- message

			// PASS RESPONSE TO CLIENT
			req.Response <- dto.NewRestResponse("", id, nil)

			// TODO: Change this to set task status or sth, then wait for response about insert from node
			// then set status again to success
			//TaskManager.DelChan <- t.Id

		case common.TASK_SELECT:
			log.Debug("Worker is processing [select] task")

			responses := handleRequestForAllNodes(idGen, req)
			if responses == nil {
				w.Done()
				continue Loop
			}

			// TODO: REDUCE RESPONSES
			responseToClient := make([]dto.Dto, 0, len(responses))
			for i := 0; i < len(responses); i++ {
				responseToClient = append(responseToClient, responses[i].Data...)
			}

			sort.Sort(dto.ByTime(responseToClient))

			// PASS REDUCED RESPONSES TO CLIENT
			req.Response <- dto.NewRestResponse("", 0, responseToClient)

		case common.TASK_INFO:
			log.Debug("Worker is processing [info] task")

			responses := handleRequestForAllNodes(idGen, req)
			if responses == nil {
				w.Done()
				continue Loop
			}

			// TODO: SET NODE INFO IN NODES
			for i := 0; i < len(responses); i++ {
				log.Finest("Get info %v", responses)
			}

		default:
			log.Error("Worker can't handle task type ", req.Type)
		}
		log.Debug("Worker is done")
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
	return <-nodeChan, nil
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

func handleRequestForAllNodes(idGen common.Int64Generator, req dto.RestRequest) []*dto.RestResponse {
	// TODO: Handle errors better than return nil

	// GET NODES
	nodes := node.NodeManager.GetNodes()
	avaliableNodes := len(nodes)
	responseChan := make(chan *dto.RestResponse, avaliableNodes)
	if avaliableNodes == 0 {
		return nil
	}

	// CREATE TASK
	id := idGen.GetId()
	t := dto.NewTask(id, req, responseChan)
	log.Fine("Created new %s", t)
	TaskManager.AddChan <- t // add task to dictionary

	// CREATE MESSAGE
	message, err := t.MakeRequestForAllGpus().Encode()
	if err != nil {
		log.Error("Error while encoding request - ", err)
		req.Response <- dto.NewRestResponse("Internal server error", 0, nil)
		return nil
	}

	// SEND MESSAGE TO ALL NODES
	node.NodeManager.SendToAllNodes(message)

	responses := make([]*dto.RestResponse, avaliableNodes)

	// WAIT FOR ALL RESPONSES
	for i := 0; i < avaliableNodes; i++ {
		responses[i] = <-responseChan
		log.Finest("Got task result [%d/%d] - %s", i, avaliableNodes, responses[i])
	}

	// REMOVE TASK
	//TaskManager.DelChan <- t.Id

	return responses
}

//Interface implementation
func (w *TaskWorker) String() string {
	return fmt.Sprintf("Worker #%d pending:%d", w.index, w.pending)
}

func (w *TaskWorker) Id() int { return w.index }

func (w *TaskWorker) IncrementPending() { w.pending++ }

func (w *TaskWorker) DecrementPending() { w.pending-- }

func (w *TaskWorker) RequestChan() chan dto.RestRequest { return w.reqChan }

func (w *TaskWorker) Done() { w.done <- w }
