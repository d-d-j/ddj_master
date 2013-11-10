package dto

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type TaskRequestHeader struct {
	Id       int32
	Code     int32
	LoadSize int32
}

type Query struct {
	TaskRequestHeader
	Response chan []	Dto
	Load     			Dto
}

func (q *TaskRequestHeader) String() string {
	return fmt.Sprintf("#%d Code: %d [%X]", q.Id, q.Code, q.LoadSize)
}

func (q *TaskRequestHeader) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, q)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (q *TaskRequestHeader) Decode(buf []byte) error {

	buffer := bytes.NewBuffer(buf)
	return binary.Read(buffer, binary.LittleEndian, q)

}

func (q *TaskRequestHeader) Size() int {

	return binary.Size(q)
}
