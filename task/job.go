package task

import (
	log "code.google.com/p/log4go"
	"github.com/d-d-j/ddj_master/common"
	"github.com/d-d-j/ddj_master/dto"
	"github.com/d-d-j/ddj_master/node"
	"github.com/d-d-j/ddj_master/reduce"
	"time"
)

type job func(dto.RestRequest) bool

func (w *TaskWorker) getJob(taskType int32) job {
	switch taskType {
	case common.TASK_INSERT:
		return w.insert
	case common.TASK_SELECT:
		return w.selectTask
	case common.TASK_INFO:
		return w.info
	case common.TASK_FLUSH:
		return w.flush
	}
	log.Error("Worker can't handle task type ", taskType)
	return func(req dto.RestRequest) bool { return false }
}

func (w *TaskWorker) insert(req dto.RestRequest) bool {
	log.Finest(w, " is processing [insert] task")

	// GET NODE FOR INSERT
	insertNode, err := w.getNodeForInsert()
	if err != nil {
		log.Warn(w, " has problem with getting node to insert, ", err)
		req.Response <- dto.NewRestResponse("No node connected", common.TASK_UNINITIALIZED, nil)
		return false
	}

	// CREATE TASK
	id := w.GetId()
	t := dto.NewTask(id, req, nil)
	log.Fine(w, " created new %s", t)
	TaskManager.AddChan <- t // add task to dictionary

	// CREATE MESSAGE
	message, err := t.MakeRequest(insertNode.PreferredDeviceId).Encode()
	if err != nil {
		log.Error(w, " encourage error while encoding request - ", err)
		req.Response <- dto.NewRestResponse("Internal server error", t.Id, nil)
		return false
	}

	// SEND MESSAGE
	log.Finest(w, " is sending message to node #%d", id, insertNode.Id)
	insertNode.Incoming <- message
	TaskManager.DelChan <- t.Id

	// PASS RESPONSE TO CLIENT
	req.Response <- dto.NewRestResponse("", id, nil)

	return true
}

func (w *TaskWorker) selectTask(req dto.RestRequest) bool {
	log.Debug(w, " is processing [select] task")

	availableNodes := node.NodeManager.GetNodesLen()
	t, responseChan := CreateTaskForRequest(req, availableNodes, w.GetId())
	if availableNodes == 0 {
		log.Error("No nodes connected")
		req.Response <- dto.NewRestResponse("No nodes connected", t.Id, nil)
		return false
	}
	if !BroadcastTaskToAllNodes(t) {
		req.Response <- dto.NewRestResponse("Internal server error", t.Id, nil)
		return false
	}

	responses := parseResults(GatherAllResponses(availableNodes, responseChan), req.Data.(*dto.Query))
	TaskManager.DelChan <- t.Id
	log.Fine("Got %d responses", len(responses))
	aggregate := reduce.GetAggregator(t.AggregationType)
	responseToClient := aggregate(responses)

	req.Response <- dto.NewRestResponse("", t.Id, responseToClient)
	return true
}

func (w *TaskWorker) info(req dto.RestRequest) bool {
	log.Debug(w, " is processing [info] task")
	// TODO: Handle errors better

	availableNodes := node.NodeManager.GetNodesLen()
	t, responseChan := CreateTaskForRequest(req, availableNodes, w.GetId())
	if availableNodes == 0 {
		log.Error("No nodes connected")
		return false
	}
	if !BroadcastTaskToAllNodes(t) {
		return false
	}

	responses := parseResultsToInfos(GatherAllResponses(availableNodes, responseChan))
	if responses == nil {
		return false
	}
	TaskManager.DelChan <- t.Id

	node.NodeManager.InfoChan <- responses

	// TODO: SET NODE INFO IN NODES
	for i := 0; i < len(responses); i++ {
		log.Finest(w, "Get info %v", responses[i])
	}

	return true
}

func parseResultsToInfos(results []*dto.Result) []*dto.Info {
	infoSize := (&dto.MemoryInfo{}).Size()
	resultsCount := len(results)
	elements := make([]*dto.Info, 0, resultsCount)

	for i := 0; i < resultsCount; i++ {
		length := len(results[i].Data) / infoSize
		for j := 0; j < length; j++ {
			var e dto.Info
			err := e.MemoryInfo.Decode(results[i].Data[j*infoSize:])
			if err != nil {
				log.Error("Problem with parsing data", err)
				continue
			}
			e.NodeId = results[i].NodeId
			elements = append(elements, &e)
		}
	}
	return elements
}

func (w *TaskWorker) flush(req dto.RestRequest) bool {
	log.Debug(w, " is processing flush task")

	availableNodes := node.NodeManager.GetNodesLen()
	t, responseChan := CreateTaskForRequest(req, availableNodes, w.GetId())
	if availableNodes == 0 {
		log.Error("No nodes connected")
		req.Response <- dto.NewRestResponse("No nodes connected", t.Id, nil)
		return false
	}
	if !BroadcastTaskToAllNodes(t) {
		req.Response <- dto.NewRestResponse("Internal server error", t.Id, nil)
		return false
	}

	GatherAllResponses(availableNodes, responseChan)

	TaskManager.DelChan <- t.Id

	log.Finest("Flush is done, sending response to client")
	req.Response <- dto.NewRestResponse("", t.Id, nil)

	return true
}

func CreateTaskForRequest(req dto.RestRequest, numResponses int, taskId int64) (*dto.Task, chan *dto.Result) {
	responseChan := make(chan *dto.Result, numResponses)

	// CREATE TASK
	t := dto.NewTask(taskId, req, responseChan)
	log.Fine("Created new %s", t)
	TaskManager.AddChan <- t // add task to dictionary

	return t, responseChan
}

func BroadcastTaskToAllNodes(t *dto.Task) bool {
	// CREATE MESSAGE
	message, err := t.MakeRequestForAllGpus().Encode()
	if err != nil {
		log.Error("Error while encoding request - ", err)
		return false
	}

	// SEND MESSAGE TO ALL NODES
	node.NodeManager.SendToAllNodes(message)
	return true
}

func GatherAllResponses(numResponses int, responseChan chan *dto.Result) []*dto.Result {
	responses := make([]*dto.Result, 0, numResponses)

	timeoutDuration := time.Duration(5000) * time.Second
	timeout := time.After(timeoutDuration)
	// WAIT FOR ALL RESPONSES
	for i := 0; i < numResponses; i++ {
		select {
		case response := <-responseChan:
			responses = append(responses, response)
		case <-timeout:
			log.Error("Timeout! Got %d/%d responses", i, numResponses)
			break
		}
	}

	return responses
}
