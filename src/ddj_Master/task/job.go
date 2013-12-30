package task

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
	"ddj_Master/node"
	"ddj_Master/reduce"
)

type job func(dto.RestRequest) bool

func (w *TaskWorker) getJob(taskType int32) job {
	switch taskType {
	case common.TASK_INSERT:
		return w.Insert
	case common.TASK_SELECT:
		return w.Select
	case common.TASK_INFO:
		return w.Info
	case common.TASK_FLUSH:
		return w.Info
	}
	log.Error("Worker can't handle task type ", taskType)
	return func(req dto.RestRequest) bool { return false }
}

func (w *TaskWorker) Insert(req dto.RestRequest) bool {
	log.Finest("Worker is processing [insert] task")

	// GET NODE FOR INSERT
	insertNode, err := w.getNodeForInsert()
	if err != nil {
		log.Warn("Problem with getting node to insert, ", err)
		req.Response <- dto.NewRestResponse("No node connected", common.TASK_UNINITIALIZED, nil)
		return false
	}

	// CREATE TASK
	id := w.GetId()
	t := dto.NewTask(id, req, nil)
	log.Fine("Created new %s", t)
	TaskManager.AddChan <- t // add task to dictionary

	// CREATE MESSAGE
	message, err := t.MakeRequest(insertNode.PreferredDeviceId).Encode()
	if err != nil {
		log.Error("Error while encoding request - ", err)
		req.Response <- dto.NewRestResponse("Internal server error", 0, nil)
		return false
	}

	// SEND MESSAGE
	log.Finest("Sending message to node #%d", id, insertNode.Id)
	insertNode.Incoming <- message

	// PASS RESPONSE TO CLIENT
	req.Response <- dto.NewRestResponse("", id, nil)

	// TODO: Change this to set task status or sth, then wait for response about insert from node
	// then set status again to success
	//TaskManager.DelChan <- t.Id

	return true
}

func (w *TaskWorker) Select(req dto.RestRequest) bool {
	log.Debug("Worker is processing [select] task")

	availableNodes := node.NodeManager.GetNodesLen()
	t, responseChan := CreateTaskForRequest(req, availableNodes, w.GetId())
	if availableNodes == 0 {
		log.Error("No nodes connected")
		req.Response <- dto.NewRestResponse("No nodes connected", 0, nil)
		return false
	}
	if BroadcastTaskToAllNodes(t, req) == -1 {
		return false
	}

	responses := parseResultsToElements(GatherAllResponses(availableNodes, responseChan))
	log.Finest("Response from node ", responses)
	aggregate := reduce.GetAggregator(t.AggregationType)
	responseToClient := aggregate(responses)
	log.Finest("Response after aggregation ", responseToClient)
	// PASS REDUCED RESPONSES TO CLIENT
	req.Response <- dto.NewRestResponse("", 0, responseToClient)
	return true
}

func parseResultsToElements(results []*dto.Result) []*dto.Element {
	elementSize := (&dto.Element{}).Size()
	resultsCount := len(results)
	elements := make([]*dto.Element, 0, resultsCount)

	for i := 0; i < resultsCount; i++ {
		length := len(results[i].Data) / elementSize
		for j := 0; j < length; j++ {
			var e dto.Element
			err := e.Decode(results[i].Data[j*elementSize:])
			if err != nil {
				log.Error("Problem with parsing data", err)
				continue
			}
			elements = append(elements, &e)
		}
	}
	return elements
}

func (w *TaskWorker) Info(req dto.RestRequest) bool {
	log.Debug("Worker is processing [info] task")
	// TODO: Handle errors better

	availableNodes := node.NodeManager.GetNodesLen()
	t, responseChan := CreateTaskForRequest(req, availableNodes, w.GetId())
	if availableNodes == 0 {
		log.Error("No nodes connected")
		req.Response <- dto.NewRestResponse("No nodes connected", 0, nil)
		return false
	}
	if BroadcastTaskToAllNodes(t, req) == -1 {
		return false
	}

	responses := parseResultsToInfos(GatherAllResponses(availableNodes, responseChan))
	if responses == nil {
		return false
	}

	// TODO: SET NODE INFO IN NODES
	for i := 0; i < len(responses); i++ {
		log.Finest("Get info %v", responses)
	}

	return true
}

//TODO: Generate Info not memoryInfo.
func parseResultsToInfos(results []*dto.Result) []*dto.MemoryInfo {
	infoSize := (&dto.MemoryInfo{}).Size()
	resultsCount := len(results)
	elements := make([]*dto.MemoryInfo, 0, resultsCount)

	for i := 0; i < resultsCount; i++ {
		length := len(results[i].Data) / infoSize
		for j := 0; j < length; j++ {
			var e dto.MemoryInfo
			err := e.Decode(results[i].Data[j*infoSize:])
			if err != nil {
				log.Error("Problem with parsing data", err)
				continue
			}
			elements = append(elements, &e)
		}
	}
	return elements
}

func (w *TaskWorker) Flush(req dto.RestRequest) bool {
	log.Debug("Worker is processing flush task")

	availableNodes := node.NodeManager.GetNodesLen()
	t, responseChan := CreateTaskForRequest(req, availableNodes, w.GetId())
	if availableNodes == 0 {
		log.Error("No nodes connected")
		req.Response <- dto.NewRestResponse("No nodes connected", 0, nil)
		return false
	}
	if BroadcastTaskToAllNodes(t, req) == -1 {
		return false
	}

	responses := GatherAllResponses(availableNodes, responseChan)
	if responses == nil {
		return false
	}

	for i := 0; i < len(responses); i++ {
		log.Finest("Get flush response %v", responses)
	}

	req.Response <- dto.NewRestResponse("", responses[0].TaskId, []dto.Dto{})

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

func BroadcastTaskToAllNodes(t *dto.Task, req dto.RestRequest) int {
	// CREATE MESSAGE
	message, err := t.MakeRequestForAllGpus().Encode()
	if err != nil {
		log.Error("Error while encoding request - ", err)
		req.Response <- dto.NewRestResponse("Internal server error", 0, nil)
		return -1
	}

	// SEND MESSAGE TO ALL NODES
	node.NodeManager.SendToAllNodes(message)
	return 0
}

func GatherAllResponses(numResponses int, responseChan chan *dto.Result) []*dto.Result {
	responses := make([]*dto.Result, numResponses)

	// WAIT FOR ALL RESPONSES
	for i := 0; i < numResponses; i++ {
		responses[i] = <-responseChan
		log.Finest("Got task result [%d/%d] - %s", i, numResponses, responses[i])
	}

	// REMOVE TASK
	//TaskManager.DelChan <- t.Id

	return responses
}
