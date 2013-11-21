package ddj_TaskManager

import "dto"

type Task struct {
	Id				int64
	Type			int32
	Data			dto.Dto
	DataSize		int32
	ResponseChan	chan<- dto.Dto
}

func NewTask() *Task {
	t := new(Task)
	return t
}
