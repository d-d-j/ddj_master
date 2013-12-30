package dto

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

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
