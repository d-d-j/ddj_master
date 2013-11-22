package dto

import (
	"fmt"
	"bytes"
	"encoding/binary"
	log "code.google.com/p/log4go"
)

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

func (r *Request) EncodeHeader() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, r.Header)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Request) DecodeHeader(buf []byte) error {

	return r.Header.Decode(buf)

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
	// TODO: CHANGE 100 to real data size
	complete := make([]byte, 100)
	copy(complete, headerBuf)
	copy(complete[len(headerBuf):], r.Data)

	return complete, err
}
