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
	reqChan chan dto.RestRequest
	pending int
	index   int
}

type Worker interface {
	Work(done chan Worker, idGen common.Int64Generator, balancer *node.LoadBalancer)
	RequestChan() chan dto.RestRequest
	IncrementPending()
	DecrementPending()
	Id() int
	String() string
}

func NewTaskWorker(idx int, jobsPerWorker int32) Worker {
	w := new(TaskWorker)
	w.index = idx
	w.pending = 0
	w.reqChan = make(chan dto.RestRequest, jobsPerWorker)
	return w
}

func (w *TaskWorker) Work(done chan Worker, idGen common.Int64Generator, balancer *node.LoadBalancer) {
Loop:
	for {
		req := <-w.reqChan // GET REQUEST

		switch req.Type {
		case common.TASK_INSERT:
			log.Finest("Worker is processing [insert] task")

			// GET NODE FOR INSERT
			insertNode := getNodeForInsert(req, balancer) // get nodeId from load balancer
			if insertNode == nil {
				done <- w
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
				done <- w
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

			responses := handleRequestForAllNodes(done, idGen, balancer, req)
			if responses == nil {
				done <- w
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

			responses := handleRequestForAllNodes(done, idGen, balancer, req)
			if responses == nil {
				done <- w
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
		done <- w
	}
}

func getNodeForInsert(req dto.RestRequest, balancer *node.LoadBalancer) *node.Node {
	nodeId := balancer.CurrentInsertNodeId
	if nodeId == common.CONST_UNINITIALIZED {
		log.Warn("No node connected")
		req.Response <- dto.NewRestResponse("No node connected", common.TASK_UNINITIALIZED, nil)
		return nil
	}
	// get node
	var insertNode *node.Node
	nodeChan := make(chan *node.Node)
	nodeReq := node.GetNodeRequest{NodeId: nodeId, BackChan: nodeChan}
	node.NodeManager.GetChan <- nodeReq
	insertNode = <-nodeChan
	return insertNode
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

func handleRequestForAllNodes(done chan Worker, idGen common.Int64Generator, balancer *node.LoadBalancer, req dto.RestRequest) []*dto.RestResponse {
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
	message, err := t.MakeRequestForAllGpu().Encode()
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
	TaskManager.DelChan <- t.Id

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
