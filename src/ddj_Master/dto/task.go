package dto

import "fmt"

type Task struct {
	Id           int64
	Type         int32
	Data         Dto
	DataSize     int32
	ResponseChan chan *RestResponse
}

type GetTaskRequest struct {
	TaskId   int64
	BackChan chan *Task
}

func NewTask(id int64, request RestRequest, response chan *RestResponse) *Task {
	t := new(Task)
	t.Id = id
	t.Type = request.Type
	t.Data = request.Data
	t.DataSize = int32(request.Data.Size())
	t.ResponseChan = request.Response
	if response != nil {
		t.ResponseChan = response
	}
	return t
}

func (t *Task) MakeRequest() *Request {
	return NewRequest(t.Id, t.Type, t.DataSize, t.Data)
}

func (t *Task) String() string {
	return fmt.Sprintf("Task #%d, type: %d, size: %d", t.Id, t.Type, t.DataSize)
}
