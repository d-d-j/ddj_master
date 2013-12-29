package task

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
	"sort"
)

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

	responses := handleRequestForAllNodes(w.GetId(), req)
	if responses == nil {
		return false
	}

	// TODO: REDUCE RESPONSES
	responseToClient := make([]dto.Dto, 0, len(responses))
	for i := 0; i < len(responses); i++ {
		responseToClient = append(responseToClient, responses[i].Data...)
	}

	sort.Sort(dto.ByTime(responseToClient))

	// PASS REDUCED RESPONSES TO CLIENT
	req.Response <- dto.NewRestResponse("", 0, responseToClient)
	return true
}

func (w *TaskWorker) Info(req dto.RestRequest) bool {
	log.Debug("Worker is processing [info] task")

	responses := handleRequestForAllNodes(w.GetId(), req)
	if responses == nil {
		return false
	}

	// TODO: SET NODE INFO IN NODES
	for i := 0; i < len(responses); i++ {
		log.Finest("Get info %v", responses)
	}

	return true
}
