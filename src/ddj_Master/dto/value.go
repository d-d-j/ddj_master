package dto

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//This is alias for type that is using to store data on node. It's representation could change over time so it should
//be used instead of current type
type Value float32

func (v Value) String() string {
	return fmt.Sprintf("[%f]", v)

}

func (v Value) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, v)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (v *Value) Decode(buf []byte) error {

	buffer := bytes.NewBuffer(buf)
	return binary.Read(buffer, binary.LittleEndian, v)
}

func (v Value) Size() int {

	return binary.Size(v)
}

func (v Value) Less(y Value) bool {
	return v < y
}
