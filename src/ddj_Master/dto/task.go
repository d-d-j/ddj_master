package dto

import (
	"ddj_Master/common"
	"fmt"
)

type Task struct {
	Id              int64
	Type            int32
	AggregationType int32
	Data            Dto
	DataSize        int32
	ResponseChan    chan *RestResponse // channel for sending response to (REST API) client
	ResultChan      chan *Result       // channel for sending result to worker
}

type GetTaskRequest struct {
	TaskId   int64
	BackChan chan *Task
}

func NewTask(id int64, request RestRequest, resultChan chan *Result) *Task {
	t := new(Task)
	t.Id = id
	t.Type = request.Type
	t.AggregationType = common.CONST_UNINITIALIZED
	t.Data = request.Data
	t.DataSize = int32(request.Data.Size())
	t.ResponseChan = request.Response
	t.ResultChan = resultChan

	if t.Type == common.TASK_SELECT {
		if query, ok := t.Data.(*Query); ok {
			t.AggregationType = query.AggregationType
		} else {
			panic("Type mismatch. TaskType select can be used only with Query data")
		}

	}

	return t
}

func (t *Task) MakeRequest(deviceId int32) *Request {
	return NewRequest(t.Id, t.Type, t.DataSize, t.Data, deviceId)
}

func (t *Task) MakeRequestForAllGpus() *Request {
	return NewRequest(t.Id, t.Type, t.DataSize, t.Data, common.ALL_GPUs)
}

func (t *Task) String() string {
	return fmt.Sprintf("Task #%d, type: %d, size: %d", t.Id, t.Type, t.DataSize)
}
