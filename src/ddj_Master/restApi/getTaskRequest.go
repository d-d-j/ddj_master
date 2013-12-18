package restApi

type GetTaskRequest struct {
	TaskId   int64
	BackChan chan chan *RestResponse
}
