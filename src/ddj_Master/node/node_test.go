package node

import (
	"ddj_Master/common"
	"ddj_Master/dto"
	"testing"
)

func Test_processResult_For_Info(t *testing.T) {
	//Prepare
	const (
		NODE_ID int32 = 0
		TASK_ID int64 = 0
	)
	// CREATE CHANNEL FOR GETTING TASKS USED BY NODE
	getTaskChan := make(chan dto.GetTaskRequest)
	// CREATE CHANNEL FOR SENDING RESULT (TO WORKER)
	resultChan := make(chan *dto.RestResponse)
	// CREATE CHANNEL FOR SENDING RESPONSE TO CLIENT (REST API)
	responseChan := make(chan *dto.RestResponse)

	task := dto.NewTask(TASK_ID, dto.RestRequest{common.TASK_INFO, &dto.EmptyElement{}, responseChan}, resultChan)
	node := NewNode(NODE_ID, nil, getTaskChan)

	// PREPARE DATA FOR TEST
	expected := &Info{NODE_ID, MemoryInfo{1, 1, 1, 1}}
	data, err := expected.Encode()
	if err != nil {
		t.Error("Error occurred", err)
	}
	result := *dto.NewResult(0, common.TASK_INFO, int32(expected.Size()), data)

	// RUN TESTED METHOD
	go node.processResult(result)

	// SIMULATE WORK
	getTaskRequest := <-getTaskChan
	if getTaskRequest.TaskId != TASK_ID {
		t.Error("Wrong task request. Expected: ", TASK_ID, " but got: ", getTaskRequest.TaskId)
	}

	getTaskRequest.BackChan <- task

	response := <-resultChan
	if response.TaskId != TASK_ID {
		t.Error("Wrong task Id in response. Expected: ", TASK_ID, " but got: ", response.TaskId)
	}
	if response.Error != "" {
		t.Error("Error occurred", response.Error)
	}
	if len(response.Data) != 1 {
		t.Error("Wrong data returned. Expected only one value")
	}

	actual := response.Data[0]
	if expected.String() != actual.String() {
		t.Error("Expected: ", expected, " but got: ", actual)
	}

}

func Test_processResult_For_Select_Without_Aggregation_Empty_Response(t *testing.T) {

	//Prepare
	const (
		NODE_ID int32 = 0
		TASK_ID int64 = 0
	)
	// CREATE CHANNEL FOR GETTING TASKS USED BY NODE
	getTaskChan := make(chan dto.GetTaskRequest)
	// CREATE CHANNEL FOR SENDING RESULT (TO WORKER)
	resultChan := make(chan *dto.RestResponse)
	// CREATE CHANNEL FOR SENDING RESPONSE TO CLIENT (REST API)
	responseChan := make(chan *dto.RestResponse)

	task := dto.NewTask(TASK_ID, dto.RestRequest{common.TASK_SELECT, &dto.EmptyElement{}, responseChan}, resultChan)
	node := NewNode(NODE_ID, nil, getTaskChan)

	// PREPARE DATA FOR TEST
	data := []byte{}
	result := *dto.NewResult(0, common.TASK_SELECT, 0, data)

	// RUN TESTED METHOD
	go node.processResult(result)

	// SIMULATE WORK
	getTaskRequest := <-getTaskChan
	if getTaskRequest.TaskId != TASK_ID {
		t.Error("Wrong task request. Expected: ", TASK_ID, " but got: ", getTaskRequest.TaskId)
	}

	getTaskRequest.BackChan <- task

	// ASSERTIONS
	response := <-resultChan
	if response.TaskId != TASK_ID {
		t.Error("Wrong task Id in response. Expected: ", TASK_ID, " but got: ", response.TaskId)
	}
	if response.Error != "" {
		t.Error("Error occurred", response.Error)
	}
	if len(response.Data) != 0 {
		t.Error("Wrong data returned. Expected no value")
	}
}

