package ddj_Dto

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Result struct {
	TaskRequestHeader
	Load []Dto
}

func (r *Result) String() string {
	return fmt.Sprintf("#%d Code: %d", r.Id, r.Code)
}

func (r *Result) EncodeHeader() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, r.TaskRequestHeader)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Result) DecodeHeader(buf []byte) error {

	return r.TaskRequestHeader.Decode(buf)

}
