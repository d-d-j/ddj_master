package restApi

import (
	"ddj_Master/dto"
	"fmt"
)

type RestResponse struct {
	Error		string
	TaskId		int64
	Data		[]dto.Dto
}

func NewRestResponse(err string, taskId int64, data []dto.Dto) *RestResponse {
	rr := new(RestResponse)
	rr.Error = err
	rr.TaskId = taskId
	rr.Data = data
	return rr
}

func (r *RestResponse) String() string {
	return fmt.Sprintf("Rest response with TaskId = %d with Error = %s and data:\n%s", r.TaskId, r.Error, r.Data)
}
