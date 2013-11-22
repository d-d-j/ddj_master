package common

import "sync/atomic"

type Int64Generator interface {
	getId() int64
}

type Int32Generator interface {
	getId() int32
}

/* * * TASK ID GENERATOR * * */

type TaskIdGenerator struct {
	nextId int64
}

func NewTaskIdGenerator() *TaskIdGenerator {
	tig := new(TaskIdGenerator)
	tig.nextId = 0
	return tig
}

func (gen TaskIdGenerator) getId() int64 {
	return atomic.AddInt64(&gen.nextId, 1)
}

/* * * NODE ID GENERATOR * * */

type NodeIdGenerator struct {
	nextId int32
}

func NewNodeIdGenerator() *NodeIdGenerator {
	nig := new(NodeIdGenerator)
	nig.nextId = 0
	return nig
}

func (gen NodeIdGenerator) getId() int32 {
	return atomic.AddInt32(&gen.nextId, 1)
}
