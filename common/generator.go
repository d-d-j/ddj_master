package common

import "sync/atomic"

//Int64Generator is interface that should return unique id whenever GetId() is called. Return type should be int64
type Int64Generator interface {
	GetId() int64
}

//Int32Generator is interface that should return unique id whenever GetId() is called. Return type should be int32
type Int32Generator interface {
	GetId() int32
}

//TaskIdGenerator is implementation of Int64Generator
type TaskIdGenerator struct {
	nextId int64
}

//NewNodeIdGenerator it constructor of TaskIdGenerator
func NewTaskIdGenerator() *TaskIdGenerator {
	tig := new(TaskIdGenerator)
	tig.nextId = 0
	return tig
}

//Returns unique int64 number. This method is thread safe.
func (gen *TaskIdGenerator) GetId() int64 {
	return atomic.AddInt64(&gen.nextId, 1)
}

//NodeIdGenerator is implementation of Int32Generator
type NodeIdGenerator struct {
	nextId int32
}

//NewNodeIdGenerator it constructor of NodeIdGenerator
func NewNodeIdGenerator() *NodeIdGenerator {
	nig := new(NodeIdGenerator)
	nig.nextId = 0
	return nig
}

//Returns unique int32 number. This method is thread safe.
func (gen *NodeIdGenerator) GetId() int32 {
	return atomic.AddInt32(&gen.nextId, 1)
}
