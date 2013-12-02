package common

//Result
const (
	RESULT_OK  int32 = 0
	RESULT_ERR int32 = 1
)

//TaskType
const (
	TASK_ERROR		int32 = 0
	TASK_INSERT     int32 = 1
	TASK_SELECT_ALL int32 = 2
	TASK_FLUSH		int32 = 3
)

//NodeStatus
const (
	NODE_ERROR			int32 = 0
	NODE_CONNECTED		int32 = 1
	NODE_READY			int32 = 2
)
