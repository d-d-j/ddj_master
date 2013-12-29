package common

//Result
const (
	RESULT_OK  int32 = 0
	RESULT_ERR int32 = 1
)

//TaskType
const (
	TASK_ERROR  int32 = 0
	TASK_INSERT int32 = 1
	TASK_SELECT int32 = 2
	TASK_FLUSH  int32 = 3
	TASK_INFO   int32 = 4
)

//Aggregation Type
const (
	AGGREGATION_NONE    int32 = 0
	AGGREGATION_ADD     int32 = 1
	AGGREGATION_AVERAGE int32 = 2
)

//NodeStatus
const (
	NODE_ERROR     int32 = 0
	NODE_CONNECTED int32 = 1
	NODE_READY     int32 = 2
)

//Constants
const (
	CONST_UNINITIALIZED = -1
	CONST_INT_MIN_VALUE = -(int(^uint(0) >> 1))
)

//DEVICE
const (
	ALL_GPUs           = CONST_UNINITIALIZED
	TASK_UNINITIALIZED = CONST_UNINITIALIZED
)
