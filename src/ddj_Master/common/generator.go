package common

import "sync/atomic"

type Int64Generator interface {
	getId() int64
}

type Int32Generator interface {
	getId() int32
}

type TaskIdGenerator struct {
	nextId int64
}

func (gen TaskIdGenerator) getId() int64 {
	return atomic.AddInt64(gen.nextId, 1)
}

type NodeIdGenerator struct {
	nextId int32
}

func (gen NodeIdGenerator) getId() int32 {
	return atomic.AddInt32(gen.nextId, 1)
}
