package dto

import (
	"fmt"
	log "code.google.com/p/log4go"
)

// Implements Encoder and Dto interface and is used to delegate a task to node
type Request struct {
	Header
	Data []byte
}

func NewRequest(id int64, ttype int32, size int32, data Dto) *Request {
	r := new(Request)
	r.Header = Header{id, ttype, size}
	var err error
	r.Data, err = data.Encode()
	if err != nil {
		log.Error(err)
	}
	return r
}

func (r *Request) String() string {
	return fmt.Sprintf("Request with type %d and task id %d", r.Header.Type, r.Header.TaskId)
}

func (r *Request) Encode() ([]byte, error) {

	var (
		headerBuf []byte
		err       error
	)

	// Encode header
	headerBuf, err = r.Header.Encode()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Merge header and data to one []byte buffer
	complete := make([]byte, int32(r.Header.Size())+r.Header.DataSize)
	copy(complete, headerBuf)
	copy(complete[len(headerBuf):], r.Data)

	return complete, err
}
