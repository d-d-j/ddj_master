package ddj_TaskManager

import "ddj_Dto"

type Task struct {
	Id				int64
	Type			int32
	Data			ddj_Dto.Dto
	DataSize		int32
	ResponseChan	chan<- ddj_Dto.Dto
}

func NewTask() *Task {
	t := new(Task)
	return t
}
