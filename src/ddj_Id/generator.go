package ddjId

import "sync/atomic"

type Int64Generator interface {
	getId() int64
}

type TaskIdGenerator struct {
	nextId int64
}

func (gen TaskIdGenerator) getId() int64 {
	return atomic.AddInt64(gen.nextId, 1)
}

