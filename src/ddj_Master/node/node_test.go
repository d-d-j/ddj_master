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
	taskChan := make(chan dto.GetTaskRequest, 1)
	responseChan := make(chan *dto.RestResponse, 1)
	node := NewNode(NODE_ID, nil, taskChan)
	expected := &Info{NODE_ID, MemoryInfo{1, 1, 1, 1}}
	data, err := expected.Encode()
	if err != nil {
		t.Error("Error occurred", err)
	}
	result := *dto.NewResult(0, common.TASK_INFO, int32(expected.Size()), data)

	//Run tested method
	go node.processResult(result)

	//Simulate normal work
	taskRequest := <-taskChan
	if taskRequest.TaskId != TASK_ID {
		t.Error("Wrong task request. Expected: ", TASK_ID, " but got: ", taskRequest.TaskId)
	}

	taskRequest.BackChan <- responseChan

	response := <-responseChan
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

func Test_processResult_For_Select_Without_aggregation(t *testing.T) {
	t.Skip("Not implemented yet")
}
