package task

import (
	"ddj_Master/dto"
	"fmt"
)

type Task struct {
	Id           int64
	Type         int32
	Data         dto.Dto
	DataSize     int32
	ResponseChan chan *dto.RestResponse
}

func NewTask(id int64, request dto.RestRequest, response chan *dto.RestResponse) *Task {
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

func (t *Task) MakeRequest() *dto.Request {
	return dto.NewRequest(t.Id, t.Type, t.DataSize, t.Data)
}

func (t *Task) String() string {
	return fmt.Sprintf("Task #%d, type: %d, size: %d", t.Id, t.Type, t.DataSize)
}
