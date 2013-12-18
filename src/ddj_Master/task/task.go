package task

import (
	"ddj_Master/dto"
	"ddj_Master/restApi"
)

type Task struct {
	Id           int64
	Type         int32
	Data         dto.Dto
	DataSize     int32
	ResponseChan chan *restApi.RestResponse
}

func NewTask(id int64, request restApi.RestRequest, respone chan *restApi.RestResponse) *Task {
	t := new(Task)
	t.Id = id
	t.Type = request.Type
	t.Data = request.Data
	t.DataSize = int32(request.Data.Size())
	t.ResponseChan = request.Response
	return t
}

func (t *Task) MakeRequest() *dto.Request {
	return dto.NewRequest(t.Id, t.Type, t.DataSize, t.Data)
}
