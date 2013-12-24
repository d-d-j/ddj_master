package dto

import (
	"ddj_Master/common"
	"fmt"
)

type Task struct {
	Id           int64
	Type         int32
	Data         Dto
	DataSize     int32
	ResponseChan chan *RestResponse // channel for sending response to (REST API) client
	ResultChan   chan *RestResponse // channel for sending result to worker
}

type GetTaskRequest struct {
	TaskId   int64
	BackChan chan *Task
}

func NewTask(id int64, request RestRequest, resultChan chan *RestResponse) *Task {
	t := new(Task)
	t.Id = id
	t.Type = request.Type
	t.Data = request.Data
	t.DataSize = int32(request.Data.Size())
	t.ResponseChan = request.Response
	t.ResultChan = resultChan
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