func Test_processResult_For_Select_Without_Aggregation_One_Element_In_Response(t *testing.T) {

	//Prepare
	const (
		NODE_ID int32 = 0
		TASK_ID int64 = 0
	)
	// CREATE CHANNEL FOR GETTING TASKS USED BY NODE
	getTaskChan := make(chan dto.GetTaskRequest)
	// CREATE CHANNEL FOR SENDING RESULT (TO WORKER)
	resultChan := make(chan *dto.RestResponse)
	// CREATE CHANNEL FOR SENDING RESPONSE TO CLIENT (REST API)
	responseChan := make(chan *dto.RestResponse)

	task := dto.NewTask(TASK_ID, dto.RestRequest{common.TASK_SELECT, &dto.EmptyElement{}, responseChan}, resultChan)
	node := NewNode(NODE_ID, nil, getTaskChan)

	// PREPARE DATA FOR TEST
	expected := dto.NewElement(1, 2, 3, 0.33)
	data, err := expected.Encode()
	if err != nil {
		t.Error("Error occurred", err)
	}
	result := *dto.NewResult(0, common.TASK_SELECT, int32(expected.Size()), data)

	// RUN TESTED METHOD
	go node.processResult(result)

	// SIMULATE WORK
	getTaskRequest := <-getTaskChan
	if getTaskRequest.TaskId != TASK_ID {
		t.Error("Wrong task request. Expected: ", TASK_ID, " but got: ", getTaskRequest.TaskId)
	}

	getTaskRequest.BackChan <- task

	// ASSERTIONS
	response := <-resultChan
	if response.TaskId != TASK_ID {
		t.Error("Wrong task Id in response. Expected: ", TASK_ID, " but got: ", response.TaskId)
	}
	if response.Error != "" {
		t.Error("Error occurred", response.Error)
	}
	if len(response.Data) != 1 {
		t.Error("Wrong data returned. Expected only one value")
	}

	actual := response.Data[0]
	if expected.String() != actual.String() {
		t.Error("Expected: ", expected, " but got: ", actual)
	}
}

func Test_processResult_For_Select_Without_Aggregation_3_Elements_In_Response(t *testing.T) {

	//Prepare
	const (
		NODE_ID int32 = 0
		TASK_ID int64 = 0
	)
	// CREATE CHANNEL FOR GETTING TASKS USED BY NODE
	getTaskChan := make(chan dto.GetTaskRequest)
	// CREATE CHANNEL FOR SENDING RESULT (TO WORKER)
	resultChan := make(chan *dto.RestResponse)
	// CREATE CHANNEL FOR SENDING RESPONSE TO CLIENT (REST API)
	responseChan := make(chan *dto.RestResponse)

	task := dto.NewTask(TASK_ID, dto.RestRequest{common.TASK_SELECT, &dto.EmptyElement{}, responseChan}, resultChan)
	node := NewNode(NODE_ID, nil, getTaskChan)

	// PREPARE DATA FOR TEST
	expected := dto.Dtos{dto.NewElement(1, 2, 3, 0.33), dto.NewElement(4, 5, 6, 0.66), dto.NewElement(7, 8, 9, 0.99)}
	data, err := expected.Encode()
	if err != nil {
		t.Error("Error occurred", err)
	}
	result := *dto.NewResult(0, common.TASK_SELECT, int32(expected.Size()), data)

	// RUN TESTED METHOD
	go node.processResult(result)

	// SIMULATE WORK
	getTaskRequest := <-getTaskChan
	if getTaskRequest.TaskId != TASK_ID {
		t.Error("Wrong task request. Expected: ", TASK_ID, " but got: ", getTaskRequest.TaskId)
	}

	getTaskRequest.BackChan <- task

	// ASSERTIONS
	response := <-resultChan
	if response.TaskId != TASK_ID {
		t.Error("Wrong task Id in response. Expected: ", TASK_ID, " but got: ", response.TaskId)
	}
	if response.Error != "" {
		t.Error("Error occurred", response.Error)
	}
	if len(response.Data) != 3 {
		t.Error("Wrong data returned. Expected ", len(expected), " values")
	}

	for index, actual := range response.Data {
		if expected[index].String() != actual.String() {
			t.Error("Expected: ", expected, " but got: ", actual)
		}
	}
}
