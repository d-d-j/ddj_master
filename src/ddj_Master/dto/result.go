package dto

import (
	"fmt"
)

type Result struct {
	Header
	Data []byte
}

func NewResult(id int64, ttype int32, size int32, data []byte) *Result {
	r := new(Result)
	r.Header = Header{id, ttype, size}
	r.Data = data
	return r
}

func (r *Result) String() string {
	return fmt.Sprintf("Result with type %d and task id %d", r.Header.Type, r.Header.TaskId)
}

func (r *Result) Decode(buf []byte) error {
	return r.Header.Decode(buf)
}
