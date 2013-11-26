package common

//Result
const (
	RESULT_OK  int32 = 0 //everything went right
	RESULT_ERR int32 = 1 // there were some problem
)

//TaskType
const (
	ERROR			int32 = 0
	TASK_INSERT     int32 = 1
	TASK_SELECT_ALL int32 = 2
	FLUSH			int32 = 3
)