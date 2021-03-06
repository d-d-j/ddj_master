package dto

import (
	"github.com/d-d-j/ddj_master/common"
	"fmt"
)

//This structure is represent RAW data returned by node
type Result struct {
	Header
	NodeId int32
	Data   []byte
}

func NewResult(id int64, nodeId int32, ttype int32, size int32, data []byte) *Result {
	r := new(Result)
	r.NodeId = nodeId
	r.Header = Header{id, ttype, size, common.CONST_UNINITIALIZED}
	r.Data = data
	return r
}

func (r *Result) String() string {
	return fmt.Sprintf("Result with type %d and task id %d", r.Header.Type, r.Header.TaskId)
}

func (r *Result) Decode(buf []byte) error {
	return r.Header.Decode(buf)
}
